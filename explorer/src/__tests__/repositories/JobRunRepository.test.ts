import { getCustomRepository } from 'typeorm'
import { JobRunRepository } from '../../repositories/JobRunRepository'
import * as JobRun from '../../entity/JobRun'
import ethtxFixture from '../fixtures/JobRun.ethtx.fixture.json'
import fixture from '../fixtures/JobRun.fixture.json'
import { createChainlinkNode } from '../../entity/ChainlinkNode'

describe('JobRunRepository tests', () => {
  describe('getFirst', () => {
    it('should return a job run with its task runs sorted', async () => {
      const [chainlinkNode] = await createChainlinkNode(
        'job-run-fromString-chainlink-node',
      )

      const jr1 = JobRun.fromString(JSON.stringify(ethtxFixture))
      jr1.chainlinkNodeId = chainlinkNode.id

      await JobRun.saveJobRunTree(jr1)

      const jrRepository = getCustomRepository(JobRunRepository)
      const fetchedJr = await jrRepository.getFirst()

      fetchedJr.taskRuns.forEach((tr, i) => expect(tr.index).toEqual(i))
    })
  })

  describe('findById', () => {
    it('should find a task run by its id with its task runs sorted', async () => {
      const [chainlinkNode] = await createChainlinkNode(
        'job-run-fromString-chainlink-node',
      )

      const jr1 = JobRun.fromString(JSON.stringify(ethtxFixture))
      jr1.chainlinkNodeId = chainlinkNode.id

      const jr2 = JobRun.fromString(JSON.stringify(fixture))
      jr2.chainlinkNodeId = chainlinkNode.id

      await JobRun.saveJobRunTree(jr1)
      await JobRun.saveJobRunTree(jr2)

      const jrRepository = getCustomRepository(JobRunRepository)

      const fetchedJr1 = await jrRepository.findById(jr1.id)
      expect(fetchedJr1.id).toEqual(jr1.id)
      fetchedJr1.taskRuns.forEach((tr, i) => expect(tr.index).toEqual(i))

      const fetchedJr2 = await jrRepository.findById(jr2.id)
      expect(fetchedJr2.id).toEqual(jr2.id)
      fetchedJr2.taskRuns.forEach((tr, i) => expect(tr.index).toEqual(i))
    })
  })
})
