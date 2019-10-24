import reducer, { State } from '../../reducers'
import {
  FetchJobRunsSucceededAction,
  FetchJobRunSucceededAction,
} from '../../reducers/actions'

const STATE = { jobRunsIndex: { items: ['replace-me'] } }

describe('reducers/jobRunsIndex', () => {
  it('returns the current state for other actions', () => {
    const action = {} as FetchJobRunsSucceededAction
    const state = reducer(STATE, action) as State

    expect(state.jobRunsIndex).toEqual(STATE.jobRunsIndex)
  })

  describe('FETCH_JOB_RUNS_SUCCEEDED', () => {
    it('can replace items', () => {
      const jobRuns = [{ id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e' }]
      const data = {
        meta: {
          currentPageJobRuns: {
            data: jobRuns,
            meta: { count: 100 },
          },
        },
        entities: {},
      }
      const action = {
        type: 'FETCH_JOB_RUNS_SUCCEEDED',
        data: data,
      } as FetchJobRunsSucceededAction
      const state = reducer(STATE, action) as State

      expect(state.jobRunsIndex).toEqual({
        items: ['9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e'],
        count: 100,
      })
    })
  })

  describe('FETCH_JOB_RUN_SUCCEEDED', () => {
    it('clears items', () => {
      const data = {
        jobRuns: {},
        meta: {
          jobRun: { meta: {} },
        },
      }
      const action = {
        type: 'FETCH_JOB_RUN_SUCCEEDED',
        data: data,
      } as FetchJobRunSucceededAction
      const state = reducer(STATE, action) as State

      expect(state.jobRunsIndex).toEqual({
        items: undefined,
      })
    })
  })
})
