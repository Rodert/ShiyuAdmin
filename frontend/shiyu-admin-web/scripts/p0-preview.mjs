import express from 'express';
import http from 'node:http';
import { resolve } from 'node:path';

const app = express();
app.use('/api', (req, res) => {
  const upstream = http.request(`http://127.0.0.1:18808${req.originalUrl}`, { method: req.method, headers: req.headers }, (reply) => {
    res.writeHead(reply.statusCode, reply.headers);
    reply.pipe(res);
  });
  upstream.on('error', () => { if (!res.headersSent) res.status(502).end(); });
  req.pipe(upstream);
});
app.use(express.static(resolve('dist')));
app.get('*', (_, res) => res.sendFile(resolve('dist/index.html')));
app.listen(18800, '127.0.0.1');
