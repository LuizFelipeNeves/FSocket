# Arquitetura SSE Multi-Cliente Multi-Canal

## Visão Geral

Este projeto implementa um servidor SSE (Server-Sent Events) que permite múltiplos clientes se conectarem a múltiplos canais, além de um endpoint para publicação de mensagens em canais específicos.

---

## Componentes

- **Servidor SSE**: Gerencia conexões e distribui eventos para clientes.
- **Clientes SSE**: Aplicações que recebem eventos via SSE.
- **Gerenciador de Canais**: Mantém a relação canal → clientes conectados.
- **Endpoint de Publicação**: API REST para envio de mensagens aos canais.

---

## Fluxo de Funcionamento

1. **Conexão do Cliente**
   - Cliente faz GET `/sse/:canal`.
   - Servidor adiciona conexão à lista do canal.
   - Eventos são enviados conforme publicados.

2. **Inscrição em Canais**
   - Cliente pode se conectar a múltiplos canais (ex: `/sse/news`, `/sse/chat`).
   - Servidor mantém mapa de canais para clientes conectados.

3. **Publicação de Mensagens**
   - POST `/publish/:canal` recebe mensagem.
   - Mensagem é distribuída para todos clientes conectados ao canal.

---

## Estrutura de Dados

- **Canais**: `map[string][]Client`
- **Mensagens**: `{ canal, mensagem, timestamp }`

---

## Endpoints

- `GET /sse/:canal` — Conecta cliente ao canal SSE.
- `POST /publish/:canal` — Publica mensagem para o canal.

---

## Tecnologias Sugeridas

- **Backend**: Go (eficiente em memória e concorrência)
- **Frontend**: Qualquer cliente com suporte a EventSource.
- **Armazenamento**: Opcional, para histórico de mensagens (Redis, etc).

---

## Escalabilidade

- Utilizar goroutines para conexões concorrentes (Go).
- Balanceamento de carga para múltiplos servidores.
- Persistência opcional via Redis.

---

## Segurança

- Autenticação nos endpoints de publicação.
- Controle de acesso por canal.

---

## Observações sobre Go

- Go utiliza goroutines leves, reduzindo uso de memória em aplicações de rede.
- Gerenciamento eficiente de conexões e baixo overhead.
- Ideal para sistemas com muitos clientes simultâneos.

---

## Sugestão de Estrutura de Pastas

```
/
├── main.go
├── channels/
├── handlers/
├── README.md
```

---

## Próximos Passos
- Desenvolver interface web simples para estatísticas gerais de uso.

---

## Documentação das Funcionalidades

### Endpoints

- `GET /sse?canal=<nome>`: Conecta o cliente ao canal SSE especificado.
- `POST /publish?canal=<nome>&msg=<mensagem>`: Publica uma mensagem para todos os clientes conectados ao canal. Requer header `Authorization: Bearer <token>`.
- `POST /publish?canal=broadcast&msg=<mensagem>`: Publica mensagem para todos os canais simultaneamente.
- `GET /health`: Retorna status do servidor.
- `GET /stats`: Retorna estatísticas gerais (canais ativos, clientes conectados, mensagens publicadas).

### Funcionalidades

- Múltiplos canais SSE simultâneos.
- Autenticação por token fixo no endpoint de publicação.
- Expiração automática de canais inativos.
- Canal de broadcast para envio global.
- Logs de eventos (conexão, publicação, desconexão).
- Health check endpoint.
- Suporte a CORS.
- Interface web simples para estatísticas gerais.

---

## Guia Básico de Uso com React e Next.js
---

## Guia de Execução com Docker

1. Certifique-se de ter o Docker instalado.
2. Construa a imagem:
   ```sh
   docker build -t sse-server .
   ```
3. Rode o container (opcionalmente defina o token de autenticação):
   ```sh
   docker run -p 8080:8080 -e AUTH_TOKEN=seu_token_fixo sse-server
   ```
4. O servidor estará disponível em `http://localhost:8080`.
5. O token de autenticação pode ser definido pela variável de ambiente `AUTH_TOKEN`.

---

### Recebendo eventos SSE

```jsx
// Exemplo de componente React para receber eventos SSE
import { useEffect, useState } from 'react';

export default function SSEClient({ canal }) {
   const [mensagens, setMensagens] = useState([]);

   useEffect(() => {
      const es = new EventSource(`/sse?canal=${canal}`);
      es.onmessage = (e) => {
         setMensagens(msgs => [...msgs, e.data]);
      };
      es.onerror = () => {
         es.close();
         // Implementar reconexão automática se necessário
      };
      return () => es.close();
   }, [canal]);

   return (
      <ul>
         {mensagens.map((msg, i) => <li key={i}>{msg}</li>)}
      </ul>
   );
}
```

### Publicando mensagens (Next.js API ou fetch)

```js
// Exemplo de publicação de mensagem
fetch('/publish?canal=meucanal&msg=Olá', {
   method: 'POST',
   headers: {
      'Authorization': 'Bearer seu_token_fixo_aqui'
   }
});
```

### Exibindo estatísticas

```js
// Exemplo de consulta às estatísticas
fetch('/stats')
   .then(res => res.json())
   .then(data => console.log(data));
```

---


---

## Ideias de Aprimoramento

- **Reconexão automática**: Lógica no cliente para reconectar em caso de desconexão.
- **Logs de eventos**: Registrar conexões, publicações e desconexões para monitoramento.
- **Canal de broadcast**: Permitir envio de mensagens para todos os canais simultaneamente.
- **Expiração de canais**: Remover canais inativos após um tempo sem clientes conectados.
- **Health check endpoint**: Adicionar endpoint para monitorar status do servidor.
- **Suporte a CORS**: Permitir conexões de diferentes origens.
- **Documentação da API**: Criar um arquivo de especificação dos endpoints.
- **Interface simples de stats de uso geral**: Página web mostrando número de canais ativos, clientes conectados e mensagens publicadas.

---

## Próximos Passos

- Implementar esqueleto do servidor em Go.
- Definir modelos de dados e handlers dos endpoints.
- Adicionar autenticação e controle de acesso.
- Adicionar logs de eventos.
- Implementar canal de broadcast.
- Implementar expiração automática de canais inativos.
- Adicionar endpoint de health check.
- Habilitar CORS.
- Criar documentação da API.
- Desenvolver interface web simples para estatísticas gerais de uso.
