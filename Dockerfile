FROM golang:1.26.5-alpine AS build

ARG version=unknown

RUN apk add --no-cache \
    build-base \
    tesseract-ocr-dev

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download && \
		go mod verify

COPY . .

RUN module_path=$(go list -m) && \
    CGO_ENABLED=1 go build \
										-trimpath \
										-o /go/bin/go-stats-tracker \
										-ldflags "-X ${module_path}/cmd/version.version=$version" \
										.

FROM alpine:3.24

RUN apk add --no-cache \
    tesseract-ocr \
    tesseract-ocr-data-eng \
    tesseract-ocr-data-deu

RUN addgroup -g 1000 app && adduser -u 1000 -G app -D app

COPY --from=build /go/bin/go-stats-tracker /usr/local/bin/go-stats-tracker

WORKDIR /var/lib/go-stats-tracker
RUN mkdir -p /var/lib/go-stats-tracker/data && chown -R app:app /var/lib/go-stats-tracker

USER app

ENTRYPOINT ["/usr/local/bin/go-stats-tracker"]
