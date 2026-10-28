const express = require('express');
const cors = require('cors');
const path = require('path');
const fs = require('fs');
const statisticsRouter = require('./routes/statistics');

const PORT = Number(process.env.PORT) || 3000;

// URL de la API Go: se inyecta en el HTML en tiempo de request.
// - Local: http://localhost:8080 (default)
// - Render: se define via env var GO_API_URL en el panel de Render
const GO_API_URL = process.env.GO_API_URL || 'http://localhost:8080';

const app = express();
app.use(cors());
app.use(express.json({ limit: '5mb' }));

// Logger minimalista con duración
app.use((req, res, next) => {
  const start = Date.now();
  res.on('finish', () => {
    console.log(
      `[node-api] ${req.method} ${req.originalUrl} ${res.statusCode} ${Date.now() - start}ms`
    );
  });
  next();
});

// Ruta raíz: sirve el index.html con la URL de Go inyectada
app.get('/', (_req, res) => {
  try {
    const htmlPath = path.join(__dirname, '..', 'public', 'index.html');
    const html = fs.readFileSync(htmlPath, 'utf8').replace('__GO_API_URL__', GO_API_URL);
    res.type('html').send(html);
  } catch (err) {
    console.error('[node-api] error sirviendo index.html:', err);
    res.status(500).send('Error cargando la página');
  }
});

// Servir otros assets estáticos (CSS, JS, imágenes, etc.)
app.use(express.static(path.join(__dirname, '..', 'public')));

// Healthcheck
app.get('/health', (_req, res) => {
  res.json({ status: 'ok', service: 'node-api' });
});

// Rutas de negocio
app.use('/api/v1', statisticsRouter);

// 404 handler
app.use((_req, res) => {
  res.status(404).json({ error: 'Ruta no encontrada' });
});

// Error handler global
// eslint-disable-next-line no-unused-vars
app.use((err, _req, res, _next) => {
  console.error('[node-api] error no manejado:', err);
  res.status(err.status || 500).json({ error: err.message || 'Error interno' });
});

const server = app.listen(PORT, () => {
  console.log(`🚀 Node API escuchando en :${PORT}`);
  console.log(`   GO_API_URL = ${GO_API_URL}`);
});

// Graceful shutdown
['SIGINT', 'SIGTERM'].forEach((sig) => {
  process.on(sig, () => {
    console.log('🛑 Apagando Node API...');
    server.close(() => process.exit(0));
  });
});