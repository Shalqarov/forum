package web

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type UserInfo struct {
	Email    string `json:"email"`
	Username string `json:"name"`
}
