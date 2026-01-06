package github

// Option -.
type Option func(*OAuth2)

// ClientID -.
func ClientID(clientID string) Option {
	return func(oauth2 *OAuth2) {
		oauth2.ClientID = clientID
	}
}

// ClientSecret -.
func ClientSecret(clientSecret string) Option {
	return func(oauth2 *OAuth2) {
		oauth2.ClientSecret = clientSecret
	}
}

// ClientID -.
func RedirectURL(redirectURL string) Option {
	return func(oauth2 *OAuth2) {
		oauth2.RedirectURL = redirectURL
	}
}
