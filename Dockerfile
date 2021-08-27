FROM golang:1.16.4-alpine as builder
WORKDIR /go/src
COPY . .
RUN go build -o user_api .

FROM alpine:3.14.0
WORKDIR /bin
COPY --from=builder /go/src .
USER nobody
EXPOSE 9090
CMD ["user_api"]
