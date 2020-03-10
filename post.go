package linkedin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// This is for POST requests to the Linkedin API.
// That includes features such as...
// * Share on Linkedin
// * Manage Company Pages
// * Other Partner Features...(TBD if possible)

const (
	postURL = "https://api.linkedin.com/v1/people/~/shares?format=json"
)

// ContentStruct represents a collection of fields describing the shared content.
type Content struct {
	Title             string `json:"title,omitempty"`
	Description       string `json:"description,omitempty"`
	SubmittedURL      string `json:"sumbitted-url,omitempty"`
	SubmittedImageURL string `json:"submitted-image-url,omitempty"`
}

// VisibilityStruct represents the visibility information about the shared post.
type Visibility struct {
	// Code is the only field specified in the documentation.
	// must be either "anyone" or "connections-only"
	Code string `json:"code"`
}

// Post represents a post that one shares on a Linkedin profile,
// not a POST request. To validate posts, call isValidPost() on one.
type Post struct {
	Content    Content    `json:"content,omitempty"`
	Comment    string     `json:"comment"`
	Visibility Visibility `json:"visibility"`
}

// Validate checks to see if this post is valid or not.
func (p *Post) Validate() error {
	if p.Visibility.Code != "anyone" && p.Visibility.Code != "connections-only" {
		return fmt.Errorf("invalid visibility codes")
	}

	if len(p.Comment) == 0 {
		return fmt.Errorf("comment can not be empty")
	}

	// Check length limits
	if len(p.Content.Title) > 200 {
		return fmt.Errorf("content title exceeds content length of 200 characters")
	}
	if len(p.Content.Description) > 256 {
		return fmt.Errorf("description exceeds content length of 256 characters")
	}
	if len(p.Comment) > 700 {
		return fmt.Errorf("comment exceeds content length of 700 characters")
	}

	// Check if URLs are FQDN
	if p.Content.SubmittedURL != "" {
		_, err := url.ParseRequestURI(p.Content.SubmittedURL)
		if err != nil {
			return fmt.Errorf("failed to validate submitted url: %w", err)
		}
	}
	if p.Content.SubmittedImageURL != "" {
		_, err := url.ParseRequestURI(p.Content.SubmittedImageURL)
		if err != nil {
			return fmt.Errorf("failed to validate submitted image url: %w", err)
		}
	}

	return nil
}

// SharePost is the primary method of sharing a post.
func SharePost(post *Post, w http.ResponseWriter, r *http.Request) (*http.Response, error) {
	if !validState(r) {
		err := fmt.Errorf("state comparison failed")
		return &http.Response{}, err
	}
	// Authenticate
	params := r.URL.Query()
	tok, err := authConf.Exchange(context.Background(), params.Get("code"))
	if err != nil {
		return &http.Response{}, err
	}
	client := authConf.Client(context.Background(), tok)
	// Upload Data
	data, err := json.Marshal(post)
	if err != nil {
		return &http.Response{}, err
	}

	resp, err := client.Post(postURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return &http.Response{}, err
	}
	if resp.StatusCode == 400 {
		return resp, fmt.Errorf("recieved 404: %v", resp)
	}
	if resp.StatusCode == 201 {
		return resp, nil
	}
	err = fmt.Errorf("post unsuccessful, response given: %v", resp)
	return &http.Response{}, err
}
