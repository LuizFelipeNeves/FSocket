# Socket.io Multi-Canal Server

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

### 3. Inicie o servidor
```sh
node server.js
```

### 4. Usando Docker
```sh
docker build -t socketio-server .
docker run -p 8080:8080 -e AUTH_TOKEN=seu_token socketio-server
```

## Endpoints

- `POST /publish/:canal` — Publica mensagem em um canal (body: `{ msg: "texto" }`, header: `Authorization: Bearer <token>`)
- `POST /publish/broadcast` — Publica mensagem para todos os canais
- `GET /stats` — Estatísticas gerais
- `GET /health` — Status do servidor

## Exemplo de uso com Socket.io Client

```js
import { io } from "socket.io-client";
const socket = io("http://localhost:8080");
socket.emit('join', 'meucanal');
socket.on('mensagem', (data) => {
  console.log(data);
});
```

## Publicando mensagens

```sh
curl -X POST http://localhost:8080/publish/meucanal \
  -H "Authorization: Bearer seu_token" \
  -H "Content-Type: application/json" \
  -d '{"msg":"Olá"}'
```

---

Para dúvidas ou melhorias, abra uma issue ou contribua!
