const http = require('http');
const fs = require('fs').promises;
const path = require('path');

const PORT = process.env.PORT || 3000;

const server = http.createServer(async (req, res) => {
  res.setHeader('Content-Type', 'application/json');
  
  try {
    if (req.method === 'GET') {
      const data = await fs.readFile(path.join(__dirname, 'get-response.json'), 'utf8');
      res.writeHead(200);
      res.end(data);
    } else if (req.method === 'POST') {
      const data = await fs.readFile(path.join(__dirname, 'post-response.json'), 'utf8');
      res.writeHead(200);
      res.end(data);
    } else {
      res.writeHead(405);
      res.end(JSON.stringify({ error: 'Method not allowed' }));
    }
  } catch (error) {
    console.error('Error reading file:', error);
    res.writeHead(500);
    res.end(JSON.stringify({ error: 'Internal server error' }));
  }
});

server.listen(PORT, () => {
  console.log(`Server running on http://localhost:${PORT}`);
  console.log('GET requests return: get-response.json');
  console.log('POST requests return: post-response.json');
});