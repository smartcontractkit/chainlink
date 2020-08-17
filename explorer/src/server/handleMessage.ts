import rpcServer from './rpcServer'
import { logger } from '../logging'
import { fromString, saveJobRunTree } from '../entity/JobRun'
import jayson from 'jayson'

export interface ServerContext {
  chainlinkNodeId: number
}

// legacy server response synonymous with upsertJobRun RPC method
const handleLegacy = async (json: string, context: ServerContext) => {
  try {
    const jobRun = fromString(json)
    jobRun.chainlinkNodeId = context.chainlinkNodeId
    await saveJobRunTree(jobRun)
    return { status: 201 }
  } catch (e) {
    logger.error(e)
    return { status: 422 }
  }
}

const handleJSONRCP = (request: string, context: ServerContext) => {
  return new Promise(resolve => {
    rpcServer.call(
      request,
      context,
      (error: jayson.JSONRPCErrorLike, result: jayson.JSONRPCResultLike) => {
        // resolve both errored and successful responses
        if (error) {
          logger.error(error.message)
          resolve(error)
        } else {
          resolve(result)
        }
      },
    )
  })
}

export const handleMessage = async (
  message: string,
  context: ServerContext,
) => {
  if (message.includes('jsonrpc')) {
    return await handleJSONRCP(message, context)
  } else {
    return await handleLegacy(message, context)
  }
}
