import { fromJSONObject, saveJobRunTree } from './entity/JobRun'
import { logger } from './logging'
import { getDb } from './database'
import { messageContext } from './handleMessage'

import jayson from 'jayson'

const { INVALID_PARAMS } = jayson.Server.errors

const methods = {
  upsertJobRun: async (
    payload: any,
    context: messageContext,
    callback: (a: any, b?: any) => void,
  ) => {
    try {
      const db = await getDb()
      const jobRun = fromJSONObject(payload)
      jobRun.chainlinkNodeId = context.chainlinkNodeId
      await saveJobRunTree(db, jobRun)
      callback(null, 'success')
    } catch {
      callback({ code: INVALID_PARAMS, message: 'invalid params' })
    }
  },
}

const serverOptions = {
  useContext: true,
}

export default new jayson.Server(methods, serverOptions)
