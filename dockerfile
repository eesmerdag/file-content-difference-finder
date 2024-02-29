FROM golang:1.21-alpine

WORKDIR /file-diff-finder

COPY . /file-diff-finder

RUN GOOS=linux GOARCH=amd64 go build -o bin/service /file-diff-finder/cmd/main.go

FROM alpine:3.18.3

COPY --from=0 /file-diff-finder/bin/service /go/bin/service

ENTRYPOINT ["/go/bin/service", "1903", "13", "This is a simple file content."]