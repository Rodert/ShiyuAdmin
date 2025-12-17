# Shiyu Admin Backend

后端服务：`backend/shiyu-admin-backend`

> 本文档主要面向后端开发者，说明认证、权限模型，以及「超级管理员」和账号状态字段的行为，方便二次开发和排查问题。

---

## 1. 认证与权限概览

- **JWT 认证**
  - 生成：`internal/service/auth/service.go` → `AuthService.Login`
  - 工具：`pkg/jwtutil`
  - 中间件：`internal/middleware/auth.go`

- **权限校验（RBAC）**
  - 关系：用户 ↔ 角色 ↔ 菜单（多对多）
  - 权限中间件：`internal/middleware/permission.go`
  - 菜单与权限标识（`perms`）：`internal/model/entity/menu.go`
  - 菜单树接口：`internal/api/v1/system/menus.go` 中的 `listMenuTree`

---

## 2. 用户模型（User）

实体定义：`internal/model/entity/user.go`

```go
// User represents a system user.
type User struct {
    ID           int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID，自增"`
    UserCode     string    `json:"user_code" gorm:"column:user_code;size:32;uniqueIndex;comment:用户编码，业务唯一标识，用于表关联"`
    Username     string    `json:"username" gorm:"size:64;uniqueIndex;comment:用户名，登录账号，唯一"`
    Nickname     string    `json:"nickname" gorm:"size:64;comment:昵称，显示名称"`
    Email        string    `json:"email" gorm:"size:128;comment:邮箱地址"`
    Phone        string    `json:"phone" gorm:"size:32;comment:手机号码"`
    Password     string    `json:"-" gorm:"size:255;comment:密码，BCrypt加密存储"`
    DeptCode     string    `json:"dept_code" gorm:"size:32;comment:部门编码，关联部门表"`
    Status       int       `json:"status" gorm:"comment:状态：1=启用，0=禁用"`
    IsSuperAdmin bool      `json:"is_super_admin" gorm:"column:is_super_admin;not null;default:false"`
    CreatedAt    time.Time `json:"created_at" gorm:"comment:创建时间"`
    UpdatedAt    time.Time `json:"updated_at" gorm:"comment:更新时间"`
    DeletedAt    *time.Time `json:"deleted_at" gorm:"index;comment:删除时间，软删除标记"`
}
```

关键字段：

- `Status`：账号状态
  - `1` = 启用
  - `0` = 禁用（锁定），**禁止登录**
- `IsSuperAdmin`：是否为超级管理员
  - `true`：在权限中间件和菜单接口中**跳过所有基于角色/菜单的权限检查**
  - `false`：按普通用户流程走 RBAC

表名：`sys_users`（见 `TableName()`）。

---

## 3. 超级管理员（Super Admin）

### 3.1 默认超管账号来源

配置文件：

- `configs/config.yaml`
- `configs/config.docker.yaml`

字段：

```yaml
bootstrap:
  admin_username: "admin"
  admin_password: "Admin@123"
  admin_nickname: "超级管理员"
```

初始化逻辑：`internal/bootstrap/bootstrap.go` → `EnsureAdminUser`

行为：

1. 启动时会根据 `bootstrap.admin_username` 查询用户表：
   - 如果已存在该用户名：
     - 只做 **同步超管标记**：
       - `is_super_admin = true`
   - 如果不存在该用户名：
     - 创建一条新用户记录：
       - `username = admin`（或配置中的值）
       - `password = Admin@123`（BCrypt 加密后存储）
       - `nickname = "超级管理员"`
       - `status = 1`
       - `is_super_admin = true`
2. 配合 `AutoMigrate`，会自动为 `sys_users` 增加 `is_super_admin` 列（如不存在）。

> 结论：默认情况下，配置中的 `admin` 账号会被视为**系统级超级管理员**。

### 3.2 JWT 中的超管标记

实现位置：`pkg/jwtutil/token.go`

- Claims 结构：

  ```go
  type Claims struct {
      UserCode     string `json:"user_code"`
      Username     string `json:"username"`
      IsSuperAdmin bool   `json:"is_super_admin"`
      jwt.RegisteredClaims
  }
  ```

- 生成 Token：`GenerateToken` 接收 `isSuperAdmin` 参数，并写入 Claims：

  ```go
  func GenerateToken(secret, issuer, userCode, username string, isSuperAdmin bool, expireSeconds int64) (string, error)
  ```

- 登录逻辑：`internal/service/auth/service.go` → `Login`：

  ```go
  token, err := jwtutil.GenerateToken(
      s.jwtSecret,
      s.jwtIssuer,
      user.UserCode,
      user.Username,
      user.IsSuperAdmin,
      s.jwtExpireIn,
  )
  ```

> 任何设置了 `IsSuperAdmin = true` 的用户，在成功登录后，其 JWT 中会携带 `"is_super_admin": true` 的 Claim。

### 3.3 权限中间件中的超管直通

实现位置：`internal/middleware/permission.go`

在以下三个中间件中，均在解析 JWT 并拿到 `claims` 后第一时间判断：

- `RequirePermission`
- `RequireAnyPermission`
- `RequireAllPermissions`

逻辑类似：

```go
claims, ok := claimsVal.(*jwtutil.Claims)
if !ok {
    // 401 处理
}
if claims.IsSuperAdmin {
    c.Next()
    return
}

// 非超管用户才继续走权限检查逻辑
```

效果：

- **超级管理员用户**：
  - 不再调用 `PermissionService` 进行权限计算，所有挂了权限中间件的接口一律通过。
- **普通用户**：
  - 仍然基于 用户 ↔ 角色 ↔ 菜单 的关联，计算是否拥有某个接口所需的 `perms` 标识。

### 3.4 菜单树中的超管视图

实现位置：`internal/api/v1/system/menus.go` → `listMenuTree`

关键逻辑（简化）：

```go
filtered := menus
if userRoleSvc != nil && roleMenuSvc != nil {
    if claimsVal, exists := c.Get(middleware.CurrentUserCtxKey); exists {
        if claims, ok := claimsVal.(*jwtutil.Claims); ok && claims.UserCode != "" {
            if claims.IsSuperAdmin {
                filtered = menus // 超管看到所有菜单
            } else {
                // 按角色/菜单关系过滤出允许的菜单编码
                allowedCodes, err := collectUserMenuCodes(...)
                // ... 过滤 + 补全父级菜单
            }
        }
    }
}
```

- 超管（`IsSuperAdmin = true`）：前端菜单树接口直接返回**完整菜单列表**。
- 普通用户：仅返回与其角色关联的菜单，并补全父级菜单形成完整树。

---

## 4. 账号状态（Status）与登录限制

### 4.1 字段含义

- 用户表字段：`status`（`int`）
  - `1` = 启用
  - `0` = 禁用（账号锁定）

实体注释（`entity.User`）：

```go
// Status 状态：1=启用，0=禁用
Status int `json:"status" gorm:"comment:状态：1=启用，0=禁用"`
```

### 4.2 登录流程中的 status 检查

实现位置：`internal/service/auth/service.go` → `Login`

关键逻辑（简化）：

```go
user, err := s.repo.FindUserByUsername(ctx, req.Username)
// ... 用户存在 & 密码校验通过后：

if user.Status != 1 {
    return nil, errors.New("账号已停用")
}

// 仅 status=1 的账号才会继续生成 Token
```

行为说明：

- 当 `user.Status != 1`（即 `0` 或其他非 1 值）时：
  - 登录接口直接返回错误 `"账号已停用"`。
  - **不会生成 JWT**，也不会进入后续权限/菜单逻辑。
- 此逻辑对所有用户生效，**包括超级管理员**：
  - 也就是说，可以通过将超管账号的 `status` 设为 `0` 来临时锁定该账号。

### 4.3 推荐使用方式

- **禁用账号**：
  - 推荐通过更新 `status = 0` 的方式实现，而不是物理删除用户记录。
  - 好处：
    - 方便后续恢复（`status = 1` 即可重新启用）
    - 便于审计和追踪历史操作人。

- **禁用超级管理员（admin）**：
  - 当确有需要锁死默认超管账号时，可直接把 `admin` 的 `status` 改为 `0`；
  - 注意同时确保还有其他可以登录的高权限账号，否则会需要通过数据库手动恢复。

---

## 5. 开发者建议

- **新增超管用户**：
  - 直接在 `sys_users` 中新增用户，并将 `is_super_admin` 设为 `true`，即可拥有与默认 admin 相同的权限行为。

- **迁移到自有项目时**：
  - 保留 `IsSuperAdmin` 与 `Status` 两个字段及其语义，有助于：
    - 提供一条“永远能救场”的超管账号（前提是不被禁用）。
    - 为运维留出“紧急锁定账号”的手段。

- **排查权限问题时**：
  - 优先检查：
    1. 当前登录用户的 `status` 是否为 `1`；
    2. 是否为 `is_super_admin = true`；
    3. 非超管用户再看角色与菜单的关联关系。

如需调整权限模型（例如改为更细粒度的按钮级权限、数据权限等），可以在保留现有 `is_super_admin` 与 `status` 语义的前提下扩展 `Menu` / `Role` 实体及相应 Service/Repository。
