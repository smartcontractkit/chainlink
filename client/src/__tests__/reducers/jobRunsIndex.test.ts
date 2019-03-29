import reducer, { IState } from '../../reducers'
import { JobRunsAction } from '../../reducers/jobRuns'

const INITIAL_STATE = { jobRunsIndex: { items: ['replace-me'] } }

describe('reducers/jobRunsIndex', () => {
  it('returns an initial state', () => {
    const action = {} as JobRunsAction
    const state = reducer({}, action) as IState

    expect(state.jobRunsIndex).toEqual({
      items: undefined
    })
  })

  describe('UPSERT_JOB_RUNS', () => {
    it('can replace items', () => {
      const ids = ['9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e']
      const data = { result: ids, entities: {} }
      const action = { type: 'UPSERT_JOB_RUNS', data: data } as JobRunsAction
      const state = reducer(INITIAL_STATE, action) as IState

      expect(state.jobRunsIndex).toEqual({
        items: ['9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e']
      })
    })
  })

  describe('UPSERT_JOB_RUN', () => {
    it('clears items', () => {
      const data = { entities: {} }
      const action = { type: 'UPSERT_JOB_RUN', data: data } as JobRunsAction
      const state = reducer(INITIAL_STATE, action) as IState

      expect(state.jobRunsIndex).toEqual({
        items: undefined
      })
    })
  })
})
