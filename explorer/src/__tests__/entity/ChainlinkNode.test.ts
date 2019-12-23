import { Connection } from 'typeorm'
import {
  ChainlinkNode,
  createChainlinkNode,
  deleteChainlinkNode,
  hashCredentials,
  uptime,
} from '../../entity/ChainlinkNode'
import { Session, createSession, closeSession } from '../../entity/Session'
import { closeDbConnection, getDb } from '../../database'

async function wait(sec: number) {
  return new Promise(res => {
    setTimeout(() => {
      res()
    }, sec * 1000)
  })
}

let db: Connection

beforeAll(async () => {
  db = await getDb()
})

afterAll(async () => closeDbConnection())

describe('createChainlinkNode', () => {
  it('returns a valid ChainlinkNode record', async () => {
    const [chainlinkNode, secret] = await createChainlinkNode(
      db,
      'new-valid-chainlink-node-record',
    )
    expect(chainlinkNode.accessKey).toHaveLength(16)
    expect(chainlinkNode.salt).toHaveLength(32)
    expect(chainlinkNode.hashedSecret).toBeDefined()
    expect(secret).toHaveLength(64)
  })

  it('reject duplicate ChainlinkNode names', async () => {
    await createChainlinkNode(db, 'identical')
    await expect(createChainlinkNode(db, 'identical')).rejects.toThrow()
  })
})

describe('deleteChainlinkNode', () => {
  it('deletes a ChainlinkNode with the specified name', async () => {
    await createChainlinkNode(db, 'chainlink-node-to-be-deleted')
    let count = await db.manager.count(ChainlinkNode)
    expect(count).toBe(1)
    await deleteChainlinkNode(db, 'chainlink-node-to-be-deleted')
    count = await db.manager.count(ChainlinkNode)
    expect(count).toBe(0)
  })
})

describe('hashCredentials', () => {
  it('returns a sha256 signature', () => {
    expect(hashCredentials('a', 'b', 'c')).toHaveLength(64)
  })
})

describe('uptime', () => {
  it('returns 0 when no sessions exist', async () => {
    const [node, _] = await createChainlinkNode(db, 'chainlink-node')
    const initialUptime = await uptime(db, node)
    expect(initialUptime).toEqual(0)
  })

  it('calculates uptime based on open and closed sessions', async () => {
    const [node, _] = await createChainlinkNode(db, 'chainlink-node')
    const session = await createSession(db, node)
    await wait(1)
    await closeSession(db, session)
    const uptime1 = await uptime(db, node)
    expect(uptime1).toBeGreaterThan(0)
    expect(uptime1).toBeLessThan(3)
    await createSession(db, node)
    await wait(1)
    const uptime2 = await uptime(db, node)
    expect(uptime2).toBeGreaterThan(uptime1)
    expect(uptime2).toBeLessThan(4)
  })
})
