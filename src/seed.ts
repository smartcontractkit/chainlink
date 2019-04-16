import { getDb } from './database'
import { JobRun } from './entity/JobRun'
import { TaskRun } from './entity/TaskRun'
import { Initiator } from './entity/Initiator'

export const JOB_RUN_A_ID = 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa'
export const JOB_RUN_B_ID = 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb'

export default async () => {
  const dbConnection = getDb()
  const count = await dbConnection.manager.count(JobRun)

  if (count === 0) {
    const jobRunA = new JobRun()
    jobRunA.runId = JOB_RUN_A_ID
    jobRunA.jobId = 'cccccccccccccccccccccccccccccccc'
    jobRunA.status = 'in_progress'
    jobRunA.createdAt = new Date(Date.parse('2019-04-08T01:00:00.000Z'))
    await dbConnection.manager.save(jobRunA)

    const initiatorA = new Initiator()
    initiatorA.type = 'runlog'
    initiatorA.jobRun = jobRunA
    initiatorA.txHash = 'txA'
    initiatorA.requestId = 'requestIdA'
    initiatorA.requester = 'requesterA'
    await dbConnection.manager.save(initiatorA)

    const taskRunA = new TaskRun()
    taskRunA.jobRun = jobRunA
    taskRunA.index = 0
    taskRunA.status = 'in_progress'
    taskRunA.type = 'httpget'
    await dbConnection.manager.save(taskRunA)

    const jobRunB = new JobRun()
    jobRunB.runId = JOB_RUN_B_ID
    jobRunB.jobId = 'dddddddddddddddddddddddddddddddd'
    jobRunB.status = 'completed'
    jobRunB.createdAt = new Date(Date.parse('2019-04-09T01:00:00.000Z'))
    await dbConnection.manager.save(jobRunB)

    const initiatorB = new Initiator()
    initiatorB.type = 'runlog'
    initiatorB.jobRun = jobRunB
    initiatorB.txHash = 'txB'
    initiatorB.requestId = 'requestIdB'
    initiatorB.requester = 'requesterB'
    await dbConnection.manager.save(initiatorB)

    const taskRunB = new TaskRun()
    taskRunB.jobRun = jobRunB
    taskRunB.index = 0
    taskRunB.status = 'completed'
    taskRunB.type = 'ethbytes32'
    await dbConnection.manager.save(taskRunB)
  }
}
