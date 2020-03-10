package linkedin

import (
	"crypto/rand"
	"testing"
)

// TODO: Automate this

// generateText is a helper to generate descriptions over the char limit
// for use in testing Post.Validate()
func generateText(t *testing.T, width int) string {
	t.Helper()

	b := make([]byte, width)
	if _, err := rand.Read(b); err != nil {
		t.Fatalf("failed to read from random: %v", err)
	}
	return string(b)
}

func TestPost_Validate(t *testing.T) {
	vis := Visibility{Code: "anyone"}
	comment := "This is a comment"

	testcases := []struct {
		name    string
		post    *Post
		wantErr bool
	}{
		{
			name:    "Empty",
			post:    &Post{},
			wantErr: true,
		},
		{
			name: "Basic",
			post: &Post{
				Visibility: vis,
				Comment:    comment,
			},
			wantErr: false,
		},
		{
			name: "invalid title",
			post: &Post{
				Visibility: vis,
				Comment:    comment,
				Content:    Content{Title: generateText(t, 201)},
			},
			wantErr: true,
		},
		{
			name: "invalid description",
			post: &Post{
				Visibility: vis,
				Comment:    comment,
				Content:    Content{Description: generateText(t, 257)},
			},
			wantErr: true,
		},
		{
			name: "invalid comment",
			post: &Post{
				Visibility: vis,
				Comment:    generateText(t, 701),
			},
			wantErr: true,
		},
		{
			name: "invalid submitted url",
			post: &Post{
				Visibility: vis,
				Comment:    comment,
				Content:    Content{SubmittedURL: "This isn't a URL"},
			},
			wantErr: true,
		},
		{
			name: "invalid submitted image url",
			post: &Post{
				Visibility: vis,
				Comment:    comment,
				Content:    Content{SubmittedImageURL: "Neither is this"},
			},
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.post.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("should error is %t, did err = %t", tc.wantErr, err != nil)
			}
		})
	}
}
