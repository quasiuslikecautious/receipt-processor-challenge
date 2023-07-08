# syntax=docker/dockerfile:1

FROM golang:1.19

WORKDIR /workdir

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /receipt-processor-challenge

EXPOSE 8080

CMD ["/receipt-processor-challenge"]
