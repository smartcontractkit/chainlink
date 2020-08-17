import http from 'http'
import httpStatus from 'http-status-codes'
import { createAdmin } from '../../../support/admin'
import {
  createChainlinkNode,
  find as findNode,
} from '../../../entity/ChainlinkNode'
import { start, stop } from '../../../support/server'
import { requestBuilder, RequestBuilder } from '../../../support/requestBuilder'

const USERNAME = 'myadmin'
const PASSWORD = 'validpassword'
const ADMIN_PATH = '/api/v1/admin'
const adminNodesPath = `${ADMIN_PATH}/nodes`

let server: http.Server
let rb: RequestBuilder

beforeAll(async () => {
  server = await start()
  rb = requestBuilder(server)
})
afterAll(done => stop(server, done))
beforeEach(async () => {
  await createAdmin(USERNAME, PASSWORD)
})

describe('POST /api/v1/admin/nodes', () => {
  it('can create a node and returns the generated information', done => {
    const data = { name: 'nodeA', url: 'http://nodea.com' }

    rb.sendPost(adminNodesPath, USERNAME, PASSWORD, data)
      .expect(httpStatus.CREATED)
      .expect(res => {
        expect(res.body.id).toBeDefined()
        expect(res.body.accessKey).toBeDefined()
        expect(res.body.secret).toBeDefined()
      })
      .end(done)
  })

  it('returns an error with invalid params', done => {
    const data = { url: 'http://nodea.com' }

    rb.sendPost(adminNodesPath, USERNAME, PASSWORD, data)
      .expect(httpStatus.UNPROCESSABLE_ENTITY)
      .expect(res => {
        const errors = res.body.errors

        expect(errors).toBeDefined()
        expect(errors.name).toEqual({
          minLength: 'must be at least 3 characters',
        })
      })
      .end(done)
  })

  it('returns an error when the node already exists', async done => {
    const [node] = await createChainlinkNode('nodeA')
    const data = { name: node.name }

    rb.sendPost(adminNodesPath, USERNAME, PASSWORD, data)
      .expect(httpStatus.CONFLICT)
      .end(done)
  })

  it('returns a 401 unauthorized with invalid admin credentials', done => {
    rb.sendPost(adminNodesPath, USERNAME, 'invalidpassword')
      .expect(httpStatus.UNAUTHORIZED)
      .end(done)
  })
})

describe('DELETE /api/v1/admin/nodes/:name', () => {
  function path(name: string): string {
    return `${adminNodesPath}/${name}`
  }

  it('can delete a node', async done => {
    const [node] = await createChainlinkNode('nodeA')

    rb.sendDelete(path(node.name), USERNAME, PASSWORD)
      .expect(httpStatus.OK)
      .expect(async () => {
        const nodeAfter = await findNode(node.id)
        expect(nodeAfter).not.toBeDefined()
      })
      .end(done)
  })

  it('returns a 401 unauthorized with invalid admin credentials', done => {
    rb.sendDelete(path('idontexist'), USERNAME, 'invalidpassword')
      .expect(httpStatus.UNAUTHORIZED)
      .end(done)
  })
})

describe('GET /api/v1/admin/nodes/:id', () => {
  function path(id: number): string {
    return `${adminNodesPath}/${id}`
  }

  it('can get a node', async done => {
    const [node] = await createChainlinkNode('nodeA')

    rb.sendGet(path(node.id), USERNAME, PASSWORD)
      .expect(httpStatus.OK)
      .expect(res => {
        expect(res.body.data.id).toBeDefined()
      })
      .end(done)
  })

  it('returns a 401 unauthorized with invalid admin credentials', async done => {
    const [node] = await createChainlinkNode('nodeA')
    const _nodePath = path(node.id)
    rb.sendGet(_nodePath, USERNAME, 'invalidpassword')
      .expect(httpStatus.UNAUTHORIZED)
      .end(done)
  })
})
