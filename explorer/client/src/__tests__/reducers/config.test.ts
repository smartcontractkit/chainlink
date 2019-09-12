import reducer, { State } from '../../reducers'
import { Action } from '../../reducers/config'

describe('reducers/config', () => {
  it('returns an initial state', () => {
    const action = {} as Action
    const state = reducer({}, action) as State

    expect(state.search).toEqual({
      etherscanHost: undefined,
    })
  })

  it('can update the search query', () => {
    const action = {
      type: 'UPSERT_JOB_RUN',
      data: {
        meta: {
          jobRun: {
            meta: {
              etherscanHost: 'ropsten.etherscan.io',
            },
          },
        },
      },
    } as Action
    const state = reducer({}, action) as State

    expect(state.config).toEqual({
      etherscanHost: 'ropsten.etherscan.io',
    })
  })
})
