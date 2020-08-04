FROM golang:latest
WORKDIR /go/src/github.com/ScoreTrak/Web
COPY pkg/ pkg/
COPY views/ views/
COPY cmd/ cmd/
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod tidy
RUN go build -o web cmd/web/main.go
RUN go build -o jobs cmd/jobs/main.go
RUN chmod +x web jobs