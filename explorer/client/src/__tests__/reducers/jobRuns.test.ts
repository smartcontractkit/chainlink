import reducer, { State } from '../../reducers'
import {
  FetchJobRunsSucceededAction,
  FetchJobRunSucceededAction,
} from '../../reducers/actions'

const STATE = {
  jobRuns: {
    items: { 'replace-me': { id: 'replace-me' } },
  },
}

describe('reducers/jobRuns', () => {
  it('returns the current state for other actions', () => {
    const action = {} as FetchJobRunsSucceededAction
    const state = reducer(STATE, action) as State

    expect(state.jobRuns).toEqual(STATE.jobRuns)
  })

  describe('FETCH_JOB_RUNS_SUCCEEDED', () => {
    it('can replace items', () => {
      const normalizedJobRuns = {
        '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e': {
          id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e',
        },
      }
      const orderedJobRuns = [{ id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e' }]
      const data = {
        jobRuns: normalizedJobRuns,
        meta: {
          currentPageJobRuns: {
            data: orderedJobRuns,
            meta: {},
          },
        },
      }
      const action = {
        type: 'FETCH_JOB_RUNS_SUCCEEDED',
        data: data,
      } as FetchJobRunsSucceededAction
      const state = reducer(STATE, action) as State

      expect(state.jobRuns).toEqual({
        items: {
          '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e': {
            id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e',
          },
        },
      })
    })
  })

  describe('FETCH_JOB_RUN_SUCCEEDED', () => {
    it('can replace items', () => {
      const normalizedJobRuns = {
        '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e': {
          id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e',
        },
      }
      const data = {
        jobRuns: normalizedJobRuns,
        meta: {
          jobRun: { meta: {} },
        },
      }
      const action = {
        type: 'FETCH_JOB_RUN_SUCCEEDED',
        data: data,
      } as FetchJobRunSucceededAction
      const state = reducer(STATE, action) as State

      expect(state.jobRuns).toEqual({
        items: {
          '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e': {
            id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e',
          },
        },
      })
    })
  })
})
