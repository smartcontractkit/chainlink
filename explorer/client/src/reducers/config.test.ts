import reducer, { INITIAL_STATE } from '../reducers'
import { FetchJobRunSucceededAction } from '../reducers/actions'

describe('reducers/config', () => {
  it('updates the current etherscan host when job runs are fetched', () => {
    const action: FetchJobRunSucceededAction = {
      type: 'FETCH_JOB_RUN_SUCCEEDED',
      data: {
        chainlinkNodes: [],
        jobRuns: [],
        taskRuns: [],
        meta: {
          jobRun: {
            meta: {
              etherscanHost: 'ropsten.etherscan.io',
            },
          },
        },
      },
    }
    const state = reducer(INITIAL_STATE, action)

    expect(state.config).toEqual({
      etherscanHost: 'ropsten.etherscan.io',
    })
  })
})
