# go_sso
## 架构
六边形架构（Hexagonal Architecture）
接口：定义应用与外部世界交互的接口，通常有两类：
输入端口（用于接收请求，如处理业务逻辑） 输出端口（用于与外部系统进行交互，如数据库操作）
适配器：负责实现端口接口的具体方式，适配外部系统的实现。适配器可以是数据库连接、API

优点：各个组件可以独立执行对应的业务逻辑，扩展性高

##
项目结构


## api
/login 
+ 请求体
```
{
    "userID": string,
    "password": string
}
```
+ 返回体
```
{
    "access_token": string,
    "refresh_token": string
}
```

/verify
+ 请求体
```
{
    "access_token": string
}
```
+ 返回体
```
{}
```

/refresh
+ 请求体
```
{
    “access_token”: string,
    "refresh_token": strinf
}
```
+ 返回体
```
{
  "access_token": string
}
```
/Update/[field]  

tips:field为要修改的字段  
包括"Email" "Name" "PhoneNumber"

+ 请求体
```
{
    "userId": string
    "[field]": string
}
```
+返回体
```
{
{
    "userId": "20233802086",
    "name": string,
    "grade": string,
    "majorClass": string,
    "email": string,
    "phoneNumber": string,
    "status": string
}
}
```
  
  