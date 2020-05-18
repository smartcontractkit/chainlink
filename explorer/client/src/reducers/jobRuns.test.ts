import { partialAsFull } from '@chainlink/ts-helpers'
import { JobRun } from 'explorer/models'
import reducer, {
  INITIAL_STATE as initialRootState,
  AppState,
} from '../reducers'
import {
  JobRunsNormalizedData,
  JobRunNormalizedData,
  FetchJobRunsBeginAction,
  FetchJobRunsSucceededAction,
  FetchJobRunsErrorAction,
  FetchJobRunBeginAction,
  FetchJobRunSucceededAction,
  FetchJobRunErrorAction,
} from '../reducers/actions'

const INITIAL_JOB_RUN = partialAsFull<JobRun>({ id: 'replace-me' })

const INITIAL_STATE: AppState = {
  ...initialRootState,
  jobRuns: {
    loading: false,
    error: false,
    items: { 'replace-me': INITIAL_JOB_RUN },
  },
}

describe('reducers/jobRuns', () => {
  const fetchJobRunsBeginAction: FetchJobRunsBeginAction = {
    type: 'FETCH_JOB_RUNS_BEGIN',
  }
  const fetchJobRunBeginAction: FetchJobRunBeginAction = {
    type: 'FETCH_JOB_RUN_BEGIN',
  }

  describe('FETCH_JOB_RUNS_BEGIN', () => {
    it('sets loading to true', () => {
      const state = reducer(INITIAL_STATE, fetchJobRunsBeginAction)

      expect(state.jobRuns.loading).toEqual(true)
    })
  })

  describe('FETCH_JOB_RUNS_SUCCEEDED', () => {
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
      data,
    }

    it('replaces items', () => {
      const state = reducer(INITIAL_STATE, action)

      expect(state.jobRuns.items).toEqual({
        '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e': {
          id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e',
        },
      })
    })

    it('sets loading to false', () => {
      let state = reducer(INITIAL_STATE, fetchJobRunsBeginAction)
      expect(state.jobRuns.loading).toEqual(true)

      state = reducer(state, action)
      expect(state.jobRuns.loading).toEqual(false)
    })
  })

  describe('FETCH_JOB_RUNS_ERROR', () => {
    const action: FetchJobRunsErrorAction = {
      type: 'FETCH_JOB_RUNS_ERROR',
      errors: [
        {
          status: 500,
          detail: 'An error',
        },
      ],
    }

    it('sets loading to false', () => {
      let state = reducer(INITIAL_STATE, fetchJobRunsBeginAction)
      expect(state.jobRuns.loading).toEqual(true)

      state = reducer(state, action)
      expect(state.jobRuns.loading).toEqual(false)
    })

    it('sets error to true & resets it to false when a new fetch starts', () => {
      let state = reducer(INITIAL_STATE, action)
      expect(state.jobRuns.error).toEqual(true)

      state = reducer(state, fetchJobRunsBeginAction)
      expect(state.jobRuns.error).toEqual(false)
    })
  })

  describe('FETCH_JOB_RUN_BEGIN', () => {
    it('sets loading to true', () => {
      const state = reducer(INITIAL_STATE, fetchJobRunBeginAction)

      expect(state.jobRuns.loading).toEqual(true)
    })
  })

  describe('FETCH_JOB_RUN_SUCCEEDED', () => {
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
      data,
    }

    it('replaces items', () => {
      const state = reducer(INITIAL_STATE, action)

      expect(state.jobRuns.items).toEqual({
        '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e': {
          id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e',
        },
      })
    })

    it('sets loading to false', () => {
      let state = reducer(INITIAL_STATE, fetchJobRunBeginAction)
      expect(state.jobRuns.loading).toEqual(true)

      state = reducer(state, action)
      expect(state.jobRuns.loading).toEqual(false)
    })
  })

  describe('FETCH_JOB_RUN_ERROR', () => {
    const action: FetchJobRunErrorAction = {
      type: 'FETCH_JOB_RUN_ERROR',
      errors: [
        {
          status: 500,
          detail: 'An error',
        },
      ],
    }

    it('sets loading to false', () => {
      let state = reducer(INITIAL_STATE, fetchJobRunsBeginAction)
      expect(state.jobRuns.loading).toEqual(true)

      state = reducer(state, action)
      expect(state.jobRuns.loading).toEqual(false)
    })

    it('sets error to true & resets it to false when a new fetch starts', () => {
      let state = reducer(INITIAL_STATE, action)
      expect(state.jobRuns.error).toEqual(true)

      state = reducer(state, fetchJobRunsBeginAction)
      expect(state.jobRuns.error).toEqual(false)
    })
  })
})
