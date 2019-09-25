/* eslint-disable @typescript-eslint/no-var-requires */
;(async function() {
  const net = require('net')

  const [customResponse] = process.argv.slice(2)
  const defaultResponse = '{"last": "3843.95"}'
  const response = customResponse || defaultResponse
  const port = process.env.CYPRESS_JOB_SERVER_PORT || 6692

  const server = new net.Server(socket => {
    socket.on('data', () => {
      socket.write(`HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: ${response.length}

${response}`)
      socket.end()
    })
  })
  server.on('close', () => {
    server.unref()
  })

  const endpoint = await server.listen(port)
  console.log(`Job Server listening on port ${endpoint.address().port}`)
})()
