package google

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
	domain            = "example.com"
	users             = []string{"user"}
)

func TestNew(t *testing.T) {
	a := assert.New(t)
	p, err := newGoogle()
	a.NoError(err)
	a.Equal(p.config.ClientID, oauthClientID)
	a.Equal(p.config.ClientSecret, oauthClientSecret)
	a.Equal(p.domain, domain)
	a.Equal(p.whitelist, map[string]bool{"user": true})
}

func TestWhitelist(t *testing.T) {
	c := &config.Auth{
		OauthClientID:     oauthClientID,
		OauthClientSecret: oauthClientSecret,
		ProviderOpts:      map[string]string{"domain": ""},
		UsersWhitelist:    []string{},
	}
	if _, err := New(c); err == nil {
		t.Error("creating a provider without a domain set should return an error")
	}
	// Set a user whitelist but no domain
	c.UsersWhitelist = users
	if _, err := New(c); err != nil {
		t.Error("creating a provider with users but no domain should not return an error")
	}
	// Unset the user whitelist and set a domain
	c.UsersWhitelist = []string{}
	c.ProviderOpts = map[string]string{"domain": domain}
	if _, err := New(c); err != nil {
		t.Error("creating a provider with a domain set but without a user whitelist should not return an error")
	}
}

func TestStartSession(t *testing.T) {
	a := assert.New(t)

	p, err := newGoogle()
	a.NoError(err)
	r := &http.Request{
		Host: oauthCallbackURL,
	}
	s := p.StartSession("test_state", r)
	a.Contains(s.AuthURL, "accounts.google.com/o/oauth2/auth")
	a.Contains(s.AuthURL, "state=test_state")
	a.Contains(s.AuthURL, fmt.Sprintf("hd=%s", domain))
	a.Contains(s.AuthURL, fmt.Sprintf("client_id=%s", oauthClientID))
}

func newGoogle() (*Config, error) {
	c := &config.Auth{
		OauthClientID:     oauthClientID,
		OauthClientSecret: oauthClientSecret,
		ProviderOpts:      map[string]string{"domain": domain},
		UsersWhitelist:    users,
	}
	return New(c)
}
