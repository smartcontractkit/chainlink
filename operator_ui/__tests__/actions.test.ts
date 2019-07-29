import * as actions from 'actions'
import jsonApiJobSpecFactory from 'factories/jsonApiJobSpec'
import configureStore from 'redux-mock-store'
import thunk from 'redux-thunk'
import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'

const middlewares = [thunk]
const mockStore = configureStore(middlewares)

describe('fetchJob', () => {
  it('maintains dashed keys', () => {
    const minuteAgo = isoDate(Date.now() - MINUTE_MS)
    const expectedTask = {
      id: 1,
      type: 'httpget',
      params: {
        headers: {
          'x-api-key': ['SOME_API_KEY']
        }
      }
    }
    const jobSpecId = 'someid'
    const jobSpecResponse = jsonApiJobSpecFactory({
      createdAt: minuteAgo,
      id: jobSpecId,
      initiators: [{ type: 'web' }],
      tasks: [expectedTask]
    })

    global.fetch.getOnce(`/v2/specs/${jobSpecId}`, jobSpecResponse)

    const store = mockStore({})
    return store.dispatch(actions.fetchJob(jobSpecId)).then(() => {
      const upsertJob = store.getActions()[1]
      const task = upsertJob.data.specs[jobSpecId].attributes.tasks[0]
      expect(task).toEqual(expectedTask)
    })
  })
})
