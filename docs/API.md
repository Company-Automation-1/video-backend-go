# API 文档

## 接口汇总表

| 方法 | 路径 | 鉴权 | 功能 | 请求体 |
|------|------|------|------|--------|
| GET | `/health` | ❌ 无 | 健康检查 | - |
| POST | `/api/v1/auth/user/login` | ❌ 无 | 用户登录 | ✅ |
| POST | `/api/v1/auth/admin/login` | ❌ 无 | 管理员登录 | ✅ |
| POST | `/api/v1/auth/logout` | ✅ 用户 | 用户登出 | - |
| POST | `/api/v1/users/send-verification-code` | ❌ 无 | 发送验证码 | ✅ |
| POST | `/api/v1/users/register` | ❌ 无 | 用户注册 | ✅ |
| GET | `/api/v1/users/profile` | 👤 本人 | 获取个人信息 | - |
| POST | `/api/v1/users/update-email` | ✅ 用户 | 更新邮箱 | ✅ |
| PUT | `/api/v1/users/:id` | 👤 本人 | 更新用户信息 | ✅ |
| DELETE | `/api/v1/users/:id` | 👤 本人 | 删除用户 | - |
| GET | `/api/v1/admin/admins` | 🔐 管理员 | 获取管理员列表 | - |
| GET | `/api/v1/admin/users` | 🔐 管理员 | 管理员获取用户列表 | - |
| GET | `/api/v1/admin/users/:id` | 🔐 管理员 | 管理员获取单个用户 | - |
| PUT | `/api/v1/admin/users/:id` | 🔐 管理员 | 管理员更新用户 | ✅ |

**鉴权说明：**
- ❌ 无：无需认证
- ✅ 用户：需要用户Token（AuthMiddleware）
- 👤 本人：需要用户Token且操作的是自己的数据（SelfMiddleware）
- 🔐 管理员：需要管理员Token（AdminMiddleware）

---

## 接口详情

### 1. 健康检查
```
GET /health
```
- 鉴权：❌ 无
- 请求体：无

---

### 2. 用户登录
```
POST /api/v1/auth/user/login
```
- 鉴权：❌ 无
- 请求体：
```json
{
  "username": "string",  // 必填
  "password": "string"   // 必填
}
```

---

### 3. 管理员登录
```
POST /api/v1/auth/admin/login
```
- 鉴权：❌ 无
- 请求体：
```json
{
  "username": "string",  // 必填
  "password": "string"   // 必填
}
```

---

### 4. 用户登出
```
POST /api/v1/auth/logout
Headers: Authorization: Bearer <access_token>
```
- 鉴权：✅ 用户
- 请求体：无

---

### 5. 发送验证码
```
POST /api/v1/users/send-verification-code
```
- 鉴权：❌ 无
- 请求体：
```json
{
  "email": "user@example.com"  // 必填，邮箱格式
}
```

---

### 6. 用户注册
```
POST /api/v1/users/register
```
- 鉴权：❌ 无
- 请求体：
```json
{
  "username": "string",     // 必填，3-100字符
  "email": "string",        // 必填，邮箱格式
  "password": "string",      // 必填，最少6位
  "captcha": "string"       // 必填，6位验证码
}
```

---

### 7. 获取个人信息
```
GET /api/v1/users/profile
Headers: Authorization: Bearer <access_token>
```
- 鉴权：👤 本人（从Token中获取用户ID，获取本人的个人信息）
- 请求体：无
- 路径参数：无

---

### 8. 更新邮箱
```
POST /api/v1/users/update-email
Headers: Authorization: Bearer <access_token>
```
- 鉴权：✅ 用户
- 请求体：
```json
{
  "email": "newemail@example.com",  // 必填，新邮箱
  "code": "123456"                  // 必填，6位验证码
}
```

---

### 9. 更新用户信息
```
PUT /api/v1/users/:id
Headers: Authorization: Bearer <access_token>
```
- 鉴权：👤 本人（路径参数id必须与Token中的用户ID一致）
- 路径参数：`id` (用户ID)
- 请求体：
```json
{
  "username": "string",  // 可选，3-100字符
  "password": "string",  // 可选，最少6位
  "points": 0            // 可选，积分值。null=置空（等同于0），0=置为0，其他数值=设置对应积分
}
```
- 说明：
  - 所有字段均为可选，只更新提供的字段
  - `points` 字段支持三种操作：
    - 不传：不更新积分字段
    - `null`：将积分置空（NULL，逻辑上等同于0）
    - `0` 或其他数值：设置对应的积分值

---

### 10. 删除用户
```
DELETE /api/v1/users/:id
Headers: Authorization: Bearer <access_token>
```
- 鉴权：👤 本人（路径参数id必须与Token中的用户ID一致）
- 路径参数：`id` (用户ID)
- 请求体：无

---

### 11. 管理员获取用户列表
```
GET /api/v1/admin/users?page=1&page_size=10
Headers: Authorization: Bearer <admin_access_token>
```
- 鉴权：🔐 管理员
- 请求体：无
- 路径参数：无
- 查询参数：
  - `page` (可选，默认1)：页码，从1开始
  - `page_size` (可选，默认10，最大100)：每页数量
- 响应体（分页格式）：
```json
{
  "code": 200,
  "success": true,
  "data": {
    "list": [
      {
        "id": 1,
        "username": "string",
        "email": "string",
        "points": 0,
        "created_at": 1234567890,
        "updated_at": 1234567890
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 10,
      "total": 100,
      "pages": 10
    }
  }
}
```

---

### 12. 管理员获取单个用户
```
GET /api/v1/admin/users/:id
Headers: Authorization: Bearer <admin_access_token>
```
- 鉴权：🔐 管理员
- 请求体：无
- 路径参数：`id` (用户ID)

---

### 13. 获取管理员列表
```
GET /api/v1/admin/admins?page=1&page_size=10
Headers: Authorization: Bearer <admin_access_token>
```
- 鉴权：🔐 管理员
- 请求体：无
- 路径参数：无
- 查询参数：
  - `page` (可选，默认1)：页码，从1开始
  - `page_size` (可选，默认10，最大100)：每页数量
- 响应体（分页格式）：
```json
{
  "code": 200,
  "success": true,
  "data": {
    "list": [
      {
        "id": 1,
        "username": "string",
        "created_at": 1234567890,
        "updated_at": 1234567890
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 10,
      "total": 100,
      "pages": 10
    }
  }
}
```

---

### 14. 管理员更新用户
```
PUT /api/v1/admin/users/:id
Headers: Authorization: Bearer <admin_access_token>
```
- 鉴权：🔐 管理员
- 路径参数：`id` (用户ID)
- 请求体：
```json
{
  "username": "string",  // 可选，3-100字符
  "password": "string",  // 可选，最少6位
  "points": 0            // 可选，积分值。null=置空（等同于0），0=置为0，其他数值=设置对应积分
}
```
- 说明：管理员可以更新任意用户的积分

---

## 注意事项

1. 所有需要鉴权的接口都需要在请求头中携带 `Authorization: Bearer <token>`
2. 请求体为 JSON 格式，Content-Type: application/json
3. 所有字段验证失败会返回 400 错误
4. 未知字段会被拒绝，返回 400 错误
5. 积分字段 `points` 为可选字段，支持 `null`、`0` 和其他数值
