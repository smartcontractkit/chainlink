import reducer, {
  INITIAL_STATE as initialRootState,
  AppState,
} from '../../reducers'
import {
  FetchJobRunsSucceededAction,
  FetchJobRunSucceededAction,
} from '../../reducers/actions'
import { JobRun } from 'explorer/models'

const INITIAL_JOB_RUN = { id: 'replace-me' } as JobRun

const INITIAL_STATE: AppState = {
  ...initialRootState,
  jobRuns: {
    items: { 'replace-me': INITIAL_JOB_RUN },
  },
}

describe('reducers/jobRuns', () => {
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
      const state = reducer(INITIAL_STATE, action)

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
      const state = reducer(INITIAL_STATE, action)

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
