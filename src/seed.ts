import { JobRun } from "./entity/JobRun"

const seed = async (dbConnection: any) => {
  const count = await dbConnection.manager.count(JobRun)

  if (count === 0) {
    const jobRunA = new JobRun()
    jobRunA.requestId = "66eb9365-6c0c-487c-9297-7b1b44d87711"
    await dbConnection.manager.save(jobRunA)
    console.log("Saved a new job run with id: " + jobRunA.id)

    const jobRunB = new JobRun()
    jobRunB.requestId = "66eb9365-6c0c-487c-9297-7b1b44d87711"
    await dbConnection.manager.save(jobRunB)
    console.log("Saved a new job run with id: " + jobRunB.id)
  }
}

export default seed
