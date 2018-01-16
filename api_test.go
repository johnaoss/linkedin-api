package golinkedinapi_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	api "github.com/johnaoss/golinkedinapi"
)

const port = ":5000"

func TestInitConfig(t *testing.T) {
	emailPermissions := []string{"r_emailaddress"}
	profilePermissions := []string{"r_basicprofile"}
	fullPermissions := []string{"r_emailaddress", "r_basicprofile"}
	invalidPermissions := []string{"This is an invalid string"}
	dummyID, dummySecret := "yes", "no"
	validURL := "https://example.com/auth"
	invalidURL := "example.com"

	// Invalid attempts
	tryConfigPanic(t, emailPermissions, dummyID, dummySecret, invalidURL)
	tryConfigPanic(t, invalidPermissions, dummyID, dummySecret, validURL)
	tryConfigPanic(t, nil, dummyID, dummySecret, validURL)

	// Valid attempts
	api.InitConfig(emailPermissions, dummyID, dummySecret, validURL)
	api.InitConfig(profilePermissions, dummyID, dummySecret, validURL)
	api.InitConfig(fullPermissions, dummyID, dummySecret, validURL)

}

// Helper function for TestInitConfig due to Golang limitations.
// Turns out you can't pass a function as an argument that has no return value.
// (i.e. testify's assert subpackage didn't work with initConfig() as a param)
func tryConfigPanic(t *testing.T, p []string, id string, secret string, url string) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("Program should have panicked")
		}
	}()
	api.InitConfig(p, id, secret, url)
}

func TestGetLoginURL(t *testing.T) {
	go initRouter()
	fullPermissions := []string{"r_emailaddress", "r_basicprofile"}
	dummyID, dummySecret := "yes", "no"
	validURL := "https://example.com/auth"
	api.InitConfig(fullPermissions, dummyID, dummySecret, validURL)
	resp, err := http.Get("http://localhost:5000")
	if err != nil {
		t.Errorf("Could not get LoginURL, given %s", err.Error())
	}
	url, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Could not read response of LoginURL, given %s", err.Error())
	}
	fmt.Println(url)

}

// Initializes a router for testing purposes.
func initRouter() {
	router := mux.NewRouter()
	router.HandleFunc("/", loginHandler)
	http.Handle("/", router)
	http.ListenAndServe(port, router)
}

// Handles a route for getLoginURL()
func loginHandler(w http.ResponseWriter, r *http.Request) {
	login := api.GetLoginURL(w, r)
	w.Write([]byte(login))
}
