import { getConnection, getCustomRepository, getRepository } from 'typeorm'
import { createChainlinkNode } from '../../entity/ChainlinkNode'
import { fromString, JobRun, saveJobRunTree } from '../../entity/JobRun'
import fixture from '../fixtures/JobRun.fixture.json'
import { JobRunRepository } from '../../repositories/JobRunRepository'

describe('entity/taskRun', () => {
  it('copies old confirmations to new column on INSERT', async () => {
    const [chainlinkNode] = await createChainlinkNode(
      'testOverwriteJobRunsErrorOnConflict',
    )

    const jr = fromString(JSON.stringify(fixture))
    jr.chainlinkNodeId = chainlinkNode.id
    await saveJobRunTree(jr)
    expect(jr.id).toBeDefined()

    // insert into old columns
    await getConnection().query(
      `
      INSERT INTO task_run("jobRunId", index, type, status, confirmations, "minimumConfirmations")
      VALUES ($1, 1, 'randomtask', 'in_progress', 1, 2);
    `,
      [jr.id],
    )

    const jobRunRepository = getCustomRepository(JobRunRepository)
    const retrieved = await jobRunRepository.getFirst()

    const task = retrieved.taskRuns[1]

    expect(task.confirmationsOld).toEqual(1)
    expect(task.minimumConfirmationsOld).toEqual(2)
    expect(task.confirmations).toEqual('1')
    expect(task.minimumConfirmations).toEqual('2')
  })

  it('copies old confirmations to new column on UPDATE', async () => {
    const [chainlinkNode] = await createChainlinkNode(
      'testOverwriteJobRunsErrorOnConflict',
    )

    const jr = fromString(JSON.stringify(fixture))
    jr.chainlinkNodeId = chainlinkNode.id
    await saveJobRunTree(jr)
    expect(jr.id).toBeDefined()
    const tr = jr.taskRuns[0]

    // update old columns
    await getConnection().query(
      `
      UPDATE task_run SET confirmations = 9, "minimumConfirmations" = 10
      WHERE id = $1;
    `,
      [tr.id],
    )

    const retrieved = await getRepository(JobRun).findOne(jr.id)
    const task = retrieved.taskRuns[0]

    expect(task.confirmationsOld).toEqual(9)
    expect(task.minimumConfirmationsOld).toEqual(10)
    expect(task.confirmations).toEqual('9')
    expect(task.minimumConfirmations).toEqual('10')
  })
})
