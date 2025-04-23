ARG GO_VERSION=1
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /usr/src/app
RUN apk add --no-cache protoc
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN make generate
RUN CGO_ENABLED=0 go build -v -o /pb-ircclient ./cmd/ircclient
RUN CGO_ENABLED=0 go build -v -o /pb-server ./cmd/server

FROM alpine:latest AS pb-ircclient
COPY --from=builder /pb-ircclient /usr/local/bin/
CMD ["pb-ircclient"]

FROM alpine:latest AS pb-server
RUN apk add --no-cache tzdata
COPY --from=builder /pb-server /usr/local/bin/
CMD ["pb-server"]
