import { bindActionCreators, Middleware } from 'redux'
import configureStore from 'redux-mock-store'
import thunk from 'redux-thunk'
import * as actions from 'actions'
import jsonApiJobSpecFactory from '../support/factories/jsonApiJobSpec'
import jsonApiJobSpecRunFactory from '../support/factories/jsonApiJobSpecRun'
import isoDate, { MINUTE_MS } from '../support/test-helpers/isoDate'
import globPath from '../support/test-helpers/globPath'

describe('fetchJob', () => {
  it('maintains dashed keys', done => {
    expect.assertions(1)

    const minuteAgo = isoDate(Date.now() - MINUTE_MS)
    const expectedTask = {
      id: 1,
      type: 'httpget',
      params: {
        headers: {
          'x-api-key': ['SOME_API_KEY'],
        },
      },
    }
    const jobSpecId = 'someid'
    const jobSpecResponse = jsonApiJobSpecFactory({
      createdAt: minuteAgo,
      id: jobSpecId,
      initiators: [{ type: 'web' }],
      tasks: [expectedTask],
    })

    global.fetch.getOnce(globPath(`/v2/specs/${jobSpecId}`), jobSpecResponse)

    const testMiddleware: Middleware = () => next => action => {
      next(action)

      if (action.type === 'UPSERT_JOB') {
        const task = action.data.specs[jobSpecId].attributes.tasks[0]
        expect(task).toEqual(expectedTask)
        done()
      }
    }

    const middlewares = [thunk, testMiddleware]
    const mockStore = configureStore(middlewares)
    const store = mockStore({})
    const fetchJob = bindActionCreators(actions.fetchJob, store.dispatch)

    fetchJob(jobSpecId)
  })
})

describe('fetchJobRun', () => {
  it('maintains dashed keys', done => {
    expect.assertions(1)

    const expectedTask = {
      type: 'noop',
      params: {
        headers: {
          'x-api-key': ['SOME_API_KEY'],
        },
      },
    }
    const taskRunA = {
      id: 'taskRunA',
      status: 'completed',
      task: expectedTask,
    }
    const runResponse = jsonApiJobSpecRunFactory({
      taskRuns: [taskRunA],
    })
    const id = runResponse.data.id
    global.fetch.getOnce(globPath(`/v2/runs/${id}`), runResponse)

    const testMiddleware: Middleware = () => next => action => {
      next(action)

      if (action.type === 'UPSERT_JOB_RUN') {
        const run = action.data.runs[id]
        const task = run.attributes.taskRuns[0].task
        expect(task).toEqual(expectedTask)
        done()
      }
    }
    const middlewares = [thunk, testMiddleware]
    const mockStore = configureStore(middlewares)
    const store = mockStore({})
    const fetchJobRun = bindActionCreators(actions.fetchJobRun, store.dispatch)

    fetchJobRun(id)
  })
})
