FROM golang:1.22-alpine as builder

COPY . /app
WORKDIR /app
RUN GOOS=linux go build cmd/cloudflare-ddns-update.go

FROM golang:1.22-alpine

COPY --from=builder /app/cloudflare-ddns-update /

ENTRYPOINT ["/cloudflare-ddns-update"]

LABEL io.artifacthub.package.readme-url=https://github.com/dploeger/cloudflare-ddns-update
LABEL org.opencontainers.image.created=2024-06-06
LABEL org.opencontainers.image.description="Dyn API to update Cloudflare DNS records"
LABEL org.opencontainers.image.documentation=https://github.com/dploeger/cloudflare-ddns-update
LABEL org.opencontainers.image.source=https://github.com/dploeger/cloudflare-ddns-update
LABEL org.opencontainers.image.title="cloudflare-ddns-update"
LABEL org.opencontainers.image.url=https://github.com/dploeger/cloudflare-ddns-update
LABEL org.opencontainers.image.vendor="Dennis Plöger"
