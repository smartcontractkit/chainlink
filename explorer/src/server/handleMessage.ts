import { callRPCServer } from './rpcServer'
import { logger } from '../logging'
import { getDb } from '../database'
import { fromString, saveJobRunTree } from '../entity/JobRun'
import jayson from 'jayson'

export interface ServerContext {
  chainlinkNodeId: number
}

const handleLegacy = async (json: string, context: ServerContext) => {
  try {
    const db = await getDb()
    const jobRun = fromString(json)
    jobRun.chainlinkNodeId = context.chainlinkNodeId
    await saveJobRunTree(db, jobRun)
    return { status: 201 }
  } catch (e) {
    logger.error(e)
    return { status: 422 }
  }
}

const handleJSONRCP = async (request: string, context: ServerContext) => {
  return await new Promise(resolve => {
    callRPCServer(
      request,
      context,
      (error: jayson.JSONRPCErrorLike, response: jayson.JSONRPCResultLike) => {
        // resolve both error and success responses
        error ? resolve(error) : resolve(response)
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
