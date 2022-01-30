FROM golang:1.17 AS builder

WORKDIR /src
ADD . /src
RUN go mod download
RUN CGO_ENABLED=0 go build -a -ldflags '-s -w -extldflags "-static"' -o /s3backup .

FROM scratch
COPY --from=builder /etc/mime.types /etc/mime.types
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /s3backup /s3backup
ENTRYPOINT ["/s3backup"]



