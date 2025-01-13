package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
)

func handleRequest(ctx context.Context, event json.RawMessage) error {

	// Unmarshal the event
	var input map[string]interface{}
	err := json.Unmarshal(event, &input)
	if err != nil {
		return errors.New("failed to unmarshal event")
	}

	// Auth related env vars
	authPath := os.Getenv("VAULT_AUTH_PROVIDER") // if not provided, Vault will default to "aws"
	authRole := os.Getenv("VAULT_AUTH_ROLE")     // if not provided, Vault will fall back on looking for a role with the IAM role name if you're using the iam auth type, or the EC2 instance's AMI id if using the ec2 auth type

	// Transform related env vars
	transformPath := os.Getenv("VAULT_TRANSFORM_PATH")
	if transformPath == "" {
		return errors.New("VAULT_TRANSFORM_PATH environment var is required")
	}
	transformRole := os.Getenv("VAULT_TRANSFORM_ROLE")
	if transformRole == "" {
		return errors.New("VAULT_TRANSFORM_ROLE environment var is required")
	}

	api, err := getVaultClientWithAWSAuthIAM(authPath, authRole)
	if err != nil {
		return errors.Wrap(err, "failed to get Vault client with AWS IAM auth")
	}

	// for reference: https://docs.aws.amazon.com/redshift/latest/dg/udf-creating-a-lambda-sql-udf.html
	arguments := input["arguments"].([]interface{})
	for _, argument := range arguments {
		// use the Vault client transform decode here
		response, err := api.Logical().Write(transformPath+"/decode/"+transformRole, map[string]interface{}{
			"value":          argument, // TODO: type may need to be converted
			"decode_format":  "TODO",
			"transformation": "TODO",
		})
		if err != nil {
			return errors.Wrap(err, "failed to decode argument")
		}

		decoded := response.Data["decoded_value"] // do something with the decoded value
		log.Printf("decoded value: %v", decoded)
	}

	return nil
}

func main() {
	lambda.Start(handleRequest)
}
