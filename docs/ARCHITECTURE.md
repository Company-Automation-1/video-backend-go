# 架构设计文档

## 各层职责

### 1. DTO层 (`src/dto/`)
**职责：入参校验和转换**
- 定义请求DTO（Request DTO）
- 通过 `binding` 标签进行参数校验
- 提供 `ToModel()` 方法将DTO转换为模型
- **不负责出参处理**

**示例：**
- `UserCreateRequest` - 创建用户请求
- `UserUpdateRequest` - 更新用户请求

### 2. VO层 (`src/vo/`)
**职责：出参格式化**
- 定义值对象（Value Object）
- 从模型转换为VO，隐藏敏感字段（如密码）
- 提供 `FromModel()` 和 `FromModelList()` 方法
- **不负责入参处理**

**示例：**
- `UserVO` - 用户值对象

### 3. 控制器层 (`src/controllers/`)
**职责：请求处理和响应转换**
- 接收HTTP请求
- 调用 `ShouldBindJSON` 进行DTO校验（参数错误会panic）
- DTO → 模型转换（调用 `dto.ToModel()`）
- 调用服务层
- 模型 → VO转换（调用 `vo.FromModel()`）
- 设置响应到上下文（`middleware.SetData()`）

### 4. 服务层 (`src/services/`)
**职责：业务逻辑和ORM调用**
- 处理业务逻辑
- 调用ORM层操作数据库
- 返回模型或错误
- **不处理DTO/VO转换**

### 5. ORM层 (`src/query/`)
**职责：数据库操作**
- GORM Gen自动生成的类型安全查询
- 直接操作数据库
- 返回模型或错误

### 6. 中间件层 (`src/middleware/`)
**职责：全局逻辑处理**
- `ErrorRecoveryMiddleware`：全局错误捕获，处理所有panic
- `ResponseMiddleware`：统一响应格式化，将VO封装为统一响应格式

## 完整数据流

### 请求流程（以Create为例）：

```
HTTP POST /api/v1/users
    ↓
1. DTO层 - 参数校验
   ShouldBindJSON(&dto.UserCreateRequest)
   ├─ 参数错误 → panic → ErrorRecoveryMiddleware → 返回错误响应
   └─ 参数正确 → 继续
    ↓
2. 控制器层 - DTO转模型
   req.ToModel() → *models.User
    ↓
3. 控制器层 - 调用服务
   service.Create(user)
    ↓
4. 服务层 - 业务逻辑（当前无额外逻辑）
    ↓
5. 服务层 - 调用ORM
   query.User.Create(user)
   ├─ 错误 → panic → ErrorRecoveryMiddleware → 返回错误响应
   └─ 成功 → 返回模型（user已包含ID等数据库生成字段）
    ↓
6. 控制器层 - 模型转VO
   vo.FromModel(user) → *vo.UserVO
    ↓
7. 控制器层 - 设置响应
   middleware.SetData(ctx, vo)
   middleware.SetCode(ctx, 201)
    ↓
8. 中间件层 - 响应封装
   ResponseMiddleware 读取 ctx 中的 data
   封装为统一格式：{code: 201, data: UserVO}
    ↓
HTTP 响应
```

### 响应流程（以GetOne为例）：

```
HTTP GET /api/v1/users/:id
    ↓
1. 控制器层 - 解析参数
   strconv.ParseUint(id)
   ├─ 参数错误 → panic → ErrorRecoveryMiddleware
   └─ 参数正确 → 继续
    ↓
2. 控制器层 - 调用服务
   service.GetOne(query.User.ID.Eq(id))
    ↓
3. 服务层 - 调用ORM
   query.User.Where(...).First()
   ├─ 记录不存在 → panic → ErrorRecoveryMiddleware → 404错误
   └─ 成功 → 返回模型
    ↓
4. 控制器层 - 模型转VO
   vo.FromModel(user) → *vo.UserVO
    ↓
5. 控制器层 - 设置响应
   middleware.SetData(ctx, vo)
    ↓
6. 中间件层 - 响应封装
   ResponseMiddleware 封装为：{code: 200, data: UserVO}
    ↓
HTTP 响应
```

## 关键设计原则

1. **职责分离**
   - DTO：只负责入参
   - VO：只负责出参
   - 控制器：负责转换和调用
   - 服务：负责业务逻辑
   - ORM：负责数据操作

2. **错误处理**
   - 所有错误通过 `panic` 抛出
   - 全局 `ErrorRecoveryMiddleware` 统一捕获
   - 自动识别错误类型（如 `gorm.ErrRecordNotFound` → 404）

3. **数据转换**
   - 入参：HTTP → DTO → 模型
   - 出参：模型 → VO → HTTP
   - 转换逻辑集中在DTO和VO层

4. **类型安全**
   - 使用GORM Gen生成类型安全的查询
   - DTO/VO与模型分离，避免暴露敏感字段

