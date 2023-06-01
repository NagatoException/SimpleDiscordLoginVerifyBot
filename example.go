//go:build ignore
// +build ignore
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type LoginResponse struct {
	Success bool `json:"success"`
	UUID string `json:"uuid"`
}
func main() {
	fmt.Println("输入你的 Discord 用户名")
	var input string
	fmt.Scanln(&input)
	username := strings.Split(input, "#")[0]
	tag := strings.Split(input, "#")[1]

	loginURL := "http://localhost/new?username=" + username + "&tag=" + tag
	req, _ := http.Get(loginURL)
	defer req.Body.Close()
	resp,_ := io.ReadAll(req.Body)
	var response LoginResponse
	json.Unmarshal(resp, &response)
	if response.Success == true {
		fmt.Println("成功发送登录请求")
		for i := 0; i < 10; i++ {
			verifyURL := "http://localhost/verify?uuid=" + response.UUID
			verifyReq, _ := http.Get(verifyURL)
			defer verifyReq.Body.Close()
			resp, _ = io.ReadAll(verifyReq.Body)
			var verifyResponse LoginResponse
			err := json.Unmarshal(resp, &verifyResponse)
			if err != nil {
				panic(err)
			}

			if verifyResponse.Success == true {
				fmt.Println("成功登录")
				return
			}

			time.Sleep(5 * time.Second)
		}
	}


	fmt.Println("登陆失败")
}
