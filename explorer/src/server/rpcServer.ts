import jayson from 'jayson'
import rpcMethods from './rpcMethods'

const serverOptions = {
  useContext: true, // permits passing extra data object to RPC methods as 'server context'
}

export default new jayson.Server(rpcMethods, serverOptions)
