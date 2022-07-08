FROM golang:1.17
RUN mkdir /app 
COPY . /app
WORKDIR /app 
RUN go build -o server ./cmd
EXPOSE 5000
CMD ["/app/server"]
