import reducer from 'connectors/redux/reducers'
import { UPSERT_ACCOUNT_BALANCE } from 'connectors/redux/reducers/accountBalances'

describe('connectors/reducers/accountBalances', () => {
  it('returns the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.accountBalances).toEqual({})
  })

  it('UPSERT_ACCOUNT_BALANCE assigns the eth & link balance', () => {
    const previousState = {
      accountBalances: {}
    }
    const action = {
      type: UPSERT_ACCOUNT_BALANCE,
      data: {
        accountBalances: {
          '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f': {
            id: '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f',
            attributes: {
              address: '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f'
            }
          }
        }
      }
    }
    const state = reducer(previousState, action)

    const balance =
      state.accountBalances['0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f']
    expect(balance).toEqual({
      id: '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f',
      attributes: {
        address: '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f'
      }
    })
  })
})
