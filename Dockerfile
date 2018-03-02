FROM golang:alpine  as builder

RUN apk add --no-cache make gcc musl-dev linux-headers

WORKDIR /go/src/seth/


ADD . /go/src/seth/

RUN CGO_ENABLED=1 GOOS=linux go build -a  -o seth .

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /go/src/seth/seth .
RUN chmod +x /root/seth

ENTRYPOINT ["/root/seth"]
#CMD ["start"]