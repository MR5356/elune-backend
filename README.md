## 公共返回参数

```json
{
  "code": 0,
  "message": "success",
  "data": Object
}
```

## 用户管理
* `/user/login` `POST` 登录 

| 接口                  | 请求方式   | 请求参数                         | 返回格式(公共返回参数中的data字段) | 介绍      |
|---------------------|--------|------------------------------|----------------------|---------|
| /user/login         | POST   | {username: "", password: ""} | {token: ""}          | 用户登录    |
| /user/logout        | DELETE | 无                            | null                 | 用户登出    | 
| /user/token/refresh | GET    | 无                            | {need: true}         | token刷新 |
| /user/token/refresh | PUT    | 无                            | {token: ""}          | token刷新 |
| /user/info          | GET    | 无                            | {username: ""}       | 用户信息    |