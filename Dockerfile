FROM golang:1.19 AS builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build

FROM alpine:latest AS production
COPY --from=builder /app .
CMD ["./little-john-store"]