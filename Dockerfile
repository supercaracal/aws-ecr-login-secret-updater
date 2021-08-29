FROM golang:1.17 as builder
WORKDIR /go/src/app
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 make build

# @see https://console.cloud.google.com/gcr/images/distroless/GLOBAL
FROM gcr.io/distroless/static-debian11:latest-amd64
WORKDIR /opt
COPY --from=builder /go/src/app/aws-ecr-login-secret-updater ./job
ENTRYPOINT ["/opt/job"]
