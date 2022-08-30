FROM golang:1.17.4-buster as builder
WORKDIR /app

COPY go.* /
RUN go mod download
COPY . .
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o webapp ./cmd/api/main.go

FROM alpine:3.12.0
RUN apk add --no-cache tzdata
ENV TZ Asia/Bangkok
ARG SUB_PROJECT
ENV SERVICE_NAME $SUB_PROJECT
WORKDIR /app
COPY --from=builder app/webapp .
EXPOSE 3000
ENTRYPOINT ["./webapp"]
