package epiJar

import (
	"fmt"
	"io/ioutil"
	"net/http"
	cookiejar "net/http/cookiejar"
	"net/url"
	"strings"
)

type (
	EpiJar struct {
		email string
		password string
		jar *cookiejar.Jar
		client *http.Client
	}
)

func New(email string, password string) *EpiJar {
	tmp := EpiJar{}
	tmp.email = email
	tmp.password = password
	tmp.jar, _ = cookiejar.New(nil)
	tmp.client = &http.Client{
    Jar: tmp.jar,
  }
	return &tmp
}

func (ej EpiJar) Auth() {
	postData := url.Values{}
	postData.Set("keyword", "尹相杰")
	postData.Set("smblog", "搜微博")
	req, _ := http.NewRequest("POST", "http://weibo.cn/search/?vt=4", strings.NewReader(postData.Encode()))
	resp, err := ej.client.Do(req)
	if err != nil {
		panic(nil)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println(string(body))
}

func GetJar() {

}
