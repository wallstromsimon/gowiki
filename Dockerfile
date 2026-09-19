FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o wiki .

FROM alpine:3.21

RUN apk --no-cache add ca-certificates \
	&& adduser -D -g '' appuser

WORKDIR /home/appuser

COPY --from=builder /app/wiki .

USER appuser

EXPOSE 8080

CMD ["./wiki"]
