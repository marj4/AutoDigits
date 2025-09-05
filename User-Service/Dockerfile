FROM golang:latest

WORKDIR /user-service

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o user-service ./cmd

CMD ["./user-service"]
