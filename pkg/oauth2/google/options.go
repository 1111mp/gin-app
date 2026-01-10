package google

// Option -.
type Option func(*Client)

// ClientID -.
func ClientID(clientID string) Option {
	return func(client *Client) {
		client.ClientID = clientID
	}
}

// ClientSecret -.
func ClientSecret(clientSecret string) Option {
	return func(client *Client) {
		client.ClientSecret = clientSecret
	}
}

// RedirectURL -.
func RedirectURL(redirectURL string) Option {
	return func(client *Client) {
		client.RedirectURL = redirectURL
	}
}
