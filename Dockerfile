FROM golang:1.21.8-alpine

WORKDIR /app

RUN go mod init vk_app
RUN go get -u github.com/lib/pq github.com/golang-jwt/jwt/v5 golang.org/x/crypto/bcrypt github.com/swaggo/http-swagger/v2

COPY ./ ./

EXPOSE 5000

ENTRYPOINT [ "go", "run", "./cmd/app" ]