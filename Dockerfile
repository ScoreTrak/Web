FROM golang:latest
WORKDIR /go/src/github.com/L1ghtman2k/ScoreTrakWeb
COPY . .
RUN go mod tidy
RUN go build -o web cmd/web/main.go
RUN chmod +x web