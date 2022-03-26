FROM golang:1.18 as builder
WORKDIR /go/src/github.com/ddddddO/memo
COPY cmd/api/main.go go.* ./
COPY api/ api/
COPY models/ models/
COPY repository/postgres/ repository/postgres/
RUN go build -o rest-api main.go

FROM gcr.io/distroless/base
COPY --from=builder /go/src/github.com/ddddddO/memo/rest-api ./
ENV PORT="8080"
CMD ["/rest-api"]
