FROM golang:latest
WORKDIR /go/src/github.com/L1ghtman2k/ScoreTrakWeb
COPY deployments .
RUN go mod tidy
RUN go build -o web cmd/web/main.go
RUN chmod +x web
RUN mv web /tmp/web
RUN rm -rf *
RUN mv /tmp/web web

#Set Context Path as ScoreTrakWeb directory