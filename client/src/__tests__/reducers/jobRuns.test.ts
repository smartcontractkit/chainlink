import reducer, { IState } from '../../reducers'
import { JobRunsAction } from '../../reducers/jobRuns'

const STATE = {
  jobRuns: {
    items: { 'replace-me': { id: 'replace-me' } }
  }
}

describe('reducers/jobRuns', () => {
  it('returns the current state for other actions', () => {
    const action = {} as JobRunsAction
    const state = reducer(STATE, action) as IState

    expect(state.jobRuns).toEqual(STATE.jobRuns)
  })

  it('sets a blank state on default init', () => {
    const action: JobRunsAction = { type: '@@redux/INIT' }
    const state = reducer(STATE, action) as IState

    expect(state.jobRuns).toEqual({
      items: undefined
    })
  })

  it('sets a blank state on dev tools init', () => {
    const action: JobRunsAction = { type: '@@INIT' }
    const state = reducer(STATE, action) as IState

    expect(state.jobRuns).toEqual({
      items: undefined
    })
  })

  describe('UPSERT_JOB_RUNS', () => {
    it('can replace items', () => {
      const jobRuns = {
        '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e': {
          id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e'
        }
      }
      const data = { entities: { jobRuns: jobRuns } }
      const action = { type: 'UPSERT_JOB_RUNS', data: data } as JobRunsAction
      const state = reducer(STATE, action) as IState

      expect(state.jobRuns).toEqual({
        items: {
          '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e': {
            id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e'
          }
        }
      })
    })
  })

  describe('UPSERT_JOB_RUN', () => {
    it('can replace items', () => {
      const jobRuns = {
        '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e': {
          id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e'
        }
      }
      const data = { entities: { jobRuns: jobRuns } }
      const action = { type: 'UPSERT_JOB_RUN', data: data } as JobRunsAction
      const state = reducer(STATE, action) as IState

      expect(state.jobRuns).toEqual({
        items: {
          '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e': {
            id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e'
          }
        }
      })
    })
  })
})
