package golinkedinapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

// This is for POST requests to the Linkedin API.
// That includes features such as...
// * Share on Linkedin
// * Manage Company Pages
// * Other Partner Features...(TBD if possible)

const (
	postURL = "https://api.linkedin.com/v1/people/~/shares?format=json"
)

// Post represents a post that one shares on a Linkedin profile,
// not a POST request. To validate posts, call isValidPost() on one.
type Post struct {
	Content    ContentStruct    `json:"content,omitempty"`
	Comment    string           `json:"comment"`
	Visibility VisibilityStruct `json:"visibility"`
}

// ContentStruct represents a collection of fields describing the shared content.
type ContentStruct struct {
	Title             string `json:"title,omitempty"`
	Description       string `json:"description,omitempty"`
	SubmittedURL      string `json:"sumbitted-url,omitempty"`
	SubmittedImageURL string `json:"submitted-image-url,omitempty"`
}

// VisibilityStruct represents the visibility information about the shared post.
type VisibilityStruct struct {
	// Code is the only field specified in the documentation.
	// must be either "anyone" or "connections-only"
	Code string `json:"code"`
}

// isValidPost checks to see if a post is valid or not.
// This is all done in accordance with https://developer.linkedin.com/docs/share-on-linkedin
// and is valid as of 2018-01-16
func isValidPost(post *Post) bool {
	// Validate field essentials.
	if (post.Visibility.Code != "anyone" && post.Visibility.Code != "connections-only") || len(post.Comment) == 0 {
		return false
	}
	// Check length limits
	if len(post.Content.Title) > 200 || len(post.Content.Description) > 256 || len(post.Comment) > 700 {
		return false
	}
	// Check if URLs are FQDN
	if post.Content.SubmittedURL != "" {
		_, err := url.ParseRequestURI(post.Content.SubmittedURL)
		if err != nil {
			return false
		}
	}
	if post.Content.SubmittedImageURL != "" {
		_, err := url.ParseRequestURI(post.Content.SubmittedImageURL)
		if err != nil {
			return false
		}
	}
	return true
}

// postToJSON converts a post struct into a JSON object in order to send to
// the API's endpoint.
func postToJSON(post *Post) []byte {
	data, err := json.Marshal(post)
	if err != nil {
		return nil
	}
	return data
}

// SharePost is the primary method of sharing a post.
func SharePost(post *Post, w http.ResponseWriter, r *http.Request) (*http.Response, error) {
	if validState(r) == false {
		err := fmt.Errorf("state comparison failed")
		return &http.Response{}, err
	}
	// Authenticate
	params := r.URL.Query()
	tok, err := authConf.Exchange(oauth2.NoContext, params.Get("code"))
	if err != nil {
		return &http.Response{}, err
	}
	client := authConf.Client(oauth2.NoContext, tok)
	// Upload Data
	data := postToJSON(post)
	resp, err := client.Post(postURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return &http.Response{}, err
	}
	if resp.StatusCode == 400 {
		err = fmt.Errorf("cannot post duplicate content")
	}
	if resp.StatusCode == 201 {
		return resp, nil
	}
	err = fmt.Errorf("post unsuccessful, response given: %v", resp)
	return &http.Response{}, err
}
