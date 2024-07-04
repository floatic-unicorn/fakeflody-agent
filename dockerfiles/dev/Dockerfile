FROM golang:1.22 AS builder

WORKDIR /build

COPY ../../ .

RUN go build ./cmd/fakeflody

FROM alpine:3.20

WORKDIR /app

RUN apk add libc6-compat

COPY --from=builder /build/fakeflody .

EXPOSE 8080

ENTRYPOINT ["./fakeflody", "--t", "2", "--web"]
