package golinkedinapi

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: Automate this

var (
	vis                 = VisibilityStruct{Code: "anyone"}
	comment             = "This is a comment"
	longTitle           = generateText(201)
	longDescription     = generateText(257)
	longComment         = generateText(701)
	validURL            = "https://google.com"
	invalidURL          = "This isn't a URL"
	basicContent        = ContentStruct{Title: "A Title"}
	emptyPost           = &Post{}
	basicPost           = &Post{Visibility: vis, Comment: comment}
	invalidPostTitle    = &Post{Visibility: vis, Comment: comment, Content: ContentStruct{Title: longTitle}}
	invalidPostDesc     = &Post{Visibility: vis, Comment: comment, Content: ContentStruct{Description: longDescription}}
	invalidPostComment  = &Post{Visibility: vis, Comment: longComment}
	invalidSubmittedURL = &Post{Visibility: vis, Comment: comment, Content: ContentStruct{SubmittedURL: invalidURL}}
	invalidImageURL     = &Post{Visibility: vis, Comment: comment, Content: ContentStruct{SubmittedImageURL: invalidURL}}
)

// generateText is a helper to generate descriptions over the char limit
// for use in testing IsValidPost()
func generateText(width int) string {
	b := make([]byte, width)
	rand.Read(b)
	return string(b)
}

func TestIsValidPost(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(false, isValidPost(emptyPost))
	assert.Equal(true, isValidPost(basicPost))
	assert.NotEqual(false, isValidPost(basicPost))
	assert.Equal(false, isValidPost(invalidPostTitle))
	assert.Equal(false, isValidPost(invalidPostDesc))
	assert.Equal(false, isValidPost(invalidPostComment))
	assert.Equal(false, isValidPost(invalidImageURL))
	assert.Equal(false, isValidPost(invalidSubmittedURL))
}

func TestPostToJSON(t *testing.T) {
	assert := assert.New(t)

	assert.NotEqual(nil, postToJSON(emptyPost))
	assert.NotEqual(nil, postToJSON(basicPost))
	// Figure out how to make json.Marshal(post) error out
}
