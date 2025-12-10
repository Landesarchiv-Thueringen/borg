FROM golang:alpine3.23 AS build
ARG BORG_VERSION=${BORG_VERSION}
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o magika_api -ldflags "-X main.version=${BORG_VERSION}" ./cmd

FROM python:3.13-slim AS prod
WORKDIR /borg/tools/magika
COPY --from=build /build/magika_api .
RUN pip install magika
ENTRYPOINT ["./magika_api"]