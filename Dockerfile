FROM golang:1.15-buster as builder
WORKDIR /haaukins

COPY . .
RUN go build -o server .

FROM gcr.io/distroless/base-debian10
COPY --from=builder /haaukins /
EXPOSE 9090
CMD ["/server"]