import { ChainlinkNode } from './entity/ChainlinkNode'
import { Connection } from 'typeorm'
import { createChainlinkNode } from '../src/entity/ChainlinkNode'
import { createJobRun } from './factories'
import { getDb } from './database'

export const JOB_RUN_B_ID = 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb'

export default async () => {
  const db = await getDb()
  const count = await db.manager.count(ChainlinkNode)

  if (count === 0) {
    const [node, _] = await createChainlinkNode(db, 'default')

    await createJobRun(db, node)
    await createJobRun(db, node)
  }
}
