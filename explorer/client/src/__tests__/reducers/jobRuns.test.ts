import reducer, {
  INITIAL_STATE as initialRootState,
  AppState,
} from '../../reducers'
import {
  JobRunsNormalizedData,
  JobRunNormalizedData,
  FetchJobRunsSucceededAction,
  FetchJobRunSucceededAction,
} from '../../reducers/actions'
import { JobRun } from 'explorer/models'
import { partialAsFull } from '../support/mocks'

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
      const data = partialAsFull<JobRunsNormalizedData>({
        jobRuns: {
          '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e': {
            id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e',
          },
        },
        meta: {
          currentPageJobRuns: {
            data: [{ id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e' }],
            meta: { count: 100 },
          },
        },
      })
      const action: FetchJobRunsSucceededAction = {
        type: 'FETCH_JOB_RUNS_SUCCEEDED',
        data: data,
      }
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
      const data = partialAsFull<JobRunNormalizedData>({
        jobRuns: {
          '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e': {
            id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e',
          },
        },
        meta: {
          jobRun: { meta: {} },
        },
      })
      const action: FetchJobRunSucceededAction = {
        type: 'FETCH_JOB_RUN_SUCCEEDED',
        data: data,
      }
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
