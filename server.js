const express = require('express');
const http = require('http');
const { Server } = require('socket.io');

const AUTH_TOKEN = process.env.AUTH_TOKEN || 'seu_token_fixo_aqui';

const app = express();
const server = http.createServer(app);
const io = new Server(server, {
  cors: {
    origin: '*',
  },
});

let messagesPublished = 0;
let channels = {};

// Middleware de autenticação para publicação
function authMiddleware(req, res, next) {
  const token = req.headers['authorization'];
  if (token !== `Bearer ${AUTH_TOKEN}`) {
    return res.status(401).json({ error: 'Unauthorized' });
  }
  next();
}

// Endpoint REST para publicar mensagem em canal
app.use(express.json());
app.post('/publish/:canal', authMiddleware, (req, res) => {
  const canal = req.params.canal;
  const msg = req.body.msg;
  if (!canal || !msg) {
    return res.status(400).json({ error: 'Canal ou mensagem não especificados' });
  }
  io.to(canal).emit('mensagem', { canal, msg, timestamp: new Date().toISOString() });
  messagesPublished++;
  return res.json({ status: 'Mensagem publicada' });
});

// Endpoint REST para broadcast
app.post('/publish/broadcast', authMiddleware, (req, res) => {
  const msg = req.body.msg;
  if (!msg) {
    return res.status(400).json({ error: 'Mensagem não especificada' });
  }
  io.emit('mensagem', { canal: 'broadcast', msg, timestamp: new Date().toISOString() });
  messagesPublished++;
  return res.json({ status: 'Mensagem publicada para todos os canais' });
});

// Endpoint de estatísticas
app.get('/stats', (req, res) => {
  const activeChannels = Object.keys(channels).length;
  let connectedClients = 0;
  Object.values(channels).forEach(arr => connectedClients += arr.length);
  res.json({ activeChannels, connectedClients, messagesPublished });
});

// Endpoint de health check
app.get('/health', (req, res) => {
  res.json({ status: 'ok' });
});

// Socket.io conexão e gerenciamento de canais
io.on('connection', (socket) => {
  socket.on('join', (canal) => {
    socket.join(canal);
    if (!channels[canal]) channels[canal] = [];
    channels[canal].push(socket.id);
  });

  socket.on('leave', (canal) => {
    socket.leave(canal);
    if (channels[canal]) {
      channels[canal] = channels[canal].filter(id => id !== socket.id);
      if (channels[canal].length === 0) delete channels[canal];
    }
  });

  socket.on('disconnect', () => {
    for (const canal in channels) {
      channels[canal] = channels[canal].filter(id => id !== socket.id);
      if (channels[canal].length === 0) delete channels[canal];
    }
  });
});

const PORT = process.env.PORT || 8080;
server.listen(PORT, () => {
  console.log(`Servidor Socket.io rodando em :${PORT}`);
});
