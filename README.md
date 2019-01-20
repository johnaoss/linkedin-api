# GoLinkedinAPI [![GoDoc](https://godoc.org/github.com/johnaoss/golinkedinapi?status.svg)](https://godoc.org/github.com/johnaoss/golinkedinapi)

This project functions as a pure Go interface for Linkedin's REST API. Currently, this provides a lightweight, documented interface to get a user's posts & profile data into a native marshallable Go struct. Requests are made secure by using OAuth2.0 authenticated requests to LinkedIn's servers.

This was my first project written in Go, and so I'd love to hear your thoughts!

This currently only supports GET requests, and soon to be POST requests.

## Installation

```bash
go get -t github.com/johnaoss/golinkedinapi
```

That's it!

## Examples

I haven't tested this with other routing packages, but this does indeed work for any program using gorilla/mux for routing purposes.

### Setting up the configurations

```go

import api "github.com/johnaoss/golinkedinapi"

func main() {
    permissions := []string{"r_basicprofile"}
    clientID := "myID"
    clientSecret := "hush"
    redirectURL := "https://example.com/totallyvalidauth"
    api.InitConfig(permissions, clientID, clientSecret, redirectURL)
}

```

### Getting a user's login url (includes state)

```go
import api "github.com/johnaoss/golinkedinapi"

func loginHandler(w http.ResponseWriter, r *http.Request) {
    login := api.GetLoginURL(w,r)
    html := "Your login is <a href=\"" + login + "\">Login here!</a>"
    w.Write([]byte(html))
}
```

### Getting a user's data

```go
import api "github.com/johnaoss/golinkedinapi"

// this handles the authorized redirect URL as specified in the Linkedin developer console
func authHandler(w http.ResponseWriter, r *http.Request) {
    userData := api.GetProfileData(w,r)
    html := "Your name is is " userData.FirstName + " " + userData.LastName
    w.Write([]byte(html))
}
```

### Sharing a comment (UNTESTED)

```go
import api "github.com/johnaoss/golinkedinapi"

func sharePost(w http.ResponseWriter, r *http.Request) {
    vis := VisibilityStruct{Code: "anyone"}
    post := &api.Post{Visibility: vis, Comment: "This is a comment"}
    resp, err := api.SharePost(post,w,r)
    if err != nil {
        w.Write([]byte("Something went wrong!"))
    } else {
        w.Write([]byte("Your post was successfully shared!"))
    }
}
```

## Limitations

This currently services only the non-partnered content as LinkedIn does not have any means of open-source developers to reliably acquire this data. Aside from that, if there are any bugs please feel free to get in touch with either an issue or an e-mail!

## TODO

1. Proper testing of functions that require a router.
1. POST features
1. Full Linkedin Partner features (if possible)
1. More info on `current-share`
