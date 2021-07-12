import { partialAsFull } from 'support/test-helpers/partialAsFull'
import * as models from 'core/store/models'
import { bindActionCreators, Middleware } from 'redux'
import configureStore from 'redux-mock-store'
import thunk from 'redux-thunk'
import * as actionCreators from '../src/actionCreators'
import { ResourceActionType } from '../src/reducers/actions'
import jsonApiJobSpecFactory from '../support/factories/jsonApiJobSpec'
import jsonApiJobSpecRunFactory from '../support/factories/jsonApiJobSpecRun'
import globPath from '../support/test-helpers/globPath'
import isoDate, { MINUTE_MS } from '../support/test-helpers/isoDate'
import { RunStatus } from 'core/store/models'

describe('fetchJob', () => {
  it('maintains dashed keys', (done) => {
    expect.assertions(1)

    const minuteAgo = isoDate(Date.now() - MINUTE_MS)
    const expectedTask = partialAsFull<models.TaskSpec>({
      type: 'httpget',
      params: {
        headers: {
          'x-api-key': ['SOME_API_KEY'],
        },
      },
    })
    const jobSpecId = 'someid'
    const jobSpecResponse = jsonApiJobSpecFactory({
      createdAt: minuteAgo,
      id: jobSpecId,
      initiators: [{ type: 'web' }],
      tasks: [expectedTask],
    })

    global.fetch.getOnce(globPath(`/v2/specs/${jobSpecId}`), jobSpecResponse)

    const testMiddleware: Middleware = () => (next) => (action) => {
      next(action)

      if (action.type === ResourceActionType.UPSERT_JOB) {
        const task = action.data.specs[jobSpecId].attributes.tasks[0]
        expect(task).toEqual(expectedTask)
        done()
      }
    }

    const middlewares = [thunk, testMiddleware]
    const mockStore = configureStore(middlewares)
    const store = mockStore({})
    const fetchJob = bindActionCreators(actionCreators.fetchJob, store.dispatch)

    fetchJob(jobSpecId)
  })
})

describe('fetchJobRun', () => {
  it('maintains dashed keys', (done) => {
    expect.assertions(1)

    const expectedTask = partialAsFull<models.TaskSpec>({
      type: 'noop',
      params: {
        headers: {
          'x-api-key': ['SOME_API_KEY'],
        },
      },
    })
    const taskRunA = partialAsFull<models.TaskRun>({
      id: 'taskRunA',
      status: 'completed' as RunStatus.COMPLETED,
      task: expectedTask,
    })
    const runResponse = jsonApiJobSpecRunFactory({
      taskRuns: [taskRunA],
    })
    const id = runResponse.data.id
    global.fetch.getOnce(globPath(`/v2/runs/${id}`), runResponse)

    const testMiddleware: Middleware = () => (next) => (action) => {
      next(action)

      if (action.type === ResourceActionType.UPSERT_JOB_RUN) {
        const run = action.data.runs[id]
        const task = run.attributes.taskRuns[0].task
        expect(task).toEqual(expectedTask)
        done()
      }
    }
    const middlewares = [thunk, testMiddleware]
    const mockStore = configureStore(middlewares)
    const store = mockStore({})
    const fetchJobRun = bindActionCreators(
      actionCreators.fetchJobRun,
      store.dispatch,
    )

    fetchJobRun(id)
  })
})
