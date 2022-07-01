package web

const (
	// ClientID     = "820533650499-t1lg2j1tl2162t2sldeo9tp3sj4itj3k.apps.githubusercontent.com"
	githubClientID = "e6775f24b7ae5d23a9bf"
	// ClientSecret = "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"
	githubClientSecret = "669c62b1390ac83f7f64ba7e497cb5b4ec0be5f1"
)

var (
	githubConfigSignIn = &Config{
		RedirectURL:  "http://localhost:5000/signin/github/callback",
		ClientID:     githubClientID,
		ClientSecret: githubClientSecret,
		Scopes: []string{
			"https://www.githubapis.com/auth/userinfo.email",
			"https://www.githubapis.com/auth/userinfo.profile",
		},
		Endpoint: githubEndPoint,
	}
	githubConfigSignUp = &Config{
		RedirectURL:  "http://localhost:5000/signup/github/callback",
		ClientID:     githubClientID,
		ClientSecret: githubClientSecret,
		Scopes: []string{
			"https://www.githubapis.com/auth/userinfo.email",
			"https://www.githubapis.com/auth/userinfo.profile",
		},
		Endpoint: githubEndPoint,
	}
	githubEndPoint = Endpoint{
		AuthURL:  "https://github.com/login/oauth/authorize",
		TokenURL: "https://github.com/login/oauth/access_token",
	}
)
