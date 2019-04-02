import { getDb } from './database'
import { JobRun } from './entity/JobRun'
import { TaskRun } from './entity/TaskRun'

export const JOB_RUN_A_ID = '66eb9365-6c0c-487c-9297-7b1b44d87711'
export const JOB_RUN_B_ID = '81369a4d-76db-45a5-9192-869a023eced0'

export default async () => {
  const dbConnection = getDb()
  const count = await dbConnection.manager.count(JobRun)

  if (count === 0) {
    const jobRunA = new JobRun()
    jobRunA.id = JOB_RUN_A_ID
    jobRunA.jobId = 'd9b0dd13-091f-4f55-b718-d9e725ab96dd'
    jobRunA.status = 'in_progress'
    jobRunA.initiatorType = 'run_at'
    await dbConnection.manager.save(jobRunA)

    const taskRunA = new TaskRun()
    taskRunA.id = '5828158b-e075-466e-aabb-f1a52827d3a5'
    taskRunA.jobRun = jobRunA
    taskRunA.index = 0
    taskRunA.status = 'in_progress'
    taskRunA.type = 'httpget'
    await dbConnection.manager.save(taskRunA)

    const jobRunB = new JobRun()
    jobRunB.id = JOB_RUN_B_ID
    jobRunB.jobId = 'dbbb5305-5ec9-46e8-9bab-0891d2ad4578'
    jobRunB.status = 'completed'
    jobRunB.initiatorType = 'run_at'
    await dbConnection.manager.save(jobRunB)

    const taskRunB = new TaskRun()
    taskRunB.id = 'd62fbeb6-dc92-4191-87a3-b3c1ac4899ee'
    taskRunB.jobRun = jobRunB
    taskRunB.index = 0
    taskRunB.status = 'completed'
    taskRunB.type = 'ethbytes32'
    await dbConnection.manager.save(taskRunB)
  }
}
