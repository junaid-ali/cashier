package github

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/nsheridan/cashier/server/config"
	"github.com/stretchr/testify/assert"
)

var (
	oauthClientID     = "id"
	oauthClientSecret = "secret"
	oauthCallbackURL  = "url"
	organization      = "exampleorg"
	users             = []string{"user"}
)

func TestNew(t *testing.T) {
	a := assert.New(t)

	p, _ := New(&config.Auth{
		OauthClientID:     oauthClientID,
		OauthClientSecret: oauthClientSecret,
		ProviderOpts:      map[string]string{"organization": organization},
		UsersWhitelist:    users,
	})
	a.Equal(p.config.ClientID, oauthClientID)
	a.Equal(p.config.ClientSecret, oauthClientSecret)
	a.Equal(p.organization, organization)
	a.Equal(p.whitelist, map[string]bool{"user": true})
}

func TestWhitelist(t *testing.T) {
	c := &config.Auth{
		OauthClientID:     oauthClientID,
		OauthClientSecret: oauthClientSecret,
		ProviderOpts:      map[string]string{"organization": ""},
		UsersWhitelist:    []string{},
	}
	if _, err := New(c); err == nil {
		t.Error("creating a provider without an organization set should return an error")
	}
	// Set a user whitelist but no domain
	c.UsersWhitelist = users
	if _, err := New(c); err != nil {
		t.Error("creating a provider with users but no organization should not return an error")
	}
	// Unset the user whitelist and set a domain
	c.UsersWhitelist = []string{}
	c.ProviderOpts = map[string]string{"organization": organization}
	if _, err := New(c); err != nil {
		t.Error("creating a provider with an organization set but without a user whitelist should not return an error")
	}
}

func TestStartSession(t *testing.T) {
	a := assert.New(t)

	p, _ := newGithub()
	r := &http.Request{
		Host: oauthCallbackURL,
	}
	s := p.StartSession("test_state", r)
	a.Contains(s.AuthURL, "github.com/login/oauth/authorize")
	a.Contains(s.AuthURL, "state=test_state")
	a.Contains(s.AuthURL, fmt.Sprintf("client_id=%s", oauthClientID))
}

func newGithub() (*Config, error) {
	c := &config.Auth{
		OauthClientID:     oauthClientID,
		OauthClientSecret: oauthClientSecret,
		ProviderOpts:      map[string]string{"organization": organization},
	}
	return New(c)
}
