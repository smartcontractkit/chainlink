/* eslint-env jest */
import { fetchJob } from 'actions'

describe('actions', () => {
  describe('.fetchJob', () => {
    it('maintains snake case keys', async () => {
      const actionsCaptured = []
      const captureArgs = (...args) => actionsCaptured.push(...args)
      const jobSpecId = 'abc19'
      const jobSpecResponse = {
        data: {
          id: jobSpecId,
          type: 'jobSpec',
          attributes: {
            id: jobSpecId,
            snake_case: 'maintained',
            camelCase: 'maintained'
          }
        }
      }
      global.fetch.getOnce(`/v2/specs/${jobSpecId}`, jobSpecResponse)

      await fetchJob(jobSpecId)(captureArgs)
      const [json] = actionsCaptured.filter(j => j.type === 'UPSERT_JOB')
      expect(json).toHaveProperty(
        `data.jobSpec.${jobSpecId}.attributes.snake_case`,
        'maintained'
      )
      expect(json).toHaveProperty(
        `data.jobSpec.${jobSpecId}.attributes.camelCase`,
        'maintained'
      )
    })
  })
})
