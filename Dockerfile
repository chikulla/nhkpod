FROM golang:1.18-alpine
RUN apk add --no-cache ffmpeg

WORKDIR /app
EXPOSE 8080

COPY go.mod ./
COPY go.sum ./
COPY *.go ./
COPY cmd ./cmd
COPY conf.yml ./
COPY assets ./assets

RUN go build ./cmd/nhkpod/...

CMD ["./nhkpod"]
