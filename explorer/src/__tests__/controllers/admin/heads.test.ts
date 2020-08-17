import http from 'http'
import httpStatus from 'http-status-codes'
import { getRepository } from 'typeorm'
import { BigNumber } from 'bignumber.js'
import { clearDb } from '../../testdatabase'
import { createAdmin } from '../../../support/admin'
import { Head } from '../../../entity/Head'
import { start, stop } from '../../../support/server'
import { requestBuilder, RequestBuilder } from '../../../support/requestBuilder'

const USERNAME = 'myadmin'
const PASSWORD = 'validpassword'
const ADMIN_PATH = '/api/v1/admin'
const adminHeadsPath = `${ADMIN_PATH}/heads`

let server: http.Server
let rb: RequestBuilder

beforeAll(async () => {
  server = await start()
  rb = requestBuilder(server)
})
afterAll(done => stop(server, done))
beforeEach(async () => {
  await clearDb()
  await createAdmin(USERNAME, PASSWORD)
})

describe('GET /api/v1/admin/heads', () => {
  it('can retrieve heads', async done => {
    await createHead()

    rb.sendGet(adminHeadsPath, USERNAME, PASSWORD)
      .expect(httpStatus.OK)
      .expect(res => {
        expect(res.body.data.length).toEqual(1)
        expect(res.body.data[0].id).toBeDefined()
      })
      .end(done)
  })

  it('returns a 401 unauthorized with invalid admin credentials', async done => {
    rb.sendGet(adminHeadsPath, USERNAME, 'invalidpassword')
      .expect(httpStatus.UNAUTHORIZED)
      .end(done)
  })
})

describe('GET /api/v1/admin/heads/:id', () => {
  function path(id: number): string {
    return `${adminHeadsPath}/${id}`
  }

  it('can get a node', async done => {
    const head = await createHead()

    rb.sendGet(path(head.id), USERNAME, PASSWORD)
      .expect(httpStatus.OK)
      .expect(res => {
        expect(res.body.data.id).toBeDefined()
      })
      .end(done)
  })

  it('returns a 401 unauthorized with invalid admin credentials', async done => {
    const head = await createHead()

    rb.sendGet(path(head.id), USERNAME, 'invalidpassword')
      .expect(httpStatus.UNAUTHORIZED)
      .end(done)
  })
})

const DEFAULT_HEAD_ATTRS: Pick<
  Head,
  | 'blockHash'
  | 'parentHash'
  | 'uncleHash'
  | 'coinbase'
  | 'root'
  | 'txHash'
  | 'receiptHash'
  | 'bloom'
  | 'difficulty'
  | 'number'
  | 'gasLimit'
  | 'gasUsed'
  | 'time'
  | 'extra'
  | 'mixDigest'
  | 'nonce'
  | 'createdAt'
> = {
  blockHash: Buffer.from('abc123'),
  parentHash: Buffer.from('abc123'),
  uncleHash: Buffer.from('abc123'),
  coinbase: Buffer.from('abc123'),
  root: Buffer.from('abc123'),
  txHash: Buffer.from('abc123'),
  receiptHash: Buffer.from('abc123'),
  bloom: Buffer.from('abc123'),
  difficulty: new BigNumber('1'),
  number: new BigNumber('1'),
  gasLimit: new BigNumber('1'),
  gasUsed: new BigNumber('1'),
  time: new BigNumber('1'),
  extra: Buffer.from('abc123'),
  mixDigest: Buffer.from('abc123'),
  nonce: Buffer.from('abc123'),
  createdAt: new Date(),
}

function createHead(attrs: Partial<Head> = {}): Promise<Head> {
  const head = Head.build({ ...DEFAULT_HEAD_ATTRS, ...attrs })
  return getRepository(Head).save(head)
}
