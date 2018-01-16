package golinkedinapi

import "net/url"

// This is for POST requests to the Linkedin API.
// That includes features such as...
// * Share on Linkedin
// * Manage Company Pages
// * Other Partner Features...(TBD if possible)

// Post represents a post that one shares on a Linkedin profile,
// not a POST request. To validate posts, call isValidPost() on one.
type Post struct {
	Content    ContentStruct    `json:"content"`
	Comment    string           `json:"comment"`
	Visibility VisibilityStruct `json:"visibility"`
}

// ContentStruct represents a collection of fields describing the shared content.
type ContentStruct struct {
	Title             string `json:"title"`
	Description       string `json:"description"`
	SubmittedURL      string `json:"sumbitted-url"`
	SubmittedImageURL string `json:"submitted-image-url"`
}

// VisibilityStruct represents the visibility information about the shared post.
type VisibilityStruct struct {
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
	_, err := url.Parse(post.Content.SubmittedURL)
	if err != nil {
		return false
	}
	_, err = url.Parse(post.Content.SubmittedImageURL)
	if err != nil {
		return false
	}
	return false
}
