import { fromJSONObject, saveJobRunTree } from '../../entity/JobRun'
import { logger } from '../../logging'
import { ServerContext } from './../handleMessage'
import jayson from 'jayson'

// default invalid params error
const { INVALID_PARAMS } = jayson.Server.errors
const invalidParamsError = new jayson.Server().error(INVALID_PARAMS)

export default async (
  payload: object,
  context: ServerContext,
  callback: jayson.JSONRPCCallbackTypePlain,
) => {
  try {
    const jobRun = fromJSONObject(payload)
    jobRun.chainlinkNodeId = context.chainlinkNodeId
    await saveJobRunTree(jobRun)
    callback(null, 'success')
  } catch (e) {
    logger.error(e)
    callback(invalidParamsError)
  }
}
