import { getRepository } from 'typeorm'
import {
  ChainlinkNode,
  createChainlinkNode,
  deleteChainlinkNode,
  hashCredentials,
  uptime,
} from '../../entity/ChainlinkNode'
import { createSession, closeSession } from '../../entity/Session'

async function wait(sec: number) {
  return new Promise(res => {
    setTimeout(() => {
      res()
    }, sec * 1000)
  })
}

describe('createChainlinkNode', () => {
  it('returns a valid ChainlinkNode record', async () => {
    const [chainlinkNode, secret] = await createChainlinkNode(
      'new-valid-chainlink-node-record',
    )
    expect(chainlinkNode.accessKey).toHaveLength(16)
    expect(chainlinkNode.salt).toHaveLength(32)
    expect(chainlinkNode.hashedSecret).toBeDefined()
    expect(secret).toHaveLength(64)
  })

  it('reject duplicate ChainlinkNode names', async () => {
    await createChainlinkNode('identical')
    await expect(createChainlinkNode('identical')).rejects.toThrow()
  })
})

describe('deleteChainlinkNode', () => {
  it('deletes a ChainlinkNode with the specified name', async () => {
    await createChainlinkNode('chainlink-node-to-be-deleted')
    let count = await getRepository(ChainlinkNode).count()
    expect(count).toBe(1)
    await deleteChainlinkNode('chainlink-node-to-be-deleted')
    count = await getRepository(ChainlinkNode).count()
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
    const [node] = await createChainlinkNode('chainlink-node')
    const initialUptime = await uptime(node)
    expect(initialUptime).toEqual(0)
  })

  it('calculates uptime based on open and closed sessions', async () => {
    const [node] = await createChainlinkNode('chainlink-node')
    const session = await createSession(node)
    await wait(1)
    await closeSession(session)
    const uptime1 = await uptime(node)
    expect(uptime1).toBeGreaterThan(0)
    expect(uptime1).toBeLessThan(3)
    await createSession(node)
    await wait(1)
    const uptime2 = await uptime(node)
    expect(uptime2).toBeGreaterThan(uptime1)
    expect(uptime2).toBeLessThan(4)
  })
})
