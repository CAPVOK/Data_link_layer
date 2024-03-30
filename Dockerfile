FROM golang:1.21.2-alpine 

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o /main

WORKDIR /app

CMD ["/main"]