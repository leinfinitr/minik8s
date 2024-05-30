const express = require('express');
const client = require('prom-client');

const app = express();
const register = new client.Registry();

// 创建一个自定义的指标
const counter = new client.Counter({
  name: 'nodejs_request_count',
  help: 'Total number of requests',
  registers: [register],
});

const gauge = new client.Gauge({
  name: 'nodejs_random_value',
  help: 'A random value for demonstration purposes',
  registers: [register],
});

// 创建一个中间件来增加请求计数
app.use((req, res, next) => {
  counter.inc();
  next();
});

// 定义一个路由来返回随机值并更新指标
app.get('/random', (req, res) => {
  const randomValue = Math.random();
  gauge.set(randomValue);
  res.json({ value: randomValue });
});

// 暴露 /metrics 端点给 Prometheus
app.get('/metrics', async (req, res) => {
  res.set('Content-Type', register.contentType);
  res.end(await register.metrics());
});

// 启动服务器
const port = process.env.PORT || 3000;
app.listen(port, () => {
  console.log(`Server is running on http://localhost:${port}`);
});