package oauth2

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
)

// State -.
type State struct {
	Lang     string `json:"lang"`
	Redirect string `json:"redirect"`
	Nonce    string `json:"nonce"`
}

// MakeState -.
func MakeState(redirect string, lang string) (string, error) {
	nonce, _ := genNonce(32)
	state := State{Lang: lang, Redirect: redirect, Nonce: nonce}
	b, err := json.Marshal(state)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// ParseState -.
func ParseState(state string) (*State, error) {
	var st State

	raw, err := base64.StdEncoding.DecodeString(state)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &st); err != nil {
		return nil, err
	}

	return &st, nil
}

func genNonce(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
