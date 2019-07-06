FROM golang:alpine

RUN mkdir -p /go/src/github.com/jongschneider/youtube-project/api
WORKDIR /go/src/github.com/jongschneider/youtube-project/api

COPY ./ ./
RUN go install ./...
CMD ["cmd"]
