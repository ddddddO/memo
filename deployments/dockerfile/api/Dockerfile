FROM golang:1.13-alpine as builder

WORKDIR /go/src/github.com/ddddddO/tag-mng
COPY cmd/api/main.go go.* ./
COPY api/ api/
COPY domain/ domain/

RUN go build -o go-api main.go

FROM alpine

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/ddddddO/tag-mng/go-api /api

ENV PORT="8080"
CMD ["/api"]
