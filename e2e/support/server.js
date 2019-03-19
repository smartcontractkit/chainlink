const net = require('net')

module.exports = {
  newServer: async response => {
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

    const endpoint = await server.listen()
    return {
      port: endpoint.address().port,
      close: () => {
        return Promise.all([endpoint.close(), server.close()])
      }
    }
  }
}
