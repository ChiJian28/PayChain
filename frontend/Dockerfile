FROM node:20-alpine AS builder
RUN apk add --no-cache g++ make python3
WORKDIR /app
COPY .npmrc ./
COPY package.json package-lock.json* pnpm-lock.yaml* yarn.lock* ./
RUN npm ci --prefer-offline --no-audit --progress=false || npm i --no-audit --progress=false
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx","-g","daemon off;"]

