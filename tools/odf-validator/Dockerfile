FROM golang:alpine3.21 AS build
ARG BORG_VERSION=${BORG_VERSION}
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o odf_validator_api -ldflags "-X main.version=${BORG_VERSION}" ./cmd

FROM alpine:3.21 AS prod
RUN apk add --no-cache openjdk21
WORKDIR /app
COPY --from=build /app/odf_validator_api .
COPY third_party third_party
CMD ["./odf_validator_api"]