package code


import (
	"context"
	"fmt"
	"github.azc.ext.hp.com/fitstation/go-service-utils/errors"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)


type IdentityProvider interface {
	Name() string
	SetName(name string)
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	Profile(ctx context.Context, token oauth2.Token) (*UserProfile, error)
	SetRedirectURL(redirectURL string)
	SetScopes(scopes []string)
}


type UserProfile struct {
	RawData           map[string]interface{}
	Provider          string
	Email             string
	Name              string
	FirstName         string
	LastName          string
	NickName          string
	Description       string
	UserID            string
	AvatarURL         string
	Location          string
	EmailVerified     bool
}


func FetchUserProfilePayload(ctx context.Context, profileURL string ,token oauth2.Token) ([]byte, error) {
	req, err := http.NewRequest("GET",profileURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("unable to get user profile")
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("unable to read user profile payload")
		return body, err
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		fmt.Println("unable to get user profile")
		return nil, errors.New("unable to get user profile")
	}
	return body, nil
}