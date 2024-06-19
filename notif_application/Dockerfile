FROM golang:1.22.4-alpine
LABEL authors="armin"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o main


EXPOSE 8001

CMD ["./main"]