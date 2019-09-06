import { Connection } from 'typeorm'
import { closeDbConnection, getDb } from '../../database'
import { createChainlinkNode } from '../../entity/ChainlinkNode'
import { fromString, JobRun, saveJobRunTree } from '../../entity/JobRun'
import fixture from '../fixtures/JobRun.fixture.json'

let db: Connection

beforeAll(async () => {
  db = await getDb()
})

afterAll(async () => closeDbConnection())

describe('entity/taskRun', () => {
  it('copies old confirmations to new column on INSERT', async () => {
    const [chainlinkNode, _] = await createChainlinkNode(
      db,
      'testOverwriteJobRunsErrorOnConflict',
    )

    const jr = fromString(JSON.stringify(fixture))
    jr.chainlinkNodeId = chainlinkNode.id
    await saveJobRunTree(db, jr)
    expect(jr.id).toBeDefined()

    // insert into old columns
    await db.manager.query(
      `
      INSERT INTO task_run("jobRunId", index, type, status, confirmations, "minimumConfirmations")
      VALUES ($1, 1, 'randomtask', 'in_progress', 1, 2);
    `,
      [jr.id],
    )

    const retrieved = await db.manager.findOne(JobRun, jr.id)
    const task = retrieved.taskRuns[1]

    expect(task.confirmationsOld).toEqual(1)
    expect(task.minimumConfirmationsOld).toEqual(2)
    expect(task.confirmations).toEqual('1')
    expect(task.minimumConfirmations).toEqual('2')
  })

  it('copies old confirmations to new column on UPDATE', async () => {
    const [chainlinkNode, _] = await createChainlinkNode(
      db,
      'testOverwriteJobRunsErrorOnConflict',
    )

    const jr = fromString(JSON.stringify(fixture))
    jr.chainlinkNodeId = chainlinkNode.id
    await saveJobRunTree(db, jr)
    expect(jr.id).toBeDefined()
    const tr = jr.taskRuns[0]

    // update old columns
    await db.manager.query(
      `
      UPDATE task_run SET confirmations = 9, "minimumConfirmations" = 10
      WHERE id = $1;
    `,
      [tr.id],
    )

    const retrieved = await db.manager.findOne(JobRun, jr.id)
    const task = retrieved.taskRuns[0]

    expect(task.confirmationsOld).toEqual(9)
    expect(task.minimumConfirmationsOld).toEqual(10)
    expect(task.confirmations).toEqual('9')
    expect(task.minimumConfirmations).toEqual('10')
  })
})
