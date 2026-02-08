package github

import (
	"context"

	go_github "github.com/google/go-github/v81/github"
	"golang.org/x/oauth2"
	oauth2gh "golang.org/x/oauth2/github"
)

// Client -.
type Client interface {
	GetAuthURL(state string) string
	GetToken(ctx context.Context, code string) (*oauth2.Token, error)
	GetUser(ctx context.Context, token *oauth2.Token) (*go_github.User, error)
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
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     oauth2gh.Endpoint,
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
func (o *clientImpl) GetUser(ctx context.Context, token *oauth2.Token) (*go_github.User, error) {
	client := go_github.NewClient(o.config.Client(ctx, token))

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}

	return user, nil
}
