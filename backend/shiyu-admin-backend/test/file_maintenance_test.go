package test

import (
	"context"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/server"
	filesvc "shiyu-admin-backend/internal/service/files"
)

func testFileMaintenance(t *testing.T, a *testApp, admin string) {
	createUser := func(name string) string {
		data := a.request(t, "POST", "/users", admin, map[string]any{"username": name, "password": adminPassword, "status": 1}, 200)["data"].(map[string]any)
		return data["user_code"].(string)
	}
	ownerCode, otherCode := createUser("files-owner"), createUser("files-other")
	owner, other := a.login(t, "files-owner", adminPassword), a.login(t, "files-other", adminPassword)
	svc := &filesvc.Service{DB: a.db}
	t.Run("retention validation and compatibility", func(t *testing.T) {
		settings := a.request(t, "GET", "/settings", admin, nil, 200)["data"].(map[string]any)
		require.EqualValues(t, 30, settings["file_retention_days"])
		for _, invalid := range []int{-1, 0, 6, 366} {
			settings["file_retention_days"] = invalid
			a.request(t, "PUT", "/settings", admin, settings, 400)
		}
		settings["file_retention_days"] = 14
		a.request(t, "PUT", "/settings", owner, settings, 403)
		a.request(t, "PUT", "/settings", admin, settings, 200)
		delete(settings, "file_retention_days")
		preserved := a.request(t, "PUT", "/settings", admin, settings, 200)["data"].(map[string]any)
		require.EqualValues(t, 14, preserved["file_retention_days"])
		settings["file_retention_days"] = 30
		a.request(t, "PUT", "/settings", admin, settings, 200)
	})
	orphan := uploadWorkFile(t, a, owner, "unused file bytes", 200)
	linked := uploadWorkFile(t, a, owner, "approval evidence", 200)
	create := func(token, title, source string, codes []string, expected int) string {
		data := a.request(t, "POST", "/work/requests", token, map[string]any{"title": title, "content": "申请说明", "approver_code": otherCode, "file_codes": codes, "source_request_id": source}, expected)
		if expected != 200 {
			return ""
		}
		return data["data"].(map[string]any)["id"].(string)
	}
	id := create(owner, "采购 50%_[预算]", "", []string{linked}, 200)
	t.Run("search filters keep participant scope and escape wildcards", func(t *testing.T) {
		for _, keyword := range []string{"采购", "50%_[", "files-owner"} {
			data := a.request(t, "GET", "/work/requests?view=mine&keyword="+url.QueryEscape(keyword)+"&status=pending", owner, nil, 200)["data"].(map[string]any)
			require.EqualValues(t, 1, data["total"])
		}
		data := a.request(t, "GET", "/work/requests?view=mine&keyword="+url.QueryEscape("采购")+"&status=pending", admin, nil, 200)["data"].(map[string]any)
		require.EqualValues(t, 0, data["total"])
		a.request(t, "GET", "/work/requests?status=invalid", owner, nil, 400)
		create(owner, "非法重复", id, nil, 400)
		a.request(t, "POST", "/work/requests/"+id+"/actions", other, map[string]string{"action": "rejected", "comment": "请补充预算"}, 200)
		create(admin, "越权重提", id, nil, 404)
		again := create(owner, "补充预算后重提", id, []string{linked}, 200)
		detail := a.request(t, "GET", "/work/requests/"+again, owner, nil, 200)["data"].(map[string]any)
		require.Equal(t, id, detail["source_request_id"])
		original := a.request(t, "GET", "/work/requests/"+id, owner, nil, 200)["data"].(map[string]any)
		require.Equal(t, "rejected", original["status"])
		require.Len(t, original["events"].([]any), 2)
	})
	t.Run("own unused uploads and protected linked attachments", func(t *testing.T) {
		data := a.request(t, "GET", "/work/attachments", owner, nil, 200)["data"].(map[string]any)
		require.EqualValues(t, 1, data["total"])
		require.Equal(t, orphan, data["items"].([]any)[0].(map[string]any)["file_code"])
		data = a.request(t, "GET", "/work/attachments", other, nil, 200)["data"].(map[string]any)
		require.EqualValues(t, 0, data["total"])
		a.request(t, "DELETE", "/work/attachments/"+orphan, other, nil, 404)
		a.request(t, "DELETE", "/work/attachments/"+linked, owner, nil, 409)
		a.request(t, "GET", "/files/cleanup", owner, nil, 403)
		a.request(t, "POST", "/files/cleanup", owner, map[string]any{"file_codes": []string{orphan}}, 403)
		var general entity.MediaFile
		general.FileCode = "legacy-own-general"
		general.ObjectKey = "legacy-own-general"
		general.UploaderCode = ownerCode
		general.StorageID = 1
		general.OriginalName = "legacy.txt"
		general.MimeType = "text/plain"
		require.NoError(t, a.db.Create(&general).Error)
		a.request(t, "DELETE", "/work/attachments/"+general.FileCode, owner, nil, 404)
		usage := a.request(t, "GET", "/work/attachments/usage", owner, nil, 200)["data"].(map[string]any)
		require.EqualValues(t, 1, usage["unused"].(map[string]any)["count"])
		require.EqualValues(t, len("unused file bytes"), usage["unused"].(map[string]any)["bytes"])
		a.request(t, "DELETE", "/work/attachments/"+orphan, owner, nil, 200)
		err := svc.Purge(context.Background(), orphan)
		require.Error(t, err)
		a.request(t, "POST", "/work/attachments/"+orphan+"/restore", owner, nil, 200)
		a.request(t, "DELETE", "/work/attachments/"+orphan, owner, nil, 200)
	})
	age := func(code string) {
		require.NoError(t, a.db.Unscoped().Model(&entity.MediaFile{}).Where("file_code = ?", code).Update("deleted_at", time.Now().AddDate(0, 0, -31)).Error)
	}
	purge := func(code string, expected bool) {
		out := a.request(t, "POST", "/files/cleanup", admin, map[string]any{"file_codes": []string{code}}, 200)["data"].([]any)
		require.Len(t, out, 1)
		require.Equal(t, expected, out[0].(map[string]any)["purged"])
	}
	t.Run("purge claims survive failure and restart", func(t *testing.T) {
		purge(orphan, false)
		age(orphan)
		var storage entity.StorageConfig
		require.NoError(t, a.db.Where("is_default = ?", true).First(&storage).Error)
		require.NoError(t, a.db.Model(&storage).Update("base_path", filepath.Join(a.cfg.DataDir, "not-mounted")).Error)
		purge(orphan, false)
		a.request(t, "POST", "/work/attachments/"+orphan+"/restore", owner, nil, 404)
		create(owner, "不可引用清理中附件", "", []string{orphan}, 400)
		require.NoError(t, a.db.Model(&storage).Update("base_path", filepath.Join(a.cfg.DataDir, "uploads")).Error)
		var err error
		a.handler, err = server.NewEngine(a.cfg)
		require.NoError(t, err)
		purge(orphan, true)
		purge(orphan, true)
		var record entity.MediaFile
		require.NoError(t, a.db.Unscoped().Where("file_code = ?", orphan).First(&record).Error)
		require.NotNil(t, record.PurgedAt)
		require.False(t, record.PurgePending)
		_, err = os.Stat(filepath.Join(a.cfg.DataDir, "uploads", record.ObjectKey))
		require.True(t, errors.Is(err, os.ErrNotExist))
		a.request(t, "POST", "/files/"+orphan+"/restore", admin, nil, 404)
		usage := a.request(t, "GET", "/work/attachments/usage", owner, nil, 200)["data"].(map[string]any)
		require.EqualValues(t, 0, usage["recycled"].(map[string]any)["count"])
	})
	t.Run("stale cleanup selection cannot delete restored or linked files", func(t *testing.T) {
		code := uploadWorkFile(t, a, owner, "retain me", 200)
		a.request(t, "DELETE", "/work/attachments/"+code, owner, nil, 200)
		age(code)
		a.request(t, "GET", "/files/cleanup", admin, nil, 200)
		a.request(t, "POST", "/work/attachments/"+code+"/restore", owner, nil, 200)
		purge(code, false)
		age(linked)
		purge(linked, false)
		require.NoError(t, a.db.Unscoped().Model(&entity.MediaFile{}).Where("file_code = ?", linked).Update("deleted_at", nil).Error)
	})
	t.Run("final database write failure is safely retryable", func(t *testing.T) {
		code := uploadWorkFile(t, a, owner, "retry proof", 200)
		require.NoError(t, svc.Recycle(context.Background(), filesvc.Scope{User: ownerCode, WorkOnly: true}, code))
		age(code)
		hook := "test:fail_purge_finalization"
		require.NoError(t, a.db.Callback().Update().Before("gorm:update").Register(hook, func(tx *gorm.DB) {
			if tx.Statement.Table == "sys_media_files" {
				if values, ok := tx.Statement.Dest.(map[string]any); ok {
					if _, ok := values["purged_at"]; ok {
						tx.AddError(errors.New("simulated database failure"))
					}
				}
			}
		}))
		err := svc.Purge(context.Background(), code)
		require.NoError(t, a.db.Callback().Update().Remove(hook))
		require.ErrorContains(t, err, "simulated")
		require.Error(t, svc.Restore(context.Background(), filesvc.Scope{User: ownerCode, WorkOnly: true}, code))
		require.NoError(t, svc.Purge(context.Background(), code))
	})
	t.Run("storage traversal is blocked", func(t *testing.T) {
		outside := filepath.Join(a.cfg.DataDir, "outside-proof")
		require.NoError(t, os.WriteFile(outside, []byte("must survive"), 0600))
		var storage entity.StorageConfig
		require.NoError(t, a.db.Where("is_default = ?", true).First(&storage).Error)
		bad := entity.MediaFile{FileCode: "traversal-test", ObjectKey: "../outside-proof", UploaderCode: ownerCode, StorageID: storage.ID, OriginalName: "outside", MimeType: "text/plain", Purpose: "work", DeletedAt: gorm.DeletedAt{Time: time.Now().AddDate(0, 0, -31), Valid: true}}
		require.NoError(t, a.db.Create(&bad).Error)
		purge(bad.FileCode, false)
		contents, err := os.ReadFile(outside)
		require.NoError(t, err)
		require.Equal(t, "must survive", string(contents))
		require.NoError(t, os.Symlink(a.cfg.DataDir, filepath.Join(storage.BasePath, "escape")))
		bad.ID = 0
		bad.FileCode = "symlink-test"
		bad.ObjectKey = "escape/outside-proof"
		require.NoError(t, a.db.Create(&bad).Error)
		purge(bad.FileCode, false)
		contents, err = os.ReadFile(outside)
		require.NoError(t, err)
		require.Equal(t, "must survive", string(contents))
	})
}
