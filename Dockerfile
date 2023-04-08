FROM golang:1.19 AS builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -o little-john-store .

FROM alpine:latest AS production
COPY --from=builder /app .
EXPOSE 8080
CMD ["./little-john-store"]