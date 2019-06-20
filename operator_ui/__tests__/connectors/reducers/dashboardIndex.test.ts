import reducer from '../../../src/connectors/redux/reducers'
import { Action } from '../../../src/connectors/redux/reducers/dashboardIndex'

describe('connectors/reducers/dashboardIndex', () => {
  it('returns an initial state', () => {
    const state = reducer(undefined, {} as Action)

    expect(state.dashboardIndex).toEqual({
      recentJobRuns: undefined,
      jobRunsCount: undefined
    })
  })

  it('UPSERT_RECENT_JOB_RUNS stores the order of recent runs and total count', () => {
    const action = {
      type: 'UPSERT_RECENT_JOB_RUNS',
      data: {
        meta: {
          recentJobRuns: {
            data: [{ id: 'b' }, { id: 'a' }],
            meta: { count: 100 }
          }
        }
      }
    }
    const state = reducer(undefined, action)

    expect(state.dashboardIndex.recentJobRuns).toEqual(['b', 'a'])
    expect(state.dashboardIndex.jobRunsCount).toEqual(100)
  })
})
