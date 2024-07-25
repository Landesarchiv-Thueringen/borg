FROM node:22-alpine3.20 AS build

WORKDIR /build
COPY package.json package-lock.json ./
RUN npm ci
COPY angular.json tsconfig*.json ./
COPY src ./src
ARG CONFIGURATION
RUN npm run build -- --configuration $CONFIGURATION

FROM nginx:1.27-alpine3.19 AS prod
COPY --from=build /build/dist/borg-gui /app/
COPY build/nginx/default.conf /etc/nginx/conf.d/default.conf
