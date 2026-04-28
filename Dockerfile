# Build all service binaries in one image; compose runs each with a different command.
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN mkdir -p /out && \
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/mcpserver ./cmd/mcpserver && \
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/onboarding-agent .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
COPY --from=build /out/mcpserver /out/onboarding-agent /usr/local/bin/
