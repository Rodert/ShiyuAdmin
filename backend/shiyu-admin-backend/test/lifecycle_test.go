package test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"shiyu-admin-backend/internal/bootstrap"
	"shiyu-admin-backend/internal/config"
	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/server"
	"shiyu-admin-backend/pkg/database"
)

const adminPassword = "Shiyu!Lifecycle2026"

type testApp struct {
	cfg     *config.Config
	handler http.Handler
	db      *gorm.DB
}

func newTestApp(t *testing.T) *testApp {
	t.Helper()
	cfg := &config.Config{}
	cfg.DataDir = t.TempDir()
	cfg.Server.Mode = "test"
	cfg.JWT.ExpireTime = 3600
	cfg.JWT.Issuer = "shiyu-test"
	cfg.Database.Driver = "sqlite"
	cfg.Database.Database = filepath.Join(cfg.DataDir, "app.db")
	if path := os.Getenv("TEST_CONFIG_FILE"); path != "" {
		loaded, err := config.Load(path)
		require.NoError(t, err)
		cfg.Database = loaded.Database
	}
	handler, err := server.NewEngine(cfg)
	require.NoError(t, err)
	db, err := database.Connect(cfg)
	require.NoError(t, err)
	return &testApp{cfg, handler, db}
}

func (a *testApp) request(t *testing.T, method, path, token string, body any, status int) map[string]any {
	t.Helper()
	var b bytes.Buffer
	if body != nil {
		require.NoError(t, json.NewEncoder(&b).Encode(body))
	}
	req := httptest.NewRequest(method, "/api/v1/system"+path, &b)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)
	require.Equal(t, status, w.Code, "%s %s: %s", method, path, w.Body.String())
	var out map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out))
	return out
}

func (a *testApp) login(t *testing.T, username, password string) string {
	t.Helper()
	r := a.request(t, "POST", "/auth/login", "", map[string]string{"username": username, "password": password}, 200)
	return r["data"].(map[string]any)["token"].(string)
}

func (a *testApp) install(t *testing.T) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(a.cfg.DataDir, "setup-token"))
	require.NoError(t, err)
	a.request(t, "POST", "/installation", "", map[string]string{"token": strings.TrimSpace(string(b)), "name": "测试管理系统", "username": "owner", "password": adminPassword}, 200)
}

func TestApplicationLifecycle(t *testing.T) {
	a := newTestApp(t)
	var admin string
	t.Run("installation authorization and one time setup", func(t *testing.T) {
		out := a.request(t, "GET", "/installation", "", nil, 200)["data"].(map[string]any)
		require.Equal(t, false, out["installed"])
		a.request(t, "POST", "/auth/login", "", map[string]string{"username": "admin", "password": "Admin@123"}, 503)
		a.request(t, "POST", "/installation", "", map[string]string{"token": "wrong", "name": "Demo", "username": "owner", "password": adminPassword}, 403)
		a.install(t)
		b, err := os.ReadFile(filepath.Join(a.cfg.DataDir, "setup-token"))
		require.NoError(t, err)
		a.request(t, "POST", "/installation", "", map[string]string{"token": strings.TrimSpace(string(b)), "name": "Overwrite", "username": "other", "password": adminPassword}, 409)
		var count int64
		require.NoError(t, a.db.Model(&entity.User{}).Count(&count).Error)
		require.EqualValues(t, 1, count)
		admin = a.login(t, "owner", adminPassword)
		a.request(t, "GET", "/profile", admin, nil, 200)
		a.request(t, "GET", "/version", admin, nil, 200)
		a.request(t, "GET", "/users", "", nil, 401)
	})
	if admin == "" {
		t.Fatal("installation failed")
	}
	var managerCode, otherCode string
	var manager string
	t.Run("department writes and protected administration", func(t *testing.T) {
		for _, dept := range []string{"DEPT_A", "DEPT_B"} {
			require.NoError(t, a.db.Create(&entity.Dept{DeptCode: dept, DeptName: dept, Status: 1}).Error)
		}
		role := entity.Role{RoleCode: "LIMITED", RoleName: "部门管理员", RoleKey: "limited", DataScope: "dept", Status: 1}
		require.NoError(t, a.db.Create(&role).Error)
		for _, name := range []string{"system-user", "system-user-create", "system-user-update", "system-user-delete", "system-user-assign-role", "system-file", "system-file-upload", "system-file-delete"} {
			require.NoError(t, a.db.Create(&entity.RoleMenu{RoleCode: "LIMITED", MenuCode: name}).Error)
		}
		create := func(username, dept string) string {
			r := a.request(t, "POST", "/users", admin, map[string]any{"username": username, "password": adminPassword, "dept_code": dept, "status": 1}, 200)
			return r["data"].(map[string]any)["user_code"].(string)
		}
		managerCode = create("manager", "DEPT_A")
		otherCode = create("outside", "DEPT_B")
		require.NoError(t, a.db.Create(&entity.UserRole{UserCode: managerCode, RoleCode: "LIMITED"}).Error)
		manager = a.login(t, "manager", adminPassword)
		rows := a.request(t, "GET", "/users", manager, nil, 200)["data"].(map[string]any)
		require.EqualValues(t, 1, rows["total"])
		a.request(t, "PUT", "/users/"+otherCode, manager, map[string]any{"nickname": "blocked"}, 404)
		a.request(t, "DELETE", "/users/"+otherCode, manager, nil, 404)
		a.request(t, "GET", "/users/"+otherCode+"/roles", manager, nil, 404)
		a.request(t, "PUT", "/users/"+managerCode, manager, map[string]any{"dept_code": "DEPT_B"}, 403)
		a.request(t, "POST", "/users", manager, map[string]any{"username": "escape", "password": adminPassword, "dept_code": "DEPT_B", "status": 1}, 403)
		a.request(t, "PUT", "/users/"+managerCode+"/roles", manager, map[string]any{"role_codes": []string{"ROLE_ADMIN"}}, 403)
		a.request(t, "PUT", "/users/"+managerCode, manager, map[string]any{"nickname": "allowed"}, 200)
		var owner entity.User
		require.NoError(t, a.db.Where("username = ?", "owner").First(&owner).Error)
		a.request(t, "DELETE", "/users/"+owner.UserCode, admin, nil, 403)
	})
	t.Run("disabled account and password revoke old tokens", func(t *testing.T) {
		a.request(t, "PUT", "/users/"+managerCode, admin, map[string]any{"status": 0}, 200)
		a.request(t, "GET", "/profile", manager, nil, 401)
		a.request(t, "PUT", "/users/"+managerCode, admin, map[string]any{"status": 1}, 200)
		a.request(t, "GET", "/profile", manager, nil, 401)
		manager = a.login(t, "manager", adminPassword)
		a.request(t, "PUT", "/profile/password", manager, map[string]string{"old_password": adminPassword, "new_password": "Shiyu!Changed2026"}, 200)
		a.request(t, "GET", "/profile", manager, nil, 401)
		manager = a.login(t, "manager", "Shiyu!Changed2026")
	})
	t.Run("file isolation search recycle and persistence", func(t *testing.T) {
		var b bytes.Buffer
		form := multipart.NewWriter(&b)
		part, err := form.CreateFormFile("file", "document.txt")
		require.NoError(t, err)
		_, err = part.Write([]byte("persistent file content"))
		require.NoError(t, err)
		require.NoError(t, form.Close())
		req := httptest.NewRequest("POST", "/api/v1/system/files/upload", &b)
		req.Header.Set("Content-Type", form.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+admin)
		w := httptest.NewRecorder()
		a.handler.ServeHTTP(w, req)
		require.Equal(t, 200, w.Code, w.Body.String())
		var out struct {
			Data entity.MediaFile `json:"data"`
		}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out))
		code := out.Data.FileCode
		require.NotEmpty(t, code)
		rows := a.request(t, "GET", "/files?keyword=DOC", admin, nil, 200)["data"].(map[string]any)
		require.EqualValues(t, 1, rows["total"])
		a.request(t, "GET", "/files/"+code+"/download", manager, nil, 404)
		a.request(t, "DELETE", "/files/"+code, manager, nil, 404)
		a.request(t, "POST", "/files/"+code+"/restore", manager, nil, 404)
		a.request(t, "DELETE", "/files/"+code, admin, nil, 200)
		a.request(t, "GET", "/files/"+code+"/download", admin, nil, 404)
		a.request(t, "POST", "/files/"+code+"/restore", admin, nil, 200)
		oldKey := a.cfg.JWT.Secret
		a.cfg.JWT.Secret = ""
		a.handler, err = server.NewEngine(a.cfg)
		require.NoError(t, err)
		require.Equal(t, oldKey, a.cfg.JWT.Secret)
		req = httptest.NewRequest("GET", "/api/v1/system/files/"+code+"/download", nil)
		req.Header.Set("Authorization", "Bearer "+admin)
		w = httptest.NewRecorder()
		a.handler.ServeHTTP(w, req)
		require.Equal(t, 200, w.Code)
		require.Equal(t, "persistent file content", w.Body.String())
	})
	t.Run("durable logout without Redis", func(t *testing.T) {
		a.request(t, "POST", "/auth/logout", manager, nil, 200)
		a.request(t, "GET", "/profile", manager, nil, 401)
		var err error
		a.handler, err = server.NewEngine(a.cfg)
		require.NoError(t, err)
		a.request(t, "GET", "/profile", manager, nil, 401)
	})
	t.Run("team settings and approval lifecycle", func(t *testing.T) { testWorkManagement(t, a, admin) })
	t.Run("file maintenance and resubmission", func(t *testing.T) { testFileMaintenance(t, a, admin) })
	t.Run("migration is idempotent and refuses downgrade", func(t *testing.T) {
		require.NoError(t, a.db.Model(&entity.Menu{}).Where("menu_code = ?", "welcome").Update("menu_name", "自定义首页").Error)
		require.NoError(t, bootstrap.Migrate(a.db, a.cfg))
		var menu entity.Menu
		require.NoError(t, a.db.Where("menu_code = ?", "welcome").First(&menu).Error)
		require.Equal(t, "自定义首页", menu.MenuName)
		require.NoError(t, a.db.Create(&bootstrap.Migration{Version: bootstrap.SchemaVersion + 1, AppliedAt: time.Now()}).Error)
		require.ErrorContains(t, bootstrap.Migrate(a.db, a.cfg), "newer")
		require.NoError(t, a.db.Delete(&bootstrap.Migration{}, bootstrap.SchemaVersion+1).Error)
	})
}

func TestSecretPersistenceAndValidation(t *testing.T) {
	cfg := &config.Config{}
	cfg.DataDir = t.TempDir()
	cfg.JWT.ExpireTime = 3600
	require.NoError(t, config.PrepareSecurity(cfg))
	original := cfg.JWT.Secret
	cfg.JWT.Secret = ""
	require.NoError(t, config.PrepareSecurity(cfg))
	require.Equal(t, original, cfg.JWT.Secret)
	info, err := os.Stat(filepath.Join(cfg.DataDir, "jwt.key"))
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0600), info.Mode().Perm())
	cfg.JWT.Secret = "shiyu-admin-scaffold-secret-key"
	require.Error(t, config.PrepareSecurity(cfg))
}

func TestLegacyUpgradePreservesDataAndRequiresPasswordChange(t *testing.T) {
	cfg := &config.Config{}
	cfg.DataDir = t.TempDir()
	cfg.Database.Driver = "sqlite"
	cfg.Database.Database = filepath.Join(cfg.DataDir, "legacy.db")
	cfg.JWT.ExpireTime = 3600
	cfg.Server.Mode = "test"
	db, err := database.Connect(cfg)
	require.NoError(t, err)
	require.NoError(t, bootstrap.AutoMigrate(db))
	hash, err := bcrypt.GenerateFromPassword([]byte("Admin@123"), bcrypt.DefaultCost)
	require.NoError(t, err)
	require.NoError(t, db.Create(&entity.User{UserCode: "legacy-admin", Username: "admin", Password: string(hash), Status: 1, IsSuperAdmin: true}).Error)
	require.NoError(t, db.Create(&entity.Menu{MenuCode: "welcome", MenuName: "旧版定制首页", MenuType: "C", Perms: "welcome:view", Status: 1}).Error)
	handler, err := server.NewEngine(cfg)
	require.NoError(t, err)
	a := &testApp{cfg, handler, db}
	status := a.request(t, "GET", "/installation", "", nil, 200)["data"].(map[string]any)
	require.Equal(t, true, status["installed"])
	var menu entity.Menu
	require.NoError(t, db.Where("menu_code = ?", "welcome").First(&menu).Error)
	require.Equal(t, "旧版定制首页", menu.MenuName)
	token := a.login(t, "admin", "Admin@123")
	profile := a.request(t, "GET", "/profile", token, nil, 200)["data"].(map[string]any)
	require.Equal(t, true, profile["require_password_change"])
	a.request(t, "GET", "/users", token, nil, 403)
	a.request(t, "PUT", "/profile/password", token, map[string]string{"old_password": "Admin@123", "new_password": adminPassword}, 200)
	a.request(t, "GET", "/profile", token, nil, 401)
	a.request(t, "GET", "/users", a.login(t, "admin", adminPassword), nil, 200)
}
