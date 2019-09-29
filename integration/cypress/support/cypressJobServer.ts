/* eslint-disable @typescript-eslint/no-var-requires */
import net from 'net'
;(async function() {
  const net = require('net')

  const [customResponse] = process.argv.slice(2)
  const defaultResponse = '{"last": "3843.95"}'
  const response = customResponse || defaultResponse
  const port = process.env.CYPRESS_JOB_SERVER_PORT || 6692

  const server = new net.Server((socket: net.Socket) => {
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
  const address = endpoint.address()
  if (address && typeof address != 'string') {
    console.log(`Job Server listening on port ${address.port}`)
  } else {
    console.error(
      'Invalid server setup. Address should be of type net.AddressInfo',
    )
    process.exit(1)
  }
})()
