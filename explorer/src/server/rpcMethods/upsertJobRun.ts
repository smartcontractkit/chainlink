import { fromJSONObject, saveJobRunTree } from '../../entity/JobRun'
import { logger } from '../../logging'
import { getDb } from '../../database'
import { serverContext } from './../handleMessage'
import jayson from 'jayson'

const { INVALID_PARAMS } = jayson.Server.errors

export default async (
  payload: object,
  context: serverContext,
  callback: jayson.JSONRPCCallbackTypePlain,
) => {
  try {
    const db = await getDb()
    const jobRun = fromJSONObject(payload)
    jobRun.chainlinkNodeId = context.chainlinkNodeId
    await saveJobRunTree(db, jobRun)
    callback(null, 'success')
  } catch (e) {
    logger.error(e)
    callback({ code: INVALID_PARAMS, message: 'invalid params' })
  }
}
