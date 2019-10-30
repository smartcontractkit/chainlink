/***** THIS IS NOT AN RPC METHOD ******/
/* THIS IS THE LEGACY SERVER FUNCTION */

import { fromString, saveJobRunTree } from '../../entity/JobRun'
import { logger } from '../../logging'
import { getDb } from '../../database'
import { messageContext } from './../handleMessage'

export default async (json: string, context: messageContext) => {
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
