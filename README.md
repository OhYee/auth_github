# QQ互联基本接口封装

[![Sync to Gitee](https://github.com/OhYee/auth_github/workflows/Sync%20to%20Gitee/badge.svg)](https://gitee.com/OhYee/auth_github) [![version](https://img.shields.io/github/v/tag/OhYee/auth_github)](https://github.com/OhYee/auth_github/tags)

[Github Auth APP](https://github.com/settings/developers)

Jump to the login page (Using `LoginPage()` get the login URL)

Then use `Auth()` get the ` access_token`, and use `access_token` get the user info

```go
conn := qq.New("Your secret id", "Your secret key", "Your redirect uri")

token, err := conn.Auth(code, state)
if err != nil {
    return
}
output.Debug("%+v", token)

info, err := conn.Info(token)
if err != nil {
    return
}
```