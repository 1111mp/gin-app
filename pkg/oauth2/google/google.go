package google

import (
	"context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleOauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

// Client -.
type Client interface {
	GetAuthURL(state string) string
	GetToken(ctx context.Context, code string) (*oauth2.Token, error)
	GetUser(ctx context.Context, token *oauth2.Token) (*googleOauth2.Userinfo, error)
}

// clientImpl -.
type clientImpl struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string

	config *oauth2.Config
}

// Setup -.
func Setup(opts ...Option) Client {
	client := &clientImpl{}

	// Custom options
	for _, opt := range opts {
		opt(client)
	}

	client.config = &oauth2.Config{
		ClientID:     client.ClientID,
		ClientSecret: client.ClientSecret,
		RedirectURL:  client.RedirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}

	return client
}

// GetAuthURL -.
func (o *clientImpl) GetAuthURL(state string) string {
	return o.config.AuthCodeURL(state)
}

// GetToken -.
func (o *clientImpl) GetToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return o.config.Exchange(ctx, code)
}

// GetUser -.
func (o *clientImpl) GetUser(ctx context.Context, token *oauth2.Token) (*googleOauth2.Userinfo, error) {
	client := o.config.Client(ctx, token)
	svc, err := googleOauth2.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	user, err := svc.Userinfo.Get().Do()
	if err != nil {
		return nil, err
	}

	return user, nil
}
