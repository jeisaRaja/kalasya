FROM golang:1.23-alpine

WORKDIR /app

RUN go install github.com/air-verse/air@v1.52.3

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o cmd/web/kalasya cmd/web/*.go

EXPOSE 8000

CMD ["air", "-c", ".air.toml"]
