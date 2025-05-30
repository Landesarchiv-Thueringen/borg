FROM golang:alpine3.21 AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o droidapi cmd/droid_api.go

FROM alpine:3.21 AS extract
RUN apk add --no-cache unzip
COPY third_party /third_party
WORKDIR /third_party
RUN unzip droid-binary-6.8.1-bin.zip
RUN rm droid-binary-6.8.1-bin.zip
RUN chmod +x droid.sh

FROM alpine:3.21 AS prod
RUN apk add --no-cache openjdk21
WORKDIR /borg/tools/droid
COPY --from=build /build/droidapi .
COPY --from=extract /third_party ./third_party
# Droid is prone to failure when being called for the first time, especially for
# concurrent calls. We call it one time when creating the image, so it sets up
# all its files and will be save to call when the container starts.
RUN /bin/ash \
    third_party/droid.sh \
    -Ns third_party/DROID_SignatureFile_V120.xml \
    -Nc third_party/container-signature-20240715.xml \
    third_party/LICENSE
CMD ["./droidapi"]