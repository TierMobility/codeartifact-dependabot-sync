FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /codeartifact-dependabot-sync

FROM gcr.io/distroless/base:nonroot

WORKDIR /

COPY --from=build /codeartifact-dependabot-sync /codeartifact-dependabot-sync

ENTRYPOINT ["/codeartifact-dependabot-sync"]
