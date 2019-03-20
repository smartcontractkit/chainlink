import { JobRun } from "./entity/JobRun"

const seed = async (dbConnection: any) => {
  const count = await dbConnection.manager.count(JobRun)

  if (count === 0) {
    const jobRunA = new JobRun()
    jobRunA.requestId = "66eb9365-6c0c-487c-9297-7b1b44d87711"
    jobRunA.jobId = "d9b0dd13-091f-4f55-b718-d9e725ab96dd" 
    await dbConnection.manager.save(jobRunA)
    console.log("Saved a new job run with id: " + jobRunA.id)

    const jobRunB = new JobRun()
    jobRunB.requestId =  "81369a4d-76db-45a5-9192-869a023eced0"
    jobRunB.jobId =   "dbbb5305-5ec9-46e8-9bab-0891d2ad4578"
    await dbConnection.manager.save(jobRunB)
    console.log("Saved a new job run with id: " + jobRunB.id)
  }
}

export default seed
