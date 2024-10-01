FROM golang:alpine3.20 AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o droidapi cmd/droid_api.go

FROM alpine:3.20 AS extract
RUN apk update && apk add unzip
COPY third_party /third_party
WORKDIR /third_party
RUN unzip droid-binary-6.8.0-bin.zip
RUN rm droid-binary-6.8.0-bin.zip
RUN chmod +x droid.sh

FROM alpine:3.20 AS prod
RUN apk update && apk add openjdk17
WORKDIR /borg/tools/droid
COPY --from=build /build/droidapi .
COPY --from=extract /third_party ./third_party
# Droid is prone to failure when being called for the first time, especially for
# concurrent calls. We call it one time when creating the image, so it sets up
# all its files and will be save to call when the container starts.
RUN /bin/ash \
    third_party/droid.sh \
    -Ns third_party/DROID_SignatureFile_V114.xml \
    -Nc third_party/container-signature-20230822.xml \
    third_party/LICENSE
CMD ["./droidapi"]