# 更改用户信息
## 需求分析
用户处于登录状态时，可以更改用户名、地址和签名这三个字段的信息
## 技术方案
### 前端传入参数
```json
{
  "username": "",
  "address": "",
  "signature": ""
}
```
### 后端存储结构
新建一个请求结构体存储这三个字段，另外还需要额外UserID字段对用户进行标识，但是并不需要序列化和反序列化
```go
package request
type UserChangeInfo struct {
	UserID    uint   `json:"-"`
	UserName  string `json:"user_name" binding:"required,max=20"`
	Address   string `json:"address" binding:"max=200"`
	Signature string `json:"signature" binding:"max=320"`
}
```
### 更改逻辑
只需要从数据库中将对应UserID的数据取出来，并用req中的字段进行更新，最终再存入数据库即可