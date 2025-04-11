ARG VERSION="dev"

FROM golang:1.23 AS build
# allow this step access to build arg
ARG VERSION
# Set the working directory
WORKDIR /build

RUN go env -w GOMODCACHE=/root/.cache/go-build

# Install dependencies
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build go mod download

COPY . ./
# Build the server
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 go build \
    -o oosa-mcp-server cmd/oosa-mcp-server/main.go

# Make a stage to run the app
FROM gcr.io/distroless/base-debian12
# Set the working directory
WORKDIR /server
# Copy the binary and config from the build stage
COPY --from=build /build/oosa-mcp-server .
COPY .oosa-mcp-server.yaml /server/.oosa-mcp-server.yaml
# Command to run the server
CMD ["./oosa-mcp-server", "serve", "--config", "/server/.oosa-mcp-server.yaml"]
