import reducer, { State } from '../../reducers'
import { FetchJobRunSucceededAction } from '../../reducers/actions'

describe('reducers/config', () => {
  it('returns an initial state', () => {
    const action = {} as FetchJobRunSucceededAction
    const state = reducer({}, action) as State

    expect(state.search).toEqual({
      etherscanHost: undefined,
    })
  })

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
    const state = reducer({}, action) as State

    expect(state.config).toEqual({
      etherscanHost: 'ropsten.etherscan.io',
    })
  })
})
