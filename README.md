# Socket.io Multi-channel Server

Este projeto implementa um servidor de eventos em tempo real usando Socket.io e Node.js, com suporte a múltiplos canais, autenticação por token e endpoints REST para publicação de mensagens.

## Funcionalidades

- Múltiplos canais (salas) para clientes
- Publicação de mensagens via endpoint REST
- Broadcast para todos os canais
- Autenticação por token via variável de ambiente `AUTH_TOKEN`
- Estatísticas de uso (`/stats`)
- Health check (`/health`)
- Suporte a CORS

## Como rodar

### 1. Instale as dependências
```sh
npm install
```

### 2. Defina o token de autenticação

No ambiente ou ao rodar o container:
```sh
export AUTH_TOKEN=seu_token
```


# Socket.io Multi-Channel Server

This project implements a real-time event server using Socket.io and Node.js, supporting multiple channels, token authentication, and REST endpoints for message publishing.

## Features

- Multiple channels (rooms) for clients
- Message publishing via REST endpoint
- Broadcast to all channels
- Token authentication via `AUTH_TOKEN` environment variable
- Usage statistics (`/stats`)
- Health check (`/health`)
- CORS support

## How to run

### 1. Install dependencies
```sh
npm install
```

### 2. Set the authentication token

In your environment or when running the container:
```sh
export AUTH_TOKEN=your_token
```

### 3. Start the server
```sh
node server.js
```

### 4. Using Docker
```sh
docker build -t socketio-server .
docker run -p 8080:8080 -e AUTH_TOKEN=your_token socketio-server
```

## Endpoints

- `POST /publish` — Publish a message to a channel (body: `{ channel: "name", msg: "text" }`, header: `Authorization: Bearer <token>`)
- `POST /publish/broadcast` — Publish a message to all channels
- `GET /stats` — General statistics
- `GET /health` — Server status

## Example usage with Socket.io Client

```js
import { io } from "socket.io-client";
const socket = io("http://localhost:8080");
socket.emit('join', 'mychannel');
socket.on('mensagem', (data) => {
  console.log(data);
});
```

## Publishing messages

```sh
curl -X POST http://localhost:8080/publish \
  -H "Authorization: Bearer your_token" \
  -H "Content-Type: application/json" \
  -d '{"channel":"mychannel","msg":"Hello"}'
```

---

For questions or improvements, open an issue or contribute!
