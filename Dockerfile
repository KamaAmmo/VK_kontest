FROM golang:1.21.8-alpine

WORKDIR /app

RUN go mod init vk_app
RUN go get -u github.com/lib/pq

COPY ./ ./

EXPOSE 5000

ENTRYPOINT [ "go", "run", "./cmd/app" ]