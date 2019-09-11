import { ChainlinkNode } from './entity/ChainlinkNode'
import { createJobRun } from './factories'
import { getDb } from './database'

// NOT FOR PRODUCTION USE
// THE VALUES IN THIS FILE ARE ONLY USED FOR A LOCAL DEVELOPMENT ENVIRONMENT
// INVOKED FROM server.ts

// Unhashed Secret: 6l69Z1c2Qws5AlUuwG8Gy2vQJmX0FpVDBbPaUdxPp8nhMjrFtbHxFUl7MBFSHboo
const CORE_NODE_HASHED_SECRET =
  'c0fe784d9aa5e3bdc78adcde938fa7d9ede40f542ee9f4c4d626da8c8b406a41'
const CORE_NODE_SALT = '5dhpzBypVnddiYNRRz0eIO9Y6ZBqqv8D'
const CORE_NODE_ACCESS_KEY = 'bL1wMDLp4GVJ5p5n'

export default async () => {
  const db = await getDb()
  const count = await db.manager.count(ChainlinkNode)

  if (count === 0) {
    const node = await db.manager.save(buildChainlinkNode())

    await createJobRun(db, node)
    await createJobRun(db, node)
  }
}

const buildChainlinkNode = () => {
  const node = new ChainlinkNode()
  node.name = 'default'
  node.accessKey = CORE_NODE_ACCESS_KEY
  node.salt = CORE_NODE_SALT
  node.hashedSecret = CORE_NODE_HASHED_SECRET

  return node
}
