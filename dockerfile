FROM golang as builder
WORKDIR /build/api
COPY go.mod ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build server.go
FROM alpine
WORKDIR /root
COPY --from=builder /build/api/server .
EXPOSE 8080
CMD ["./server" , "db"]