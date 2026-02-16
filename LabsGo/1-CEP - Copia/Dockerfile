FROM golang:1.24-alpine

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./main.go 

CMD ["/out/app"]
