FROM golang:1.13

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go build -o app main.go signup.go signin.go db.go user.go

CMD ["./app"]