ARG GO_VERSION=1
FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -v -o /pb-ircclient ./cmd/ircclient
RUN CGO_ENABLED=0 go build -v -o /pb-server ./cmd/server

FROM alpine:latest as pb-ircclient
COPY --from=builder /pb-ircclient /usr/local/bin/
CMD ["pb-ircclient"]

FROM alpine:latest as pb-server
RUN apk add --no-cache tzdata
COPY --from=builder /pb-server /usr/local/bin/
CMD ["pb-server"]
