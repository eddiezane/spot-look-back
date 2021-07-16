FROM --platform=$BUILDPLATFORM "golang:latest" AS builder

ARG BUILDPLATFORM

ARG TARGETPLATFORM

ENV GOPROXY=direct

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY main.go ./

RUN case ${TARGETPLATFORM} in \
         "linux/amd64")  GOOS=linux GOARCH=amd64 ;; \
         "linux/arm64")  GOOS=linux GOARCH=arm64 ;; \
         "linux/arm/v7") GOOS=linux GOARCH=arm GOARM=7 ;; \
         "linux/arm/v6") GOOS=linux GOARCH=arm GOARM=6 ;; \
         *) exit 1 ;; \
    esac \
    && CGO_ENABLED=0 go build -a -ldflags '-s' -o spot-look-back

FROM scratch

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /app/spot-look-back ./

ENTRYPOINT ["./spot-look-back"]
