FROM golang:1.17 AS builder
RUN mkdir /app 
COPY . /app
WORKDIR /app 
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd
EXPOSE 5000


FROM alpine:3.10 AS production
COPY --from=builder /app .
CMD ["./app/server", "--addr=:5000"]