package main

import (
	jar "github.com/tsauzeau/authIntra/intra-auth/epiJar"
	"net/http"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	jar := jar.New(os.Getenv("EMAIL"), os.Getenv("PASS"))
	jar.Auth()
	req, _ := http.NewRequest("GET", "https://intra.epitech.eu/?format=json", nil)
	resp, err := jar.GetClient().Do(req)
	if err != nil {
		panic(nil)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
