FROM golang:alpine3.23 AS build
ARG BORG_VERSION=${BORG_VERSION}
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o mediainfo_api -ldflags "-X main.version=${BORG_VERSION}" ./cmd

FROM alpine:3.23 AS prod
WORKDIR /borg/tools/mediainfo
COPY third_party/de.csv .
COPY --from=build /build/mediainfo_api .
RUN apk add --no-cache mediainfo
ENTRYPOINT ["./mediainfo_api"]