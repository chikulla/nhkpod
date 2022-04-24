FROM jrottenberg/ffmpeg:4.1-alpine
FROM golang:1.18-alpine

COPY --from=0 / /

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
