FROM golang:1.24.6 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./main.go

FROM scratch

COPY --from=builder /out/app /app

EXPOSE 8080

ENTRYPOINT ["/app"]
