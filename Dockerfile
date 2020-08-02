FROM golang:latest
WORKDIR /go/src/github.com/L1ghtman2k/ScoreTrakWeb
COPY pkg/ pkg/
COPY cmd/ cmd/
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod tidy
RUN go build -o web cmd/web/main.go
RUN chmod +x web