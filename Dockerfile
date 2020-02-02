FROM "golang:latest" AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY main.go ./

RUN GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -a -ldflags '-s' -o spot-look-back

FROM scratch

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /app/spot-look-back ./

ENTRYPOINT ["./spot-look-back"]
