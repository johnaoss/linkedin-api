package linkedin

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
	if len(s) != 32 {
		t.Errorf("Returned improper size")
	}
}

func TestValidState(t *testing.T) {
	// Test empty state comparison
	request := new(http.Request)
	if ok := validState(request); !ok {
		t.Errorf("validState should have returned true, instead got %t", ok)
	}

	session, err := store.Get(request, "golinkedinapi")
	if err != nil {
		t.Errorf("received error while receiving state: %v", err)
	}
	session.Values["state"] = "Example bad state"
	if ok := validState(request); ok {
		t.Errorf("validState should have returned false, instead got %t", ok)
	}
}
