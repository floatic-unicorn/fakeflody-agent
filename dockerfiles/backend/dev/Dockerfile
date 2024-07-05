FROM golang:1.22 AS builder

WORKDIR /build

COPY ../../../go.mod ../../../go.sum .

COPY ../../../docs ./docs

COPY ../../../src ./src

RUN go build ./src/cmd/fakeflody

FROM alpine:3.20

RUN apk add libc6-compat

WORKDIR /app

COPY --from=builder /build/fakeflody .

EXPOSE 8080

ENTRYPOINT ["./fakeflody", "--t", "2", "--web"]
