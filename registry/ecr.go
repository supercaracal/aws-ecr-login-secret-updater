package registry

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// Credential is a parameter for Kubernetes secret of docker-registry type
type Credential struct {
	Server   string
	UserName string
	Password string
	Email    string
}

// ECRClient is a client for AWS ECR service
type ECRClient struct {
	svc    *ecr.ECR
	region string
}

const awsUserNameForRegistry = "AWS"

// NewECRClient is a constructor
func NewECRClient(region, endpointURL string) (*ECRClient, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS client's session: %w", err)
	}

	config := aws.NewConfig().WithRegion(region)
	if endpointURL != "" {
		config = config.WithEndpoint(endpointURL)
	}

	return &ECRClient{svc: ecr.New(sess, config), region: region}, nil
}

// Login is authorization for AWS ECR
func (c *ECRClient) Login(accountID, email string) (*Credential, error) {
	input := &ecr.GetAuthorizationTokenInput{RegistryIds: []*string{aws.String(accountID)}}

	result, err := c.svc.GetAuthorizationToken(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecr.ErrCodeServerException:
				return nil, fmt.Errorf("%s %w", ecr.ErrCodeServerException, aerr)
			case ecr.ErrCodeInvalidParameterException:
				return nil, fmt.Errorf("%s %w", ecr.ErrCodeInvalidParameterException, aerr)
			}
		}
		return nil, err
	}
	if len(result.AuthorizationData) == 0 {
		return nil, fmt.Errorf("failed to get auth token from AWS ECR")
	}

	token, err := base64.StdEncoding.DecodeString(*result.AuthorizationData[0].AuthorizationToken)
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
