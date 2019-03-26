import reducer, { IState } from '../../reducers'
import { JobRunsAction } from '../../reducers/jobRuns'

describe('reducers/search', () => {
  it('returns an initial state', () => {
    const action = {} as JobRunsAction
    const state = reducer({}, action) as IState

    expect(state.jobRuns).toEqual({
      items: undefined
    })
  })

  it('can update the search query', () => {
    const action = { type: 'UPSERT_JOB_RUNS', items: [] } as JobRunsAction
    const state = reducer({}, action) as IState

    expect(state.jobRuns).toEqual({
      items: []
    })
  })
})
