package github

import (
	"context"

	go_github "github.com/google/go-github/v81/github"
	"golang.org/x/oauth2"
	oauth2gh "golang.org/x/oauth2/github"
)

type OAuth2Inter interface {
	GetAuthURL(state string) string
	GetToken(ctx context.Context, code string) (*oauth2.Token, error)
	GetUser(ctx context.Context, token *oauth2.Token) (*go_github.User, error)
}

// GithubOAuth2 -.
type OAuth2 struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string

	config *oauth2.Config
}

// Setup -.
func Setup(opts ...Option) *OAuth2 {
	githubOAuth2 := &OAuth2{}

	// Custom options
	for _, opt := range opts {
		opt(githubOAuth2)
	}

	githubOAuth2.config = &oauth2.Config{
		ClientID:     githubOAuth2.ClientID,
		ClientSecret: githubOAuth2.ClientSecret,
		RedirectURL:  githubOAuth2.RedirectURL,
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     oauth2gh.Endpoint,
	}

	return githubOAuth2
}

// GetAuthURL -.
func (o *OAuth2) GetAuthURL(state string) string {
	return o.config.AuthCodeURL(state)
}

// GetToken -.
func (o *OAuth2) GetToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return o.config.Exchange(ctx, code)
}

// GetUser -.
func (o *OAuth2) GetUser(ctx context.Context, token *oauth2.Token) (*go_github.User, error) {
	client := go_github.NewClient(o.config.Client(ctx, token))

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}

	return user, nil
}
