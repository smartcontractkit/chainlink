import Dockerode from 'dockerode'
import MockDockerode from '../__mocks__/dockerode'
import execa from 'execa'
import { assert } from 'chai'
import ChainlinkClient from './chainlinkClient'

let clClient: ChainlinkClient
const CLIENT_NAME = 'test_name'
const CLIENT_ENDPOINT = 'http://test.com'
const CONTAINER_NAME = 'test_container'
const ENV = { env: { CLIENT_NODE_URL: 'http://test.com', ROOT: '~/test_name' } }
const mock = Dockerode as typeof Dockerode & typeof MockDockerode

beforeEach(() => {
  clClient = new ChainlinkClient(CLIENT_NAME, CLIENT_ENDPOINT, CONTAINER_NAME)
})

afterEach(() => {
  jest.clearAllMocks()
})

describe('ChainlinkClient', () => {
  describe('#constructor', () => {
    it('sets the correct properties on the instance', () => {
      expect(mock.getContainer).toHaveBeenCalledTimes(1)
      expect(mock.getContainer).toHaveBeenCalledWith(CONTAINER_NAME)
      assert.equal(clClient.name, CLIENT_NAME)
      assert.equal(clClient.clientURL, CLIENT_ENDPOINT)
      assert.equal(clClient.clientURL, CLIENT_ENDPOINT)
      assert.equal(clClient.rootDir, `~/${CLIENT_NAME}`)
    })
  })

  describe('#pause', () => {
    it('calls the correct docker functions', async () => {
      mock.setContainerState({ Paused: false })
      await clClient.pause()
      expect(mock.container.pause).toHaveBeenCalledTimes(1)
    })
  })

  describe('#unpause', () => {
    it('calls the correct docker functions', async () => {
      mock.setContainerState({ Paused: true })
      await clClient.unpause()
      expect(mock.container.unpause).toHaveBeenCalledTimes(1)
    })
  })

  describe('#login', () => {
    it('runs the correct CLI command', async () => {
      clClient.login()
      expect(execa.sync).toHaveBeenCalledTimes(1)
      expect(execa.sync).toHaveBeenCalledWith(
        'chainlink',
        ['-j', 'admin', 'login', '--file', '/run/secrets/apicredentials'],
        ENV,
      )
    })
  })

  describe('#getJobs', () => {
    it('runs the correct CLI command', async () => {
      clClient.getJobs()
      expect(execa.sync).toHaveBeenCalledTimes(1)
      expect(execa.sync).toHaveBeenCalledWith(
        'chainlink',
        ['-j', 'job_specs', 'list'],
        ENV,
      )
    })
  })

  describe('#getJobRuns', () => {
    it('runs the correct CLI command', async () => {
      clClient.getJobRuns()
      expect(execa.sync).toHaveBeenCalledTimes(1)
      expect(execa.sync).toHaveBeenCalledWith(
        'chainlink',
        ['-j', 'runs', 'list'],
        ENV,
      )
    })
  })

  describe('#createJob', () => {
    it('runs the correct CLI command', async () => {
      const jobString = '{...}'
      clClient.createJob(jobString)
      expect(execa.sync).toHaveBeenCalledTimes(1)
      expect(execa.sync).toHaveBeenCalledWith(
        'chainlink',
        ['-j', 'job_specs', 'create', jobString],
        ENV,
      )
    })
  })

  describe('#archiveJob', () => {
    it('runs the correct CLI command', async () => {
      const jobID = '1234'
      clClient.archiveJob(jobID)
      expect(execa.sync).toHaveBeenCalledTimes(1)
      expect(execa.sync).toHaveBeenCalledWith(
        'chainlink',
        ['-j', 'job_specs', 'archive', jobID],
        ENV,
      )
    })
  })

  describe('#getAdminInfo', () => {
    it('runs the correct CLI command', async () => {
      clClient.getAdminInfo()
      expect(execa.sync).toHaveBeenCalledTimes(1)
      expect(execa.sync).toHaveBeenCalledWith(
        'chainlink',
        ['-j', 'keys', 'eth', 'list'],
        ENV,
      )
    })
  })
})
