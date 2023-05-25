FROM golang:latest as builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN go build -o server ./main.go
FROM debian:buster-slim
COPY --from=builder /app/server /app/server
EXPOSE 50051
CMD ["/app/server"]