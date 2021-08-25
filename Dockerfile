FROM golang:1.15-alpine as builder
WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o esm -a -ldflags '-w -extldflags "-static"' .

FROM scratch
COPY --from=builder /app/esm /
EXPOSE 50095
CMD ["/esm"]