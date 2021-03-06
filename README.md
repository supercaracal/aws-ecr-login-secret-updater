![](https://github.com/supercaracal/aws-ecr-login-secret-updater/workflows/Test/badge.svg)
![](https://github.com/supercaracal/aws-ecr-login-secret-updater/workflows/Docker/badge.svg)

AWS ECR login secret updater
============================

[Kubernetes - How to access AWS ECR](https://dpjanes.medium.com/kubernetes-how-to-accessaws-ecr-bd1e6e6c061)

> Note that the login is only good for 12 hours.

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
          serviceAccountName: ecr-login-updater
          containers:
          - name: aws-ecr-login-secret-updater
            image: ghcr.io/supercaracal/aws-ecr-login-secret-updater:latest
            envFrom:
              - secretRef:
                  name: sample-iam-secret
            env:
              - name: TZ
                value: "Asia/Tokyo"
              - name: AWS_REGION
                value: "ap-northeast-1"
              - name: AWS_ACCOUNT_ID
                value: "000000000000"
              - name: EMAIL
                value: "foo@example.com"
              - name: SECRET
                value: "sample-ecr-login-secret"
              - name: NAMESPACE
                value: "default"
```

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: ecr-login-updater
  namespace: default
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: secret-updater
  namespace: default
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: secrets-manager
subjects:
- kind: ServiceAccount
  name: ecr-login-updater
  namespace: default
```

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: sample-iam-secret
  namespace: default
type: Opaque
data:
  AWS_ACCESS_KEY_ID: **********base64 encoded text**********
  AWS_SECRET_ACCESS_KEY: **********base64 encoded text**********
```

```
$ kubectl get secrets sample-ecr-login-secret -o json | jq -r .data.'".dockerconfigjson"' | base64 -d | jq .
{
  "auths": {
    "https://000000000000.dkr.ecr.ap-northeast-1.amazonaws.com": {
      "username": "AWS",
      "password": "*****************************************",
      "email": "foo@example.com",
      "auth": "*****************************************"
    }
  }
}
```
