FROM golang:1.18 AS builder
RUN mkdir /app 
COPY . /app
WORKDIR /app 
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/
EXPOSE 5000

FROM alpine:latest AS production
COPY --from=builder /app .
CMD ["./server"]