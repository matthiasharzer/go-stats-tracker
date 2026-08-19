FROM golang:1.26.5-alpine3.24 AS build

ARG version=unknown

RUN apk update && \
		apk add git

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download && \
		go mod verify

COPY . .

RUN module_path=$(go list -m) && \
	go build \
		-o /go/bin/go-stats-tracker \
		-ldflags "-X ${module_path}/cmd/version.version=$version" \
		.

FROM alpine:3.24

RUN addgroup -S app && adduser -S -G app app

COPY --from=build /go/bin/go-stats-tracker /usr/local/bin/go-stats-tracker

WORKDIR /var/lib/go-stats-tracker
RUN chown app:app /var/lib/go-stats-tracker

USER app

ENTRYPOINT ["/usr/local/bin/go-stats-tracker"]
