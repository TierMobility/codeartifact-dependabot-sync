package main

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/go-github/v42/github"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/oauth2"
)

func newGitHubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client
}

type GithubToken struct {
	Token               string      `json:"token"`
	ExpiresAt           time.Time   `json:"expires_at"`
	Permissions         Permissions `json:"permissions"`
	RepositorySelection string      `json:"repository_selection"`
}

type Permissions struct {
	OrganizationDependabotSecrets string `json:"organization_dependabot_secrets"`
	DependabotSecrets             string `json:"dependabot_secrets"`
	Metadata                      string `json:"metadata"`
}

func getJWT() (*string, error) {
	pemBytes := []byte(viper.GetString("GITHUB_PRIVATE_KEY"))

	block, _ := pem.Decode(pemBytes)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("creating JWT parsing private key: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": time.Now().Add(-1 * time.Minute).Unix(),
		"exp": time.Now().Add(time.Minute * 10).Unix(),
		"iss": viper.GetString("GITHUB_APP_ID"),
	})

	signedToken, err := token.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("creating JWT signing token: %w", err)
	}

	return &signedToken, nil
}

func setupGitHubAppClient(ctx context.Context) (*github.Client, error) {
	signedToken, err := getJWT()
	if err != nil {
		return nil, fmt.Errorf("setting up GitHub App client: %w", err)
	}

	tempClient := newGitHubClient(ctx, *signedToken)

	inst, _, err := tempClient.Apps.FindOrganizationInstallation(ctx, "TierMobility")
	if err != nil {
		return nil, fmt.Errorf("setting up GitHub App client: %w", err)
	}

	ghAppToken, err := retrieveGHAppToken(*signedToken, *inst.AccessTokensURL)
	if err != nil {
		return nil, fmt.Errorf("setting up GitHub App client: %w", err)
	}

	finalClient := newGitHubClient(ctx, ghAppToken.Token)

	return finalClient, nil
}

func retrieveGHAppToken(signedToken, accessTokensURL string) (*GithubToken, error) {
	req, err := http.NewRequest(http.MethodPost, accessTokensURL, nil)
	if err != nil {
		return nil, fmt.Errorf("retrieving GH app token: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+signedToken)

	nRes, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("retrieving GH app token: %w", err)
	}

	b, err := io.ReadAll(nRes.Body)
	if err != nil {
		return nil, fmt.Errorf("retrieving GH app token: %w", err)
	}

	var gt GithubToken
	err = json.Unmarshal(b, &gt)
	if err != nil {
		return nil, fmt.Errorf("retrieving GH app token: %w", err)
	}

	return &gt, nil
}

// encryptSecret is required by the github API for the dependabot/secrets endpoint
// as described in the github documentation:
// https://docs.github.com/en/rest/reference/dependabot#create-or-update-a-repository-secret
//
// The documentation requires the libsodium encryption of the key
// which is implemented in golang via the NACL implementation in the crypto
// package. We implement the interface as slim as possible by creating
// a useless nonce and copy the public und private key (github token) into
// a fixed size byte array, which is required from the implementation
func encryptSecret(plainSecret, key, tok string) (*string, error) {
	// copy publicKey byteSlice into fixed size byte array
	var publicKey [32]byte
	_, err := base64.StdEncoding.Decode(publicKey[:], []byte(key))
	if err != nil {
		return nil, fmt.Errorf("encrypt secret decode key")
	}
	// copy(publicKey[:], []byte(key))

	enc, err := box.SealAnonymous(nil, []byte(plainSecret), &publicKey, nil)
	if err != nil {
		return nil, err
	}

	strEnc := base64.StdEncoding.EncodeToString(enc)

	return &strEnc, nil
}

func createOrUpdateDependabotSecret(ctx context.Context, ghClient *github.Client, secret string) error {
	var (
		org   = viper.GetString("DEPENDABOT_ORG")
		token = viper.GetString("GITHUB_APP_TOKEN")
	)

	pk, _, err := ghClient.Dependabot.GetOrgPublicKey(ctx, org)
	if err != nil {
		return fmt.Errorf("creating dependabot secret: %w", err)
	}

	strEnc, err := encryptSecret(secret, pk.GetKey(), token)
	if err != nil {
		return fmt.Errorf("creating dependabot secret in encryption: %w", err)
	}

	// list all repositories for the authenticated user
	res, err := ghClient.Dependabot.CreateOrUpdateOrgSecret(ctx, org, &github.EncryptedSecret{
		Name:           "CodeArtifactSecret",
		KeyID:          *pk.KeyID,
		EncryptedValue: *strEnc,
		Visibility:     "all",
	})
	if err != nil {
		return fmt.Errorf("creating dependabot secret: %w", err)
	}

	parseGHResponse(res)

	up.Add(1)

	return nil
}

func parseGHResponse(res *github.Response) {
	if res.StatusCode == http.StatusCreated {
		logrus.Info("created new dependabot secret")
	} else if res.StatusCode == http.StatusNoContent {
		logrus.Info("succesfully update existing secret")
	} else {
		logrus.Error("creating/Updating the dependabot secret failed")
	}
}
