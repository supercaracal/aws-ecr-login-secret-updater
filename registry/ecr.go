package registry

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

// @see https://pkg.go.dev/github.com/aws/aws-sdk-go-v2

// Credential is a parameter for Kubernetes secret of docker-registry type
type Credential struct {
	Server   string
	UserName string
	Password string
	Email    string
}

// ECRClient is a client for AWS ECR service
type ECRClient struct {
	svc    *ecr.Client
	region string
}

const (
	awsUserNameForRegistry = "AWS"
	timeout                = 10 * time.Second
)

// NewECRClient is a constructor
func NewECRClient(region, endpointURL string) (*ECRClient, error) {
	cfg, err := loadAWSConfig(region, endpointURL)
	if err != nil {
		return nil, err
	}

	return &ECRClient{svc: ecr.NewFromConfig(cfg), region: region}, nil
}

func loadAWSConfig(region, endpointURL string) (aws.Config, error) {
	if endpointURL == "" {
		return config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	}

	return config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithEndpointResolver(
			aws.EndpointResolverFunc(
				func(service, region string) (aws.Endpoint, error) {
					return aws.Endpoint{
						PartitionID:   "aws",
						URL:           endpointURL,
						SigningRegion: region,
					}, nil
				},
			),
		),
	)
}

// Login is authorization for AWS ECR
func (c *ECRClient) Login(accountID, email string) (*Credential, error) {
	input := &ecr.GetAuthorizationTokenInput{RegistryIds: []string{accountID}}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	output, err := c.svc.GetAuthorizationToken(ctx, input)
	if err != nil {
		return nil, err
	}
	if len(output.AuthorizationData) == 0 {
		return nil, fmt.Errorf("failed to get auth token from AWS ECR")
	}

	token, err := base64.StdEncoding.DecodeString(*output.AuthorizationData[0].AuthorizationToken)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(token), ":")
	if len(parts) < 2 {
		return nil, fmt.Errorf("failed to parse auth token of AWS ECR")
	}

	return &Credential{
		Server:   fmt.Sprintf("https://%s.dkr.ecr.%s.amazonaws.com", accountID, c.region),
		UserName: awsUserNameForRegistry,
		Password: parts[1],
		Email:    email,
	}, nil
}
