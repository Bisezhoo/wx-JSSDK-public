# wx-JSSDK-public是一个用GO语言写的获取微信JSAPI signature的轻后端
由于go足够轻，只需要更改代码内的appID与secret，即可完成配置。
执行编译
--SET CGO_ENABLED=0
--SET GOOS=linux 
--SET GOARCH=amd64
--go build main.go
丢到服务器上完事
