FROM golang:latest
WORKDIR /go/src/github.com/L1ghtman2k/ScoreTrakWeb
COPY deployments .
RUN go mod tidy
RUN go build -o web cmd/web/main.go
RUN chmod +x web

#Set Context Path as ScoreTrakWeb directory