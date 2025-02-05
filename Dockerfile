FROM golang:alpine as builder

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY ./client ./client
COPY --from=builder /app/main ./
CMD ["/app/main"]