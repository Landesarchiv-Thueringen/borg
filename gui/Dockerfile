FROM node:24-alpine3.23 AS build

RUN corepack enable
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
WORKDIR /build
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
COPY angular.json tsconfig*.json ./
COPY src ./src
COPY public ./public
ARG CONFIGURATION
RUN pnpm run build --configuration $CONFIGURATION

FROM nginx:1.29-alpine3.23 AS prod
COPY --from=build /build/dist/borg-gui /app/
COPY build/nginx/default.conf.template /etc/nginx/templates/default.conf.template
