# not pinning the exact image digest to support multiple architecture builds without multiple Dockerfiles

FROM golang:1.26@sha256:3aff6657219a4d9c14e27fb1d8976c49c29fddb70ba835014f477e1c70636647 AS builder
COPY . /var/app
WORKDIR /var/app
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o app .

FROM alpine:3.24@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b
RUN apk update && apk add ca-certificates
COPY --from=builder /var/app/app /var/app/app
CMD ["/var/app/app"]
