import reducer, { INITIAL_STATE } from '../../reducers'
import { FetchJobRunSucceededAction } from '../../reducers/actions'

describe('reducers/config', () => {
  it('can update the search query', () => {
    const action = {
      type: 'FETCH_JOB_RUN_SUCCEEDED',
      data: {
        meta: {
          jobRun: {
            meta: {
              etherscanHost: 'ropsten.etherscan.io',
            },
          },
        },
      },
    } as FetchJobRunSucceededAction
    const state = reducer(INITIAL_STATE, action)

    expect(state.config).toEqual({
      etherscanHost: 'ropsten.etherscan.io',
    })
  })
})
