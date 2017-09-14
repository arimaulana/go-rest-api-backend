FROM golang:alpine

WORKDIR /go/src/github.com/arimaulana/GoMeds/api

COPY . .

RUN apk add --update git

# RUN go get ./
# RUN go build

CMD go run api.go

EXPOSE 3000
