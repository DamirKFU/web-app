FROM golang:1.24-bullseye

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=1
RUN go build -o main .

CMD ["./main"]