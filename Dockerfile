FROM golang:latest AS builder
ADD . /app
WORKDIR /app

RUN  go build -a -o /main .


FROM alpine:latest

COPY --from=builder /main /app/main
RUN chmod +x /app/main
CMD ["/app/main"]

EXPOSE 8086