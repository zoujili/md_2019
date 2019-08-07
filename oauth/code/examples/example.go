package main

import (
	"IdentityService/config"
	"IdentityService/pkg/oauth2"
	"IdentityService/pkg/oauth2/providers/hpid"
	"context"
	"fmt"
	"github.com/gorilla/pat"
	"html/template"
	"log"
	"net/http"
	"sort"
)

type Providers map[string]oauth2.IdentityProvider
var providers = Providers{}

func UseProviders(viders ...oauth2.IdentityProvider) {
	for _, provider := range viders {
		providers[provider.Name()] = provider
	}
}

func main() {

	provider:=hpid.New(config.GetConfig())
	provider.SetRedirectURL("https://k9-dev.fitstation4me.com/hpid/callback")
	UseProviders(
		provider,
	)

	m := make(map[string]string)
	m["hpid"] = "HPID"
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	providerIndex := &ProviderIndex{Providers: keys, ProvidersMap: m}
	p := pat.New()
	p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
		//case this local callback not register in HPID ,so copy the  response code here
		token, err := provider.Exchange(context.Background(),  "AWc7nCRAeesaa3d7C6NSzroEywK1AAAAAAAAAABLw3sRTVXGXwpYECcxUHGCeVV6ILwVWaMLHbqYTuAiyKyPdOjltkCr_ldOwZfQQUX-gNIpqDKBj-8q7sx9JHaHL70c_uRCx3czgedjjOczYw")
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}
		user,err:=provider.Profile(context.Background(),*token)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(res, user)
	})

	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		url:= provider.AuthCodeURL("test")
		http.Redirect(res, req, url, http.StatusTemporaryRedirect)
	})

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.New("foo").Parse(indexTemplate)
		t.Execute(res, providerIndex)
	})
	log.Println("listening on localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", p))
}


type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

var indexTemplate = `{{range $key,$value:=.Providers}}
    <p><a href="/auth/{{$value}}">Log in with {{index $.ProvidersMap $value}}</a></p>
{{end}}`

var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
`
