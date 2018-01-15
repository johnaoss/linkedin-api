package golinkedinapi_test

import (
	"testing"

	api "github.com/johnaoss/golinkedinapi"
)

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
