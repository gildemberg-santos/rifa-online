FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM node:22-alpine AS frontend-builder
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

FROM nginx:stable-alpine
RUN apk add --no-cache ca-certificates
COPY --from=backend-builder /server /server
COPY --from=frontend-builder /app/dist /usr/share/nginx/html
COPY nginx.unified.conf /etc/nginx/conf.d/default.conf
COPY start.sh /start.sh
RUN chmod +x /start.sh
EXPOSE 80
ENTRYPOINT ["/start.sh"]
