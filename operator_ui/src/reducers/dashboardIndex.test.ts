import { partialAsFull } from 'support/test-helpers/partialAsFull'
import reducer, { INITIAL_STATE } from '../../src/reducers'
import {
  UpsertRecentJobRunsAction,
  ResourceActionType,
} from '../../src/reducers/actions'

describe('connectors/reducers/dashboardIndex', () => {
  it('UPSERT_RECENT_JOB_RUNS stores the order of recent runs and total count', () => {
    const action = partialAsFull<UpsertRecentJobRunsAction>({
      type: ResourceActionType.UPSERT_RECENT_JOB_RUNS,
      data: {
        runs: {},
        meta: {
          recentJobRuns: {
            data: [{ id: 'b' }, { id: 'a' }],
            meta: { count: 100 },
          },
        },
      },
    })
    const state = reducer(INITIAL_STATE, action)

    expect(state.dashboardIndex.recentJobRuns).toEqual(['b', 'a'])
    expect(state.dashboardIndex.jobRunsCount).toEqual(100)
  })
})
