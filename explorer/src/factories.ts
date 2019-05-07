import { Connection } from 'typeorm'
import { v4 as uuid } from 'uuid'
import { getDb } from './database'
import { ChainlinkNode, createChainlinkNode } from './entity/ChainlinkNode'
import { JobRun } from './entity/JobRun'
import { TaskRun } from './entity/TaskRun'

export const createJobRun = async (
  db: Connection,
  chainlinkNode: ChainlinkNode
): Promise<JobRun> => {
  const jobRun = new JobRun()
  jobRun.chainlinkNodeId = chainlinkNode.id
  jobRun.runId = uuid()
  jobRun.jobId = uuid()
  jobRun.status = 'in_progress'
  jobRun.type = 'runlog'
  jobRun.txHash = 'tx' + uuid()
  jobRun.requestId = 'requestId' + uuid()
  jobRun.requester = 'requester' + uuid()
  jobRun.createdAt = new Date('2019-04-08T01:00:00.000Z')
  await db.manager.save(jobRun)

  const tr = new TaskRun()
  tr.jobRun = jobRun
  tr.index = 0
  tr.status = 'in_progress'
  tr.type = 'httpget'
  await db.manager.save(tr)

  return jobRun
}
