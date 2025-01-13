package main

import (
	"context"
	"fmt"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/aws"
)

// Return a Vault client with AWS IAM auth enabled.
// one of two auth methods used to authenticate with AWS (the other is EC2 auth).
func getVaultClientWithAWSAuthIAM(authPath string, authRole string) (*vault.Client, error) {
	config := vault.DefaultConfig() // modify for more granular configuration

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	awsAuth, err := auth.NewAWSAuth(
		auth.WithMountPath(authPath),
		auth.WithRole(authRole),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize AWS auth method: %w", err)
	}

	authInfo, err := client.Auth().Login(context.Background(), awsAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to login to AWS auth method: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no auth info was returned after login")
	}

	return client, nil
}
