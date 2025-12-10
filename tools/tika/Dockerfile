FROM golang:alpine3.23 AS build
ARG BORG_VERSION=${BORG_VERSION}
WORKDIR /borg/tools/tika
COPY go.mod go.sum ./
RUN go mod download
COPY cmd cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o tika_api -ldflags "-X main.version=${BORG_VERSION}" ./cmd

FROM alpine:3.23 AS prod
WORKDIR /borg/tools/tika
RUN apk add --no-cache openjdk17
COPY --from=build /borg/tools/tika/tika_api .
COPY third_party third_party
CMD ["./tika_api"]