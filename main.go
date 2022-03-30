package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codeartifact"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
)

func setupViper() {
	var (
		organization            = flag.String("DEPENDABOT_ORG", os.Getenv("DEPENDABOT_ORG"), "the GitHub organization for which the secret should be created")
		githubSecret            = flag.String("GITHUB_PRIVATE_KEY", os.Getenv("GITHUB_PRIVATE_KEY"), "GitHub secret for GitHub App authentication")
		githubAppID             = flag.String("GITHUB_APP_ID", os.Getenv("GITHUB_APP_ID"), "the ID of the GitHub App used for authentication")
		organizationOwner       = flag.String("DEPENDABOT_OWNER", os.Getenv("DEPENDABOT_OWNER"), " owner of the GitHub organization")
		tokenDuration           = flag.String("CODEARTIFACT_DURATION", os.Getenv("CODEARTIFACT_DURATION"), "duration of the AWS CodeArtifact authToken")
		codeartifactDomain      = flag.String("CODEARTIFACT_DOMAIN", os.Getenv("CODEARTIFACT_DOMAIN"), "AWS CodeArtifact Domain for which access is required")
		codeartifactDomainOwner = flag.String("CODEARTIFACT_DOMAIN_OWNER", os.Getenv("CODEARTIFACT_DOMAIN_OWNER"), "owner (AWS acc) for the AWS CodeArtifact domain")
	)

	flag.Parse()

	viper.Set("GITHUB_APP_ID", githubAppID)
	viper.Set("GITHUB_PRIVATE_KEY", githubSecret)
	viper.Set("DEPENDABOT_ORG", organization)
	viper.Set("DEPENDABOT_OWNER", organizationOwner)
	viper.Set("CODEARTIFACT_DURATION", tokenDuration)
	viper.Set("CODEARTIFACT_DOMAIN", codeartifactDomain)
	viper.Set("CODEARTIFACT_DOMAIN_OWNER", codeartifactDomainOwner)
}

func main() {
	setupViper()

	ctx, cancel := context.WithCancel(context.Background())

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		select {
		case <-sigint:
			cancel()
		case <-ctx.Done():
			if ctx.Err() != nil {
				logrus.Errorf("received error in context: %v", ctx.Err())
			}
			logrus.Info("context closed. Shutting down ... ")
			return
		}
	}()

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())

		http.ListenAndServe("0.0.0.0:8701", mux)
	}()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logrus.Fatalf("AWS config could not be created: %v", err)
	}

	run(ctx, cfg)

	for range time.NewTicker(time.Hour * 10).C {
		run(ctx, cfg)
	}
}

func run(ctx context.Context, cfg aws.Config) {
	// Get AWS CodeArtifact secret
	secret, err := getCodeArtifactSecret(ctx, cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	ghClient, err := setupGitHubAppClient(ctx)
	if err != nil {
		logrus.Errorf("run failed with: %v", err)
	}

	if err := createOrUpdateDependabotSecret(ctx, ghClient, *secret); err != nil {
		logrus.Error(err)
	}
}

func getCodeArtifactSecret(ctx context.Context, cfg aws.Config) (*string, error) {
	var (
		domain      string = viper.GetString("CODEARTIFACT_DOMAIN")
		domainOwner string = viper.GetString("CODEARTIFACT_DOMAIN_OWNER")
		duration    int64  = viper.GetInt64("CODEARTIFACT_DURATION")
	)

	client := codeartifact.NewFromConfig(cfg)
	out, err := client.GetAuthorizationToken(ctx, &codeartifact.GetAuthorizationTokenInput{
		DurationSeconds: &duration,
		Domain:          &domain,
		DomainOwner:     &domainOwner,
	})
	if err != nil {
		return nil, fmt.Errorf("retrieving code artifact secret: %w", err)
	}

	return out.AuthorizationToken, nil
}
