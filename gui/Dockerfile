FROM node:22-alpine3.20 AS build

RUN corepack enable
# next command can be removed if issue: https://github.com/pnpm/pnpm/issues/9029 was fixed
RUN corepack prepare pnpm@10.0.0 --activate
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
WORKDIR /build
COPY package.json pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
COPY angular.json tsconfig*.json ./
COPY src ./src
COPY public ./public
ARG CONFIGURATION
RUN pnpm run build --configuration $CONFIGURATION

FROM nginx:1.27-alpine3.20 AS prod
COPY --from=build /build/dist/borg-gui /app/
COPY build/nginx/default.conf.template /etc/nginx/templates/default.conf.template
