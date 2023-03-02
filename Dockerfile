FROM alpine:latest
WORKDIR /app
COPY ./go_temp /app
CMD ["/app/go_temp"]