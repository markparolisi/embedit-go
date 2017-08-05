# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.8.3

WORKDIR /go/src/embedit


ADD . .

RUN go get
RUN go build

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/src/embedit/embedit

EXPOSE 8080