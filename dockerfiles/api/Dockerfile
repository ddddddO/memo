FROM golang:1.18 as builder
WORKDIR /go/src/github.com/ddddddO/memo
COPY cmd/api/main.go go.* ./
COPY api/ api/
COPY domain/ domain/
COPY repository/ repository/
RUN go build -o rest-api main.go

FROM gcr.io/distroless/base
COPY --from=builder /go/src/github.com/ddddddO/memo/rest-api ./
ENV PORT="8080"
CMD ["/rest-api"]
