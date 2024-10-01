FROM golang:alpine3.20 AS build
WORKDIR /borg/tools/droid
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o droidapi cmd/droid_api.go

FROM alpine:3.20 AS prod
WORKDIR /borg/tools/droid
RUN apk update
RUN apk add unzip openjdk17
COPY --from=build /borg/tools/droid/ .
RUN unzip third_party/droid-binary-6.8.0-bin.zip -d third_party/
RUN chmod +x third_party/droid.sh
# Droid is prone to failure when being called for the first time, especially for
# concurrent calls. We call it one time when creating the image, so it sets up
# all its files and will be save to call when the container starts.
RUN /bin/ash \
    third_party/droid.sh \
    -Ns third_party/DROID_SignatureFile_V114.xml \
    -Nc third_party/container-signature-20230822.xml \
    third_party/LICENSE
CMD ["./droidapi"]