FROM golang:1.22.4-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o main


EXPOSE 8000

CMD ["./main"]