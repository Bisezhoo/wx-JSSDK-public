package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"crypto/sha1"
	"encoding/hex"
)

var client = http.Client{
	Timeout: 10 * time.Second,
}

var lastTime int
var jsapiTicket string
var accessToken string
var signature string

func main() {
	http.HandleFunc("/getSignature", getSignature)
	http.ListenAndServe("127.0.0.1:5000", nil)
}

//SHA1加密
func SHA1(s string) string {
	o := sha1.New()
	o.Write([]byte(s))
	return hex.EncodeToString(o.Sum(nil))
}

// handler函数
func getSignature(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	response.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	var responseBodyStr string

	keys := request.URL.Query()
	queryTimeStr := keys.Get("time")
	queryTime,_ :=  strconv.Atoi(queryTimeStr)
	nonceStr := keys.Get("nonceStr")
	url := keys.Get("url")

	//判断参数是否为空
	if 0 == strings.Compare("",queryTimeStr) ||  0 == strings.Compare("",nonceStr) ||  0 == strings.Compare("",url){
		return
	}

	//判断间隔，微信key过期时间为7200秒
	intervalTime := queryTime - lastTime

	if 7200 < intervalTime || lastTime == 0  {
		//重置上次更新时间
		lastTime = queryTime

		//获取accessToken
		requestAccessTokenUrl := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=wx0000000000000000&secret=00000000000000000000000000000000"
		accessTokenResp, _ := client.Get(requestAccessTokenUrl)

		//构建接收返回值的map
		formData := make(map[string]string )
		// 调用json包的解析，解析请求body
		json.NewDecoder(accessTokenResp.Body).Decode(&formData)
		//取得accessToken
		accessToken = formData["access_token"]
		fmt.Println(time.Now().Format("2006-01-02 15:04:05") ," - 获取新accessToken：" ,accessToken)

		//获取jsapiTicket
		requestJsapiTicketUrl := "https://api.weixin.qq.com/cgi-bin/ticket/getticket?type=jsapi&access_token=" + accessToken
		jsapiTicketResp, _ := client.Get(requestJsapiTicketUrl)
		// 构建接收返回值的map
		formData = make(map[string]string )
		// 调用json包的解析，解析请求body
		json.NewDecoder(jsapiTicketResp.Body).Decode(&formData)
		//取得jsapiTicket
		jsapiTicket = formData["ticket"]
		fmt.Println(time.Now().Format("2006-01-02 15:04:05") ," - 获取新jsapiTicket：" ,jsapiTicket)

	}
	var shaStrBuilder  strings.Builder
	shaStrBuilder.WriteString("jsapi_ticket=")
	shaStrBuilder.WriteString(jsapiTicket)
	shaStrBuilder.WriteString("&noncestr=")
	shaStrBuilder.WriteString(nonceStr)
	shaStrBuilder.WriteString("&timestamp=")
	shaStrBuilder.WriteString(strconv.Itoa(queryTime))
	shaStrBuilder.WriteString("&url=")
	shaStrBuilder.WriteString(url)
	signature = SHA1(shaStrBuilder.String())
	responseBodyStr = signature
	fmt.Println(time.Now().Format("2006-01-02 15:04:05") ," - 获取了signature" ,signature)
	//返回signature
	response.Write([]byte(fmt.Sprintf("%s" ,responseBodyStr)))
}