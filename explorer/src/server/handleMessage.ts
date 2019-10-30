import rpcServer from './rpcServer'
import handleLegacy from './rpcMethods/legacy'

export type messageContext = {
  chainlinkNodeId: number
}

const handleJSONRCP = async (request: string, context: messageContext) => {
  return await new Promise((resolve, reject) => {
    // @ts-ignore - broken typing for server.call - should be able to accept 3 arguments
    // https://github.com/tedeh/jayson#server-context
    // https://github.com/tedeh/jayson/pull/152
    rpcServer.call(request, context, (error: any, response: any) => {
      // resolve both error and success responses
      error ? resolve(error) : resolve(response)
    })
  })
}

export const handleMessage = async (
  message: string,
  context: messageContext,
) => {
  if (message.includes('jsonrpc')) {
    return await handleJSONRCP(message, context)
  } else {
    return await handleLegacy(message, context)
  }
}
