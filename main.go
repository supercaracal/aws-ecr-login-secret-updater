package main

import (
	"flag"
	"log"
	"os"

	"github.com/supercaracal/aws-ecr-login-secret-updater/kube"
	"github.com/supercaracal/aws-ecr-login-secret-updater/registry"
)

type config struct {
	awsRegion      string
	awsEndpointURL string
	awsAccountID   string
	email          string
	secret         string
	namespace      string
	masterURL      string
	kubeconfig     string
}

const defaultNamespace = "default"

var cfg config

func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	if cfg.awsRegion == "" || cfg.awsAccountID == "" || cfg.email == "" || cfg.secret == "" {
		logger.Fatalln("lack of required env vars")
	}

	ecrCli, err := registry.NewECRClient(cfg.awsRegion, cfg.awsEndpointURL)
	if err != nil {
		logger.Fatalln(err)
	}

	credential, err := ecrCli.Login(cfg.awsAccountID, cfg.email)
	if err != nil {
		logger.Fatalln(err)
	}

	kubeCli, err := kube.NewClient(cfg.masterURL, cfg.kubeconfig, cfg.namespace)
	if err != nil {
		logger.Fatalln(err)
	}

	err = kubeCli.UpdateSecret(
		cfg.secret,
		credential.Server,
		credential.UserName,
		credential.Password,
		credential.Email,
	)
	if err != nil {
		logger.Fatalln(err)
	}

	os.Exit(0)
}

func init() {
	cfg.awsRegion = os.Getenv("AWS_REGION")
	cfg.awsEndpointURL = os.Getenv("AWS_ENDPOINT_URL")
	cfg.awsAccountID = os.Getenv("AWS_ACCOUNT_ID")
	cfg.email = os.Getenv("EMAIL")
	cfg.secret = os.Getenv("SECRET")

	cfg.namespace = os.Getenv("NAMESPACE")
	if cfg.namespace == "" {
		cfg.namespace = defaultNamespace
	}

	flag.StringVar(
		&cfg.masterURL,
		"master",
		"",
		"The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.",
	)

	flag.StringVar(
		&cfg.kubeconfig,
		"kubeconfig",
		"",
		"Path to a kubeconfig. Only required if out-of-cluster.",
	)
}
