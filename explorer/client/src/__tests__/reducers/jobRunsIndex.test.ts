import reducer, {
  INITIAL_STATE as initialRootState,
  AppState,
} from '../../reducers'
import {
  FetchJobRunsSucceededAction,
  FetchJobRunSucceededAction,
} from '../../reducers/actions'

const INITIAL_STATE: AppState = {
  ...initialRootState,
  jobRunsIndex: { items: ['replace-me'] },
}

describe('reducers/jobRunsIndex', () => {
  it('returns the current state for other actions', () => {
    const action = {} as FetchJobRunsSucceededAction
    const state = reducer(INITIAL_STATE, action)

    expect(state.jobRunsIndex).toEqual(INITIAL_STATE.jobRunsIndex)
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
      }
      const state = reducer(INITIAL_STATE, action)

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
      const state = reducer(INITIAL_STATE, action)

      expect(state.jobRunsIndex).toEqual({
        items: undefined,
      })
    })
  })
})
