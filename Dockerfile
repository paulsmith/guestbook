FROM golang:1.8 as builder
MAINTAINER ops@adhocteam.us
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 go build -o app .

FROM alpine:latest
MAINTAINER ops@adhocteam.us
WORKDIR /root/
COPY --from=0 /go/src/app/app .
CMD ["./app"]
