module github.com/supercaracal/aws-ecr-login-secret-updater

go 1.15

require (
	github.com/aws/aws-sdk-go v1.35.33
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	golang.org/x/oauth2 v0.0.0-20201109201403-9fd604954f58 // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	k8s.io/api v0.19.4
	k8s.io/apimachinery v0.19.4
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.19.4
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.5-rc.0
	k8s.io/apiserver => k8s.io/apiserver v0.19.4
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.19.4
	k8s.io/client-go => k8s.io/client-go v0.19.4
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.19.5-0.20201113181133-070bf588610e
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.19.5-0.20201113181704-acbca43bf834
	k8s.io/code-generator => k8s.io/code-generator v0.19.5-rc.0
	k8s.io/component-base => k8s.io/component-base v0.19.4
)
