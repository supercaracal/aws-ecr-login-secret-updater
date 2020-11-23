![](https://github.com/supercaracal/aws-ecr-login-secret-updater/workflows/Test/badge.svg)
![](https://github.com/supercaracal/aws-ecr-login-secret-updater/workflows/Docker/badge.svg)

AWS ECR login secret updater
============================

[Kubernetes - How to access AWS ECR](https://dpjanes.medium.com/kubernetes-how-to-accessaws-ecr-bd1e6e6c061)

> Note that the login is only good for 12 hours.

```
$ kind create cluster
$ kubectl cluster-info --context kind-kind
```

```yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: sample-ecr-login-secret-updater
  namespace: default
spec:
  schedule: "0 */8 * * *"
  successfulJobsHistoryLimit: 2
  failedJobsHistoryLimit: 2
  jobTemplate:
    spec:
      backoffLimit: 0
      template:
        spec:
          terminationGracePeriodSeconds: 0
          restartPolicy: Never
          containers:
          - name: aws-ecr-login-secret-updater
            image: ghcr.io/supercaracal/aws-ecr-login-secret-updater:latest
            envFrom:
              - secretRef:
                  name: sample-ecr-login-secret-updater-secret
            env:
              - name: TZ
                value: "Asia/Tokyo"
              - name: SECRET
                value: "sample-ecr-login-secret"
              - name: NAMESPACE
                value: "default"
```

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: sample-ecr-login-secret-updater-secret
  namespace: default
type: Opaque
data:
  AWS_REGION: **********base64 encoded text**********
  AWS_ACCOUNT_ID: **********base64 encoded text**********
  AWS_ACCESS_KEY_ID: **********base64 encoded text**********
  AWS_SECRET_ACCESS_KEY: **********base64 encoded text**********
  EMAIL: **********base64 encoded text**********
```
