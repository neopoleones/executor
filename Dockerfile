FROM golang:1.21

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o executor cmd/executor/main.go

EXPOSE 8080
CMD ["./executor"]