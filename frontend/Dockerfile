# --- build stage ---
FROM node:20-alpine AS builder
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

# Vite bakes VITE_* env vars directly into the built JS as literal
# strings AT BUILD TIME — there is no "runtime config" for a static
# site. That means these have to arrive as build ARGs (set via
# docker-compose.yml's `build.args`), not as normal container
# environment variables, which only affect a running process and
# would have no effect on files that were already built.
ARG VITE_BACKEND_URL
ARG VITE_PROMETHEUS_URL
ARG VITE_GRAFANA_URL
ENV VITE_BACKEND_URL=$VITE_BACKEND_URL
ENV VITE_PROMETHEUS_URL=$VITE_PROMETHEUS_URL
ENV VITE_GRAFANA_URL=$VITE_GRAFANA_URL

RUN npm run build

# --- final stage ---
# The build stage's entire node_modules/toolchain is discarded here.
# All that survives is /app/dist — plain HTML/CSS/JS — served by
# nginx, a tiny web server with nothing Node-specific about it.
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80
