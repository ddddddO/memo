FROM golang:1.13-alpine as builder

WORKDIR /go/src/github.com/ddddddO/memo
COPY cmd/api/main.go go.* ./
COPY api/ api/
COPY domain/ domain/
COPY repository/ repository/

RUN go build -o rest-api main.go

FROM alpine

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/ddddddO/memo/rest-api /api

ENV PORT="8080"
CMD ["/api"]
