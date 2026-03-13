# FSocket - Servidor SSE em Go

Servidor de notifications em tempo real usando Server-Sent Events (SSE) e Go.

## Funcionalidades

- Server-Sent Events (SSE) para推送 de mensagens em tempo real
- Múltiplos canais (rooms) para clientes
- Publicação de mensagens via endpoint REST
- Broadcast para todos os canais
- Autenticação por token via variável de ambiente `AUTH_TOKEN`
- Estatísticas de uso (`/stats`)
- Health check (`/health`)
- Alta performance com Goroutines

## Como rodar

### Build

```sh
go build -o fsocket ./cmd/server
```

### Executar

```sh
./fsocket
```

Ou com variáveis de ambiente:

```sh
PORT=8080 AUTH_TOKEN=seu_token ./fsocket
```

### Docker

```sh
docker build --build-arg AUTH_TOKEN=$AUTH_TOKEN -t fsocket .
docker run -p 8080:8080 -e AUTH_TOKEN=$AUTH_TOKEN fsocket
```

## Endpoints

- `GET /sse?channel=xxx` — Conexão SSE para um canal
- `POST /publish` — Publicar mensagem em um canal
- `POST /publish/broadcast` — Broadcast para todos os canais
- `GET /stats` — Estatísticas (canais ativos, clientes conectados, mensagens publicadas)
- `GET /health` — Health check

## Cliente (Frontend)

### Conexão SSE

```js
const eventSource = new EventSource('http://localhost:8080/sse?channel=store_123')

eventSource.addEventListener('new_order', (event) => {
  const data = JSON.parse(event.data)
  console.log('Novo pedido:', data)
})

eventSource.addEventListener('order_updated', (event) => {
  const data = JSON.parse(event.data)
  console.log('Pedido atualizado:', data)
})

eventSource.addEventListener('order_edited', (event) => {
  const data = JSON.parse(event.data)
  console.log('Pedido editado:', data)
})
```

### Publicar mensagem

```sh
curl -X POST http://localhost:8080/publish \
  -H "Authorization: Bearer seu_token" \
  -H "Content-Type: application/json" \
  -d '{"channel":"store_123","msg":"Novo pedido","eventType":"new_order","extra":{"orderId":"123","orderNumber":456}}'
```

## Variáveis de Ambiente

| Variável | Padrão | Descrição |
|----------|--------|-----------|
| PORT | 8080 | Porta do servidor |
| AUTH_TOKEN | - | Token para autenticação nas rotas de publish |

## Arquitetura

```
cmd/server/main.go       # Entry point
internal/
├── config/              # Configurações
├── hub/                 # Gerenciador de clientes/canais
├── handler/             # HTTP handlers
└── middleware/          # Auth middleware
```
