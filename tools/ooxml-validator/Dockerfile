FROM golang:alpine3.23 AS build
ARG BORG_VERSION=${BORG_VERSION}
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o ooxml_validator_api -ldflags "-X main.version=${BORG_VERSION}" ./cmd

FROM alpine:3.23 AS extract
RUN apk add --no-cache unzip
COPY third_party /third_party
WORKDIR /third_party
RUN unzip linux-x64.zip
RUN chmod +x OOXMLValidatorCLI

FROM alpine:3.23 AS prod
WORKDIR /app
RUN apk add --no-cache unzip libgcc libstdc++ gcompat icu-libs
COPY --from=build /build/ooxml_validator_api .
COPY --from=extract /third_party/OOXMLValidatorCLI ./third_party/
CMD ["./ooxml_validator_api"]