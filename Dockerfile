# Stats service

FROM golang:1.23.4 AS production-env

WORKDIR /app

COPY . ./



RUN #go env -w GOPROXY=https://goproxy.cn,https://gocenter.io,https://goproxy.io,direct
RUN go env -w GOPROXY=https://proxy.golang.org,direct

RUN go mod download


CMD ["go", "run", "/app/cmd/app/main.go"]

