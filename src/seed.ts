import { getDb } from './database'
import { JobRun } from './entity/JobRun'
import { TaskRun } from './entity/TaskRun'
import { Client, createClient } from './entity/Client'

export const JOB_RUN_A_ID = 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa'
export const JOB_RUN_B_ID = 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb'

export default async () => {
  const dbConnection = await getDb()
  const count = await dbConnection.manager.count(Client)

  if (count === 0) {
    const [client, _] = await createClient(dbConnection, 'default')

    const jobRunA = new JobRun()
    jobRunA.clientId = client.id
    jobRunA.runId = JOB_RUN_A_ID
    jobRunA.jobId = 'cccccccccccccccccccccccccccccccc'
    jobRunA.status = 'in_progress'
    jobRunA.type = 'runlog'
    jobRunA.txHash = 'txA'
    jobRunA.requestId = 'requestIdA'
    jobRunA.requester = 'requesterA'
    jobRunA.createdAt = new Date(Date.parse('2019-04-08T01:00:00.000Z'))
    await dbConnection.manager.save(jobRunA)

    const taskRunA = new TaskRun()
    taskRunA.jobRun = jobRunA
    taskRunA.index = 0
    taskRunA.status = 'in_progress'
    taskRunA.type = 'httpget'
    await dbConnection.manager.save(taskRunA)

    const jobRunB = new JobRun()
    jobRunB.clientId = client.id
    jobRunB.runId = JOB_RUN_B_ID
    jobRunB.jobId = 'dddddddddddddddddddddddddddddddd'
    jobRunB.status = 'completed'
    jobRunB.type = 'web'
    jobRunB.createdAt = new Date(Date.parse('2019-04-09T01:00:00.000Z'))
    await dbConnection.manager.save(jobRunB)

    const taskRunB = new TaskRun()
    taskRunB.jobRun = jobRunB
    taskRunB.index = 0
    taskRunB.status = 'completed'
    taskRunB.type = 'ethbytes32'
    await dbConnection.manager.save(taskRunB)
  }
}
