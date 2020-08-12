FROM golang:latest
WORKDIR /go/src/github.com/ScoreTrak/Web
RUN apt-get update
RUN curl -sL https://deb.nodesource.com/setup_14.x | bash -
RUN apt-get update && apt-get install -y nodejs
RUN npm install yarn -g
COPY pkg/ pkg/
COPY views/ views/
COPY cmd/ cmd/
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod tidy
RUN go build -o web cmd/web/main.go
RUN go build -o jobs cmd/jobs/main.go
RUN chmod +x web jobs
RUN cd views && yarn install && yarn build