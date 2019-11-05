import jayson from 'jayson'
import * as rpcMethods from './rpcMethods'

const serverOptions = {
  useContext: true, // permits passing extra data object to RPC methods as 'server context'
}

export const rpcServer = new jayson.Server(rpcMethods, serverOptions)

// broken typing for server.call - should be able to accept 3 arguments
// https://github.com/tedeh/jayson#server-context
// https://github.com/tedeh/jayson/pull/152

type callFunction = (
  request: jayson.JSONRPCRequestLike | Array<jayson.JSONRPCRequestLike>,
  context: object,
  originalCallback?: jayson.JSONRPCCallbackType,
) => void

export const callRPCServer: callFunction = rpcServer.call.bind(rpcServer)
