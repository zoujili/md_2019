package hpid

import (
	"context"
	"encoding/json"
	"fmt"
	"github.azc.ext.hp.com/fitstation/go-service-utils/errors"
	"golang.org/x/oauth2"
	oauthclient "IdentityService/pkg/oauth2"
)


const (
	ScopeUserRead string = "user.profile.read"
	ScopeOfficeAccess string = "offline_access"
)

type Configuration interface {
	GetHPIDOAuthProviderClientID() string
	GetHPIDOAuthProviderClientSecret() string
	GetHPIDOAuthProviderEndpointAuth() string
	GetHPIDOAuthProviderEndpointToken() string
	GetHPIDOAuthProviderEndpointUserInfo() string
}

type Provider struct {
	oauth2.Config
	name  string
	ProfileURL string
}

type IDPUser struct {
	ID        string `json:"id"`
	Emails    []Email `json:"emails"`
	Name      Name `json:"name"`
	UserName string `json:"userName"`
}

type Email struct {
	Primary  bool  `json:"primary"`
	Value    string  `json:"value"`
	Verified  bool   `json:"verified"`
}

type Name struct {
	FamilyName string  `json:"familyName"`
	GivenName  string  `json:"givenName"`
	MiddleName string  `json:"middleName"`
}

func New(config  Configuration,scopes ...string) *Provider {
	p := &Provider{}
	p.ProfileURL = config.GetHPIDOAuthProviderEndpointUserInfo()
	p.ClientID = config.GetHPIDOAuthProviderClientID()
	p.ClientSecret = config.GetHPIDOAuthProviderClientSecret()
	p.Endpoint = oauth2.Endpoint{AuthURL: config.GetHPIDOAuthProviderEndpointAuth(), TokenURL: config.GetHPIDOAuthProviderEndpointToken()}
	p.name = "hpid"
	if len(scopes) > 0 {
		for _, scope := range scopes {
			p.Scopes = append(p.Scopes, scope)
		}
	} else {
		p.Scopes = []string{ScopeUserRead}
	}
	return p
}


func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) SetName(name string) {
	p.name = name
}

func (p *Provider) SetRedirectURL(redirectURL string) {
	p.RedirectURL = redirectURL

}

func (p *Provider) SetScopes(scopes []string) {
	p.Scopes = scopes
}


func (p *Provider) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	return p.Config.AuthCodeURL(state,oauth2.SetAuthURLParam("max_age","0"))

}
func (p *Provider) Profile(ctx context.Context, token oauth2.Token) (*oauthclient.UserProfile, error) {
	body, err := oauthclient.FetchUserProfilePayload(ctx,p.ProfileURL ,token)
	if err != nil {
		return nil, err
	}
	var u oauthclient.UserProfile
	var idpUser IDPUser
	err = json.Unmarshal(body, &idpUser)
	if err != nil {
		return nil, err
	}
	email :=""
	emailVerified := false
	for _,v:=range idpUser.Emails{
		if v.Primary == true{
			email = v.Value
			emailVerified = v.Verified
			break
		}
	}
	if email == ""{
		fmt.Println("user don't have primary email")
		return nil, errors.New("unable to get user profile")
	}
	u = oauthclient.UserProfile{
		Email:         email,
		EmailVerified: emailVerified,
		Name:          idpUser.UserName,
		FirstName:     idpUser.Name.FamilyName,
		LastName:      idpUser.Name.GivenName,
	}
	return &u, nil
}



