FROM golang:1.23 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o mock main.go

FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/mock /app/mock

CMD ["/app/mock"]
