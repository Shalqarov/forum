package web

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type Endpoint struct {
	AuthURL  string
	TokenURL string
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type UserInfo struct {
	Email    string `json:"email"`
	Username string `json:"name"`
}

type Config struct {
	ClientID     string
	ClientSecret string
	Endpoint     Endpoint
	RedirectURL  string
	Scopes       []string
}

var state = uuid.NewV4().String()

func (c *Config) AuthCodeURL(state string) string {
	var buf bytes.Buffer
	buf.WriteString(c.Endpoint.AuthURL)
	v := url.Values{
		"response_type": {"code"},
		"client_id":     {c.ClientID},
	}
	if c.RedirectURL != "" {
		v.Set("redirect_uri", c.RedirectURL)
	}
	if len(c.Scopes) > 0 {
		v.Set("scope", strings.Join(c.Scopes, " "))
	}
	if state != "" {
		v.Set("state", state)
	}
	if strings.Contains(c.Endpoint.AuthURL, "?") {
		buf.WriteByte('&')
	} else {
		buf.WriteByte('?')
	}
	buf.WriteString(v.Encode())
	return buf.String()
}

func (c *Config) GetTokenByCode(code string) (*Token, error) {
	v := url.Values{
		"grant_type": {"authorization_code"},
		"code":       {code},
	}
	if c.RedirectURL != "" {
		v.Set("redirect_uri", c.RedirectURL)
	}
	return retrieveToken(c, v)
}

func retrieveToken(c *Config, v url.Values) (*Token, error) {
	req, err := newTokenRequest(c.Endpoint.TokenURL, c.ClientID, c.ClientSecret, v)
	if err != nil {
		return nil, err
	}
	token, err := requestDo(req)
	if err != nil {
		req, _ := newTokenRequest(c.Endpoint.TokenURL, c.ClientID, c.ClientSecret, v)
		token, err = requestDo(req)
	}
	return token, err
}

func newTokenRequest(tokenURL, clientID, clientSecret string, v url.Values) (*http.Request, error) {
	v = cloneURLValues(v)
	if clientID != "" {
		v.Set("client_id", clientID)
	}
	if clientSecret != "" {
		v.Set("client_secret", clientSecret)
	}
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func cloneURLValues(v url.Values) url.Values {
	v2 := make(url.Values, len(v))
	for k, vv := range v {
		v2[k] = append([]string(nil), vv...)
	}
	return v2
}

func requestDo(req *http.Request) (*Token, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	token := &Token{}
	json.Unmarshal(respBody, &token)
	return token, nil
}
