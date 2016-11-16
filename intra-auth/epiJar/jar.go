package epiJar

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	cookiejar "net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"html"
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
	Office365url := "https://login.microsoftonline" + string(m[2])
	log.Println("Office365 URL: " + Office365url)

	// Request URL for CONTEXT
	req, _ = http.NewRequest("GET", Office365url, nil)
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
	re = regexp.MustCompile("action=\"/adfs([^\"]+)")
	m = re.FindStringSubmatch(epitechLoginPage)
	loginUrl := "https://sts.epitech.eu/adfs" + string(m[1])
	log.Println("loginUrl: " + loginUrl)

	//Check Password
	postData := url.Values{}
	postData.Set("UserName", ej.email)
	postData.Set("Password", ej.password)
	postData.Set("Kmsi", "true")
	postData.Set("AuthMethod", "FormsAuthentication")
	req, _ = http.NewRequest("POST", loginUrl, strings.NewReader(postData.Encode()))
	resp, err = ej.client.Do(req)
	if err != nil {
		panic(nil)
	}
	body, _ = ioutil.ReadAll(resp.Body)
	re = regexp.MustCompile("action=\"([^\"]+)")
	m = re.FindStringSubmatch(string(body))
	microsoftLoginUrl := string(m[1])

	re = regexp.MustCompile("name=\"wa\" value=\"([^\"]+)")
	m = re.FindStringSubmatch(string(body))
	wa := m[1]

	re = regexp.MustCompile("name=\"wresult\" value=\"([^\"]+)")
	m = re.FindStringSubmatch(string(body))
	wresult := m[1]

	re = regexp.MustCompile("name=\"wctx\" value=\"([^\"]+)")
	m = re.FindStringSubmatch(string(body))
	wctx := m[1]

	postData = url.Values{}
	postData.Set("wa", wa)
	postData.Set("wresult", html.UnescapeString(wresult))
	postData.Set("wctx", html.UnescapeString(wctx))
	req, _ = http.NewRequest("POST", microsoftLoginUrl, strings.NewReader(postData.Encode()))
	resp, err = ej.client.Do(req)
	if err != nil {
		panic(nil)
	}
	body, _ = ioutil.ReadAll(resp.Body)

	// Check Grant
	resp.Body.Close()
}

func (ej EpiJar) GetClient() *http.Client {
	return ej.client
}
