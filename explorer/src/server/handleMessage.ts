import { saveJobRunTree, fromString } from '../entity/JobRun'
import { logger } from '../logging'
import rpcServer from './rpcServer'
import { getDb } from '../database'

// todo make jsonRPC
export const legacyErrorResponse = { status: 422 }

export type messageContext = {
  chainlinkNodeId: number
}

const handleLegacy = async (json: string, context: messageContext) => {
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
