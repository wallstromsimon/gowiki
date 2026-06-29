FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o wiki .

FROM alpine:3.20

RUN apk --no-cache add ca-certificates \
	&& adduser -D -g '' appuser

WORKDIR /home/appuser

COPY --from=builder /app/wiki .
COPY --from=builder /app/html_templates html_templates

USER appuser

EXPOSE 8080

CMD ["./wiki"]
