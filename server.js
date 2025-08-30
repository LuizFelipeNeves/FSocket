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

// Endpoint REST para publicar mensagem em channel
app.use(express.json());
app.post('/publish', authMiddleware, (req, res) => {
  const channel = req.body.channel;
  const msg = req.body.msg;
  const eventType = req.body.eventType || 'message';
  const extra = req.body.extra || {};
  if (!channel || !msg) {
    return res.status(400).json({ error: 'channel or message not specified' });
  }
  io.to(channel).emit(eventType, { msg, timestamp: new Date().toISOString(), ...extra  });
  messagesPublished++;
  return res.json({ status: 'Message published' });
});

// Endpoint REST para broadcast
app.post('/publish/broadcast', authMiddleware, (req, res) => {
  const msg = req.body.msg;
  const eventType = req.body.eventType || 'message';
  const extra = req.body.extra || {};
  if (!msg) {
    return res.status(400).json({ error: 'Message not specified' });
  }
  io.emit(eventType, { msg, timestamp: new Date().toISOString(), ...extra });
  messagesPublished++;
  return res.json({ status: 'Message published to all channels' });
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
  socket.on('join', (channel) => {
    socket.join(channel);
    if (!channels[channel]) channels[channel] = [];
    channels[channel].push(socket.id);
  });

  socket.on('leave', (channel) => {
    socket.leave(channel);
    if (channels[channel]) {
      channels[channel] = channels[channel].filter(id => id !== socket.id);
      if (channels[channel].length === 0) delete channels[channel];
    }
  });

  socket.on('disconnect', () => {
    for (const channel in channels) {
      channels[channel] = channels[channel].filter(id => id !== socket.id);
      if (channels[channel].length === 0) delete channels[channel];
    }
  });
});

const PORT = process.env.PORT || 8080;
server.listen(PORT, () => {
  console.log(`Servidor Socket.io rodando em :${PORT}`);
});
