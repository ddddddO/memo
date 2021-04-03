FROM golang:1.13-alpine as builder

WORKDIR /go/src/github.com/ddddddO/tag-mng
COPY cmd/api/main.go go.* ./
COPY api/ api/
COPY domain/ domain/
COPY repository/ repository/

RUN go build -o go-api main.go

FROM alpine

#RUN apk add --no-cache --repository http://dl-cdn.alpinelinux.org/alpine/v3.11/main ca-certificates
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/ddddddO/tag-mng/go-api /api

ENV PORT="8080"
CMD ["/api"]
