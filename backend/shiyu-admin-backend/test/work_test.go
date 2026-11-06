package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"shiyu-admin-backend/internal/service/work"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"shiyu-admin-backend/internal/bootstrap"
	"shiyu-admin-backend/internal/config"
	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/server"
	"shiyu-admin-backend/pkg/database"
)

func uploadWorkFile(t *testing.T, a *testApp, token, content string, status int) string {
	t.Helper()
	var b bytes.Buffer
	form := multipart.NewWriter(&b)
	part, err := form.CreateFormFile("file", "approval-evidence.txt")
	require.NoError(t, err)
	_, err = part.Write([]byte(content))
	require.NoError(t, err)
	require.NoError(t, form.Close())
	req := httptest.NewRequest("POST", "/api/v1/system/work/attachments/upload", &b)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", form.FormDataContentType())
	w := httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)
	require.Equal(t, status, w.Code, w.Body.String())
	if status != 200 {
		return ""
	}
	var out struct {
		Data entity.MediaFile `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out))
	return out.Data.FileCode
}

func testWorkManagement(t *testing.T, a *testApp, admin string) {
	var applicantCode, approverCode string
	for _, username := range []string{"work-applicant", "work-approver"} {
		out := a.request(t, "POST", "/users", admin, map[string]any{"username": username, "nickname": username, "password": adminPassword, "status": 1}, 200)
		code := out["data"].(map[string]any)["user_code"].(string)
		if username == "work-applicant" {
			applicantCode = code
		} else {
			approverCode = code
		}
	}
	applicant, approver := a.login(t, "work-applicant", adminPassword), a.login(t, "work-approver", adminPassword)
	t.Run("settings validation authorization audit and upload enforcement", func(t *testing.T) {
		a.request(t, "GET", "/settings", applicant, nil, 403)
		a.request(t, "GET", "/settings/changes", approver, nil, 403)
		settings := map[string]any{"name": "团队工作系统", "logo_url": "/logo-v2.png", "timezone": "UTC", "max_upload_mb": 1}
		a.request(t, "PUT", "/settings", applicant, settings, 403)
		for _, logo := range []string{"javascript:alert(1)", "//attacker.invalid/x", "/\\attacker.invalid/x"} {
			settings["logo_url"] = logo
			a.request(t, "PUT", "/settings", admin, settings, 400)
		}
		settings["logo_url"] = "/logo-v2.png"
		settings["timezone"] = "Fake/Zone"
		a.request(t, "PUT", "/settings", admin, settings, 400)
		settings["timezone"] = "UTC"
		settings["max_upload_mb"] = 0
		a.request(t, "PUT", "/settings", admin, settings, 400)
		settings["max_upload_mb"] = 1
		a.request(t, "PUT", "/settings", admin, settings, 200)
		info := a.request(t, "GET", "/installation", "", nil, 200)["data"].(map[string]any)
		require.Equal(t, "团队工作系统", info["name"])
		require.Equal(t, "UTC", info["timezone"])
		changes := a.request(t, "GET", "/settings/changes", admin, nil, 200)["data"].(map[string]any)
		require.EqualValues(t, 1, changes["total"])
		uploadWorkFile(t, a, applicant, strings.Repeat("a", 1024*1024+1), 400)
	})
	fileCode := uploadWorkFile(t, a, applicant, "signed application evidence", 200)
	create := func(status int, who, reviewer string, files []string) string {
		out := a.request(t, "POST", "/work/requests", who, map[string]any{"title": "采购事项", "content": "申请采购办公设备", "approver_code": reviewer, "file_codes": files}, status)
		if status != 200 {
			return ""
		}
		return out["data"].(map[string]any)["id"].(string)
	}
	create(400, applicant, applicantCode, nil)
	create(400, approver, applicantCode, []string{fileCode})
	create(400, applicant, approverCode, []string{fileCode, fileCode})
	create(400, applicant, "missing-user", nil)
	id := create(200, applicant, approverCode, []string{fileCode})
	path := "/work/requests/" + id
	t.Run("participant privacy and retained attachments", func(t *testing.T) {
		a.request(t, "GET", path, admin, nil, 404)
		a.request(t, "POST", path+"/actions", admin, map[string]string{"action": "approved"}, 404)
		a.request(t, "POST", path+"/actions", applicant, map[string]string{"action": "approved"}, 403)
		a.request(t, "POST", path+"/actions", approver, map[string]string{"action": "withdrawn"}, 403)
		a.request(t, "GET", path+"/attachments/"+fileCode, admin, nil, 404)
		a.request(t, "GET", path+"/attachments/not-linked", approver, nil, 404)
		a.request(t, "DELETE", "/files/"+fileCode, admin, nil, 409)
		req := httptest.NewRequest("GET", "/api/v1/system"+path+"/attachments/"+fileCode, nil)
		req.Header.Set("Authorization", "Bearer "+approver)
		w := httptest.NewRecorder()
		a.handler.ServeHTTP(w, req)
		require.Equal(t, 200, w.Code)
		require.Equal(t, "signed application evidence", w.Body.String())
		people := a.request(t, "GET", "/work/people?keyword=work-applicant", approver, nil, 200)["data"].([]any)
		require.Len(t, people, 1)
		require.Len(t, people[0].(map[string]any), 3)
		mine := a.request(t, "GET", "/work/requests?view=mine", approver, nil, 200)["data"].(map[string]any)
		require.EqualValues(t, 0, mine["total"])
		pending := a.request(t, "GET", "/work/requests?view=pending", approver, nil, 200)["data"].(map[string]any)
		require.EqualValues(t, 1, pending["total"])
	})
	t.Run("notifications are private and reads idempotent", func(t *testing.T) {
		data := a.request(t, "GET", "/work/notifications", approver, nil, 200)["data"].(map[string]any)
		require.EqualValues(t, 1, data["unread"])
		notice := data["items"].([]any)[0].(map[string]any)
		readPath := fmt.Sprintf("/work/notifications/%.0f/read", notice["id"].(float64))
		a.request(t, "PUT", readPath, applicant, nil, 404)
		a.request(t, "PUT", readPath, approver, nil, 200)
		a.request(t, "PUT", readPath, approver, nil, 200)
		data = a.request(t, "GET", "/work/notifications?unread=true", approver, nil, 200)["data"].(map[string]any)
		require.EqualValues(t, 0, data["unread"])
		require.EqualValues(t, 0, data["total"])
	})
	t.Run("concurrent decisions produce one outcome and one notification", func(t *testing.T) {
		a.request(t, "POST", path+"/actions", approver, map[string]string{"action": "rejected"}, 400)
		var wg sync.WaitGroup
		results := make(chan int, 2)
		for _, action := range []string{"approved", "rejected"} {
			wg.Add(1)
			go func(action string) {
				defer wg.Done()
				body, _ := json.Marshal(map[string]string{"action": action, "comment": "处理意见"})
				req := httptest.NewRequest("POST", "/api/v1/system"+path+"/actions", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+approver)
				w := httptest.NewRecorder()
				a.handler.ServeHTTP(w, req)
				results <- w.Code
			}(action)
		}
		wg.Wait()
		close(results)
		statuses := []int{}
		for status := range results {
			statuses = append(statuses, status)
		}
		require.ElementsMatch(t, []int{200, 409}, statuses)
		detail := a.request(t, "GET", path, applicant, nil, 200)["data"].(map[string]any)
		require.Len(t, detail["events"].([]any), 2)
		data := a.request(t, "GET", "/work/notifications", applicant, nil, 200)["data"].(map[string]any)
		require.EqualValues(t, 1, data["total"])
		a.request(t, "POST", path+"/actions", applicant, map[string]string{"action": "withdrawn"}, 409)
	})
	t.Run("notification failure rolls back decision and event", func(t *testing.T) {
		rollbackID := create(200, applicant, approverCode, nil)
		hook := "test:reject_notification"
		require.NoError(t, a.db.Callback().Create().Before("gorm:create").Register(hook, func(tx *gorm.DB) {
			if tx.Statement.Table == "sys_notifications" {
				tx.AddError(errors.New("simulated notification storage failure"))
			}
		}))
		svc := &work.Service{DB: a.db}
		err := svc.Act(context.Background(), approverCode, rollbackID, "approved", "")
		require.NoError(t, a.db.Callback().Create().Remove(hook))
		require.ErrorContains(t, err, "simulated notification")
		detail := a.request(t, "GET", "/work/requests/"+rollbackID, applicant, nil, 200)["data"].(map[string]any)
		require.Equal(t, "pending", detail["status"])
		require.Len(t, detail["events"].([]any), 1)
		a.request(t, "POST", "/work/requests/"+rollbackID+"/actions", approver, map[string]string{"action": "approved"}, 200)
	})
	t.Run("withdraw rejection and persistence", func(t *testing.T) {
		withdrawn := create(200, applicant, approverCode, nil)
		a.request(t, "POST", "/work/requests/"+withdrawn+"/actions", applicant, map[string]string{"action": "withdrawn"}, 200)
		a.request(t, "POST", "/work/requests/"+withdrawn+"/actions", approver, map[string]string{"action": "approved"}, 409)
		rejected := create(200, applicant, approverCode, nil)
		a.request(t, "POST", "/work/requests/"+rejected+"/actions", approver, map[string]string{"action": "rejected", "comment": "请补充预算"}, 200)
		a.request(t, "PUT", "/work/notifications/read-all", approver, nil, 200)
		var err error
		a.handler, err = server.NewEngine(a.cfg)
		require.NoError(t, err)
		detail := a.request(t, "GET", "/work/requests/"+rejected, applicant, nil, 200)["data"].(map[string]any)
		require.Equal(t, "rejected", detail["status"])
		data := a.request(t, "GET", "/work/notifications", approver, nil, 200)["data"].(map[string]any)
		require.EqualValues(t, 0, data["unread"])
		settings := a.request(t, "GET", "/settings", admin, nil, 200)["data"].(map[string]any)
		require.Equal(t, "团队工作系统", settings["name"])
		a.request(t, "PUT", "/users/"+approverCode, admin, map[string]any{"status": 0}, 200)
		create(400, applicant, approverCode, nil)
	})
}

func TestUpgradeFromSchemaOne(t *testing.T) {
	cfg := &config.Config{}
	cfg.DataDir = t.TempDir()
	cfg.Database.Driver = "sqlite"
	cfg.Database.Database = filepath.Join(cfg.DataDir, "v1.db")
	db, err := database.Connect(cfg)
	require.NoError(t, err)
	require.NoError(t, bootstrap.AutoMigrate(db))
	require.NoError(t, db.AutoMigrate(&bootstrap.Migration{}))
	// The old installation table intentionally has no settings columns.
	require.NoError(t, db.Exec("CREATE TABLE sys_installation (id integer PRIMARY KEY, name text NOT NULL, installed_at datetime)").Error)
	require.NoError(t, db.Exec("INSERT INTO sys_installation (id, name, installed_at) VALUES (?, ?, ?)", 1, "旧系统名称", time.Now()).Error)
	require.NoError(t, db.Create(&bootstrap.Migration{Version: 1, AppliedAt: time.Now()}).Error)
	require.NoError(t, db.Create(&entity.Menu{MenuCode: "custom-menu", MenuName: "定制菜单", Status: 1}).Error)
	require.NoError(t, bootstrap.Migrate(db, cfg))
	require.NoError(t, bootstrap.Migrate(db, cfg))
	var settings entity.Installation
	require.NoError(t, db.First(&settings, 1).Error)
	require.Equal(t, "旧系统名称", settings.Name)
	require.Equal(t, "Asia/Shanghai", settings.Timezone)
	require.Equal(t, 2048, settings.MaxUploadMB)
	var count int64
	require.NoError(t, db.Model(&bootstrap.Migration{}).Count(&count).Error)
	require.EqualValues(t, 2, count)
	require.NoError(t, db.Model(&entity.Menu{}).Where("menu_code = ?", "custom-menu").Count(&count).Error)
	require.EqualValues(t, 1, count)
	require.True(t, db.Migrator().HasTable(&entity.WorkRequest{}))
}

func TestSQLiteBackupRestore(t *testing.T) {
	cfg := &config.Config{}
	cfg.DataDir = t.TempDir()
	cfg.Database.Driver = "sqlite"
	cfg.Database.Database = filepath.Join(cfg.DataDir, "app.db")
	cfg.JWT.ExpireTime = 3600
	cfg.JWT.Issuer = "restore-test"
	cfg.Server.Mode = "test"
	handler, err := server.NewEngine(cfg)
	require.NoError(t, err)
	db, err := database.Connect(cfg)
	require.NoError(t, err)
	a := &testApp{cfg, handler, db}
	a.install(t)
	revoked := a.login(t, "owner", adminPassword)
	file := uploadWorkFile(t, a, revoked, "backup restores this attachment", 200)
	a.request(t, "POST", "/auth/logout", revoked, nil, 200)
	token := a.login(t, "owner", adminPassword)
	restoredRoot := t.TempDir()
	// SQLite's own snapshot API includes committed WAL contents.
	restoredDB := filepath.Join(restoredRoot, "app.db")
	require.NoError(t, db.Exec("VACUUM INTO ?", restoredDB).Error)
	for _, name := range []string{"jwt.key", "setup-token"} {
		data, err := os.ReadFile(filepath.Join(cfg.DataDir, name))
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(filepath.Join(restoredRoot, name), data, 0600))
	}
	var metadata entity.MediaFile
	require.NoError(t, db.Where("file_code = ?", file).First(&metadata).Error)
	originalFile := filepath.Join(cfg.DataDir, "uploads", metadata.ObjectKey)
	target := filepath.Join(restoredRoot, "uploads", metadata.ObjectKey)
	content, err := os.ReadFile(originalFile)
	require.NoError(t, err)
	require.NoError(t, os.MkdirAll(filepath.Dir(target), 0750))
	require.NoError(t, os.WriteFile(target, content, 0600))
	require.NoError(t, os.Remove(originalFile)) // Must serve restored bytes, not the original storage.
	restoredCfg := *cfg
	restoredCfg.DataDir = restoredRoot
	restoredCfg.Database.Database = restoredDB
	restoredCfg.JWT.Secret = ""
	restoredConnection, err := database.Connect(&restoredCfg)
	require.NoError(t, err)
	// This isolated restore uses a different mount path; remap the local storage location.
	require.NoError(t, restoredConnection.Model(&entity.StorageConfig{}).Where("driver = ?", "local").Update("base_path", filepath.Join(restoredRoot, "uploads")).Error)
	restoredHandler, err := server.NewEngine(&restoredCfg)
	require.NoError(t, err)
	restored := &testApp{&restoredCfg, restoredHandler, restoredConnection}
	require.Equal(t, cfg.JWT.Secret, restoredCfg.JWT.Secret)
	restored.request(t, "GET", "/profile", revoked, nil, 401)
	restored.request(t, "GET", "/profile", token, nil, 200)
	req := httptest.NewRequest("GET", "/api/v1/system/files/"+file+"/download", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	restored.handler.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)
	require.Equal(t, "backup restores this attachment", w.Body.String())
}
