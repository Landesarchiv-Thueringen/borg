FROM golang:alpine3.20 AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o ooxml_validator_api cmd/ooxml_validator_api.go

FROM alpine:3.20 AS extract
RUN apk update && apk add unzip
COPY third_party /third_party
WORKDIR /third_party
RUN unzip linux-x64.zip
RUN chmod +x OOXMLValidatorCLI

FROM alpine:3.20 AS prod
WORKDIR /app
RUN apk update
RUN apk add unzip libgcc libstdc++ gcompat icu-libs
COPY --from=build /build/ooxml_validator_api .
COPY --from=extract /third_party/OOXMLValidatorCLI ./third_party/
CMD ["./ooxml_validator_api"]