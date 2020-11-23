FROM golang:1.15 as builder

WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 make build

FROM alpine:3.12

WORKDIR /opt

COPY --from=builder /go/src/app/aws-ecr-login-secret-updater ./job

ENTRYPOINT ["/opt/job"]
