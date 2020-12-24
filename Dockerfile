FROM golang:1.14

WORKDIR /go/src/app
COPY ./ .
EXPOSE 8080

RUN go get -d -v ./...
RUN go install -v ./...


CMD ["app"]
   