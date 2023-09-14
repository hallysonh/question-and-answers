############################
# STEP 1 Create build base to cache build tools
############################
FROM golang:alpine as build_base

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
# GCC and Muls-dev is required to build some dependencies
# Tzdata is required to provide timezone data
# Create appuser
RUN set -eux; apk add git bash ca-certificates tzdata gcc g++ libc-dev; \
    update-ca-certificates 2>/dev/null || true; \
    adduser -D -g '' appuser

WORKDIR /app

# Force the go compiler to use modules
ENV GO111MODULE=on

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

# Fetch dependencies. Using go mod with go 1.11
RUN go mod download

############################
# STEP 2 build executable binary
############################
FROM build_base AS builder

# Here we copy the rest of the source code
COPY . .

# Build the binary
RUN cd cmd/api &&\
    CGO_ENABLED=1 GOOS=linux \
    go build -installsuffix 'static' -ldflags="-w -s" -o app

############################
# STEP 3 build a small image
############################
FROM alpine as final

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/cmd/api/app /api
COPY --from=builder /app /app

# Config timezone
RUN cd /tmp &&\
    apk update &&\
    apk add --no-cache libaio libnsl &&\
    rm -rf /var/cache/apk/* &&\
    rm -rf /etc/localtime &&\
    ln -s /usr/share/zoneinfo/Brazil/West /etc/localtime

ENV APP_ENV prod

# Use an unprivileged user.
USER appuser

# Run the main binary.
ENTRYPOINT /api -log-level=INFO
