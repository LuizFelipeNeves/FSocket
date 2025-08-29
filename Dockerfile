# Dockerfile para Socket.io server em Node.js
FROM node:latest
WORKDIR /app
COPY package*.json ./
RUN npm install --production
COPY . .
EXPOSE 8080
CMD ["node", "server.js"]
