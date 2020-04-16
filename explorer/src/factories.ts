import { getRepository } from 'typeorm'
import { v4 as uuid } from 'uuid'
import { ChainlinkNode } from './entity/ChainlinkNode'
import { JobRun } from './entity/JobRun'
import { TaskRun } from './entity/TaskRun'

export const createJobRun = async (
  chainlinkNode: ChainlinkNode,
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
  const jobRunRepo = getRepository(JobRun)
  await jobRunRepo.save(jobRun)

  const tr = new TaskRun()
  tr.jobRun = jobRun
  tr.index = 0
  tr.status = 'in_progress'
  tr.type = 'httpget'
  tr.confirmations = '1'
  tr.minimumConfirmations = '3'
  const taskRunRepo = getRepository(TaskRun)
  await taskRunRepo.save(tr)

  return jobRun
}
