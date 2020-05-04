FROM golang:1.13-alpine as builder
 
WORKDIR /go/src/api
COPY deployments/files/api .

# RUN go get -d -v ./...
# RUN go install -v ./...
# RUN go build -o api main.go && \
#     mv api $GOPATH/bin
RUN go build -o api main.go

FROM alpine

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/api/api /api

#ENV DBDSN="host=localhost dbname=tag-mng user=postgres password=postgres sslmode=disable"
ENV DBDSN="XXX"
ENV PORT=8080
CMD ["/api"]
