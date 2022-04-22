FROM golang:1.17.5-alpine AS builder

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server cmd/api/main.go

FROM alpine:3.15.4

WORKDIR /api

COPY --from=builder /go/src/server .

EXPOSE 3000
CMD [ "./server" ]