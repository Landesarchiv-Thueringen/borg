FROM golang:alpine AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o ooxml_validator_api cmd/ooxml_validator_api.go

FROM alpine:3.19 AS prod
WORKDIR /app
RUN apk update
RUN apk add libgcc libstdc++ gcompat icu-libs
COPY --from=build /build/ooxml_validator_api .
COPY third_party third_party
CMD ["./ooxml_validator_api"]