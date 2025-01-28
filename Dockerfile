# CONTAINER FOR BUILDING BINARY
FROM golang:1.21 AS build

# INSTALL DEPENDENCIES
WORKDIR /src
RUN go install github.com/gobuffalo/packr/v2/packr2@v2.8.3
COPY go.mod go.sum ./
RUN go mod download

# BUILD BINARY
COPY . .
RUN packr2 && make build

# CONTAINER FOR RUNNING BINARY
FROM alpine:3.18

# INSTALL REQUIRED PACKAGES
RUN apk add --no-cache postgresql15-client

# COPY BUILT BINARY AND CONFIG FILE
COPY --from=build /src/dist/zkevm-node /app/zkevm-node
COPY --from=build /src/config/environments/testnet/node.config.toml /app/example.config.toml

# SET EXECUTION ENVIRONMENT
WORKDIR /app
EXPOSE 8123
CMD ["./zkevm-node", "run"]
