package github

// Option -.
type Option func(*clientImpl)

// ClientID -.
func ClientID(clientID string) Option {
	return func(client *clientImpl) {
		client.ClientID = clientID
	}
}

// ClientSecret -.
func ClientSecret(clientSecret string) Option {
	return func(client *clientImpl) {
		client.ClientSecret = clientSecret
	}
}

// RedirectURL -.
func RedirectURL(redirectURL string) Option {
	return func(client *clientImpl) {
		client.RedirectURL = redirectURL
	}
}
