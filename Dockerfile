FROM golang:latest as builder
WORKDIR /go/src/go-url-shortener
RUN go get -d -v github.com/go-sql-driver/mysql
RUN go get -d -v github.com/gorilla/mux
RUN go get -d -v github.com/spf13/viper
COPY app.go .
COPY config.json .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o short

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/go-url-shortener/config.json .
COPY --from=builder /go/src/go-url-shortener/short .

CMD ["./short"]

EXPOSE 8080
