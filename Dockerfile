FROM golang:1.17 
RUN mkdir /app 
COPY . /app
WORKDIR /app 
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/
EXPOSE 5000