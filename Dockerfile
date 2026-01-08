# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod tidy && go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-w -s" \
    -a \
    -installsuffix cgo \
    -o /app/dist/abb-exporter ./cmd/abb-exporter

FROM gcr.io/distroless/static-debian13:nonroot
ARG TARGETARCH

ENV TZ=Europe/Rome

WORKDIR /app
COPY --from=builder /app/dist/abb-exporter /app/abb-exporter
COPY ./config.prod.yml /app/config.yml

EXPOSE 8080

ENTRYPOINT ["/app/abb-exporter"]