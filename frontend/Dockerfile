FROM node:18-alpine AS base
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci

FROM base AS dev
RUN apk add --no-cache nginx supervisor
COPY nginx.dev.conf /etc/nginx/http.d/default.conf
COPY supervisord.conf /etc/supervisord.conf
RUN mkdir -p /var/lib/nginx /var/log/nginx /run/nginx && \
    chown -R node:node /app /var/lib/nginx /var/log/nginx /run/nginx && \
    chown -R nginx:nginx /var/lib/nginx
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]

FROM base AS builder
COPY . .
RUN npm run build

FROM nginx:stable-alpine AS prod
# Копируем собранную статику и production-конфиг nginx
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.prod.conf /etc/nginx/conf.d/default.conf
CMD ["nginx", "-g", "daemon off;"]