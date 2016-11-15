package epiJar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	cookiejar "net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

type (
	EpiJar struct {
		email    string
		password string
		jar      *cookiejar.Jar
		client   *http.Client
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

func (ej EpiJar) GetHTTP() *http.Client {
	return ej.client
}

func (ej EpiJar) Auth() {
	// Request Intra for office365 URL
	req, _ := http.NewRequest("GET", "https://intra.epitech.eu", nil)
	resp, err := ej.client.Do(req)
	if err != nil {
		panic(nil)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	re := regexp.MustCompile("(https://login.microsoftonline([^\"]+))")
	m := re.FindStringSubmatch(string(body))
	url := "https://login.microsoftonline" + string(m[2])
	log.Println("Office365 URL: " + url)

	// Request URL for CONTEXT
	req, _ = http.NewRequest("GET", url, nil)
	resp, err = ej.client.Do(req)
	if err != nil {
		panic(nil)
	}
	body, _ = ioutil.ReadAll(resp.Body)
	re = regexp.MustCompile("Constants.CONTEXT = '([^']+)")
	m = re.FindStringSubmatch(string(body))
	urlRealm := "https://login.microsoftonline.com/common/userrealm/?user=" + ej.email + "&api-version=2.1&stsRequest=" + string(m[1]) + "&checkForMicrosoftAccount=true"
	log.Println("urlRealm: " + urlRealm)

	// Get AuthURL
	req, _ = http.NewRequest("GET", urlRealm, nil)
	resp, err = ej.client.Do(req)
	if err != nil {
		panic(nil)
	}
	body, _ = ioutil.ReadAll(resp.Body)
	var dat map[string]interface{}
	if err := json.Unmarshal(body, &dat); err != nil {
		panic(err)
	}
	authURL, _:= dat["AuthURL"].(string)
	log.Println("AuthURL: " + authURL)

	// Get Epitech Login Page
	req, _ = http.NewRequest("GET", authURL, nil)
	resp, err = ej.client.Do(req)
	if err != nil {
		panic(nil)
	}
	body, _ = ioutil.ReadAll(resp.Body)
	epitechLoginPage := string(body)

	fmt.Println(epitechLoginPage)

	resp.Body.Close()
}

func (ej EpiJar) TMP() {
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
