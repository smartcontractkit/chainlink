import reducer from 'connectors/redux/reducers'
import { UPSERT_RECENT_JOB_RUNS } from 'connectors/redux/reducers/dashboardIndex'

describe('connectors/reducers/dashboardIndex', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.dashboardIndex).toEqual({
      recentJobRuns: null
    })
  })

  it('UPSERT_RECENT_JOB_RUNS upserts items along with the current page and count', () => {
    const action = {
      type: UPSERT_RECENT_JOB_RUNS,
      data: {
        meta: {
          recentJobRuns: {
            data: [{ id: 'b' }, { id: 'a' }]
          }
        }
      }
    }
    const state = reducer(undefined, action)

    expect(state.dashboardIndex.recentJobRuns).toEqual(['b', 'a'])
  })
})
