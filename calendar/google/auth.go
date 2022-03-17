package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(ctx context.Context, account *Account) (*http.Client, error) {
	token, err := getToken(account.TokenPath)
	if err != nil {
		token, err = requestToken(ctx, account.oauth)
		if err != nil {
			return nil, err
		}

		if err := saveToken(account.TokenPath, token); err != nil {
			return nil, err
		}
	}

	return account.oauth.Client(ctx, token), nil
}

// getToken retrieves the token from a local file.
func getToken(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	token := &oauth2.Token{}
	if err := json.NewDecoder(f).Decode(token); err != nil {
		return nil, errors.Wrap(err, "decoding oauth token")
	}

	return token, nil
}

// requestToken prompts the user with a URL to create an authorization code.
func requestToken(ctx context.Context, oauthConfig *oauth2.Config) (*oauth2.Token, error) {
	authURL := oauthConfig.AuthCodeURL("comeet", oauth2.AccessTypeOffline)
	fmt.Printf("Visit this URL to get an authorization code:\n%s\n", authURL)

	fmt.Printf("Paste the authorization code for client %q here: ", oauthConfig.ClientID)
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, errors.Wrap(err, "reading the authorization code")
	}

	token, err := oauthConfig.Exchange(ctx, authCode)
	if err != nil {
		return nil, errors.Wrap(err, "converting authorization code into token")
	}

	return token, nil
}

// saveToken saves a token to a local file.
func saveToken(path string, token *oauth2.Token) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Wrap(err, "saving OAuth token")
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(token); err != nil {
		return errors.Wrap(err, "enconding token")
	}

	return nil
}
