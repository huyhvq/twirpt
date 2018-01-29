FROM golang:1.9 AS builder

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && chmod +x /usr/local/bin/dep

RUN mkdir -p /go/src/github.com/huyhvq/twirpt
WORKDIR /go/src/github.com/huyhvq/twirpt

COPY . ./
RUN dep ensure -vendor-only

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /go/src/github.com/huyhvq/twirpt/app .
EXPOSE 8080
CMD ["./app"]