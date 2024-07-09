FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git bash make gcc musl-dev upx
RUN go install github.com/pressly/goose/v3/cmd/goose@latest &&  \
    go install github.com/go-task/task/v3/cmd/task@latest

WORKDIR /go/src/rooms-service

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o /go/bin/rooms-service ./cmd/rooms/main.go
RUN upx -9 /go/bin/rooms-service

FROM alpine AS runner

RUN apk add --no-cache bash

COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY ./migrations ./migrations

COPY --from=builder /go/bin/task /usr/local/bin/task
COPY Taskfile.yaml ./

COPY --from=builder /go/bin/rooms-service /usr/local/bin/rooms-service
COPY /configs/local.yaml /configs/local.yaml
ENV CONFIG_PATH=/configs/local.yaml

EXPOSE 4404

ENTRYPOINT ["task", "migrate-and-run"]