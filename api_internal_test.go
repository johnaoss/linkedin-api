package golinkedinapi

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestParseJSON(t *testing.T) {
	badData := []string{"This is not JSON", "{{}", "{no}", "{name: \"lol\"}"}
	goodData := []string{"{}", "{\"first_name\": \"John\"}"}
	exampleResponse, err := ioutil.ReadFile("exampleresponse.json")
	if err != nil {
		t.Errorf("Could not find exampleresponse.json")
	}
	goodDataResponses := []*LinkedinProfile{&LinkedinProfile{}, &LinkedinProfile{FirstName: "John"}}
	for _, elem := range badData {
		get, err := parseJSON(elem)
		if err == nil {
			t.Errorf("ParseJSON should have errored, instead given: %v", get)
		}
	}
	for num, elem := range goodData {
		get, err := parseJSON(elem)
		if err != nil {
			t.Errorf("ParseJSON threw an error on valid data: %v", err)
		}
		if !reflect.DeepEqual(goodDataResponses[num], get) {
			t.Errorf("ParseJSON errored, wanted %v, given %v", goodDataResponses[num], get)
		}
	}

	_, err = parseJSON(string(exampleResponse))
	if err != nil {
		t.Errorf("ParseJSON threw an error on a valid response, given: %s\n:", err.Error())
	}
}

func TestGetSessionValue(t *testing.T) {
	exampleMap := map[string]interface{}{
		"test": "test",
	}
	value := getSessionValue(exampleMap["test"])
	if value != "test" {
		t.Errorf("Failed to get session string, wanted %v, got %v\n", "test", value)
	}
	value = getSessionValue(nil)
	if value != "" {
		t.Errorf("Failed to get value of nil, wanted \"\", got %s\n", value)
	}
}

func TestGenerateState(t *testing.T) {
	s := generateState()
	if reflect.TypeOf(s).Name() != "string" {
		t.Errorf("GenerateState needed to return string, instead got %s\n", reflect.TypeOf(s).Name())
	}
}

func TestValidState(t *testing.T) {
	// Test empty state comparison
	request := new(http.Request)
	get := validState(request)
	if get == false {
		t.Errorf("validState should have returned true, instead got %v\n", get)
	}
	session, _ := store.Get(request, "golinkedinapi")
	session.Values["state"] = "Example bad state"
	get = validState(request)
	if get != false {
		t.Errorf("validState should have returned false, instead got %v\n", get)
	}

}
