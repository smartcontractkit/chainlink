/* eslint-disable @typescript-eslint/no-var-requires */
;(async function() {
  const net = require('net')

  const [response] = process.argv.slice(2)
  const port = process.env.CYPRESS_JOB_SERVER_PORT

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
