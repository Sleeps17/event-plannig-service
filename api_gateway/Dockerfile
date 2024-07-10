FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git bash make gcc musl-dev upx
RUN go install github.com/go-task/task/v3/cmd/task@latest

WORKDIR /go/src/api-gateway

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o /go/bin/api-gateway ./cmd/api_gateway/main.go
RUN upx -9 /go/bin/api-gateway

FROM alpine AS runner

RUN apk add --no-cache bash

COPY --from=builder /go/bin/task /usr/local/bin/task
COPY Taskfile.yaml ./

COPY --from=builder /go/bin/api-gateway /usr/local/bin/api-gateway
COPY /configs/local.yaml /configs/local.yaml
ENV CONFIG_PATH=/configs/local.yaml

EXPOSE 8080

ENTRYPOINT ["task", "run"]