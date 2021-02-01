import { partialAsFull } from 'support/test-helpers/partialAsFull'
import reducer, { INITIAL_STATE } from '../../src/reducers'
import {
  UpsertAccountBalanceAction,
  ResourceActionType,
} from '../../src/reducers/actions'

describe('reducers/accountBalances', () => {
  it('UPSERT_ACCOUNT_BALANCE assigns the eth & link balance', () => {
    const action = partialAsFull<UpsertAccountBalanceAction>({
      type: ResourceActionType.UPSERT_ACCOUNT_BALANCE,
      data: {
        eThKeys: {
          '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f': {
            id: '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f',
            attributes: {
              address: '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f',
            },
          },
        },
      },
    })
    const state = reducer(INITIAL_STATE, action)

    const balance =
      state.accountBalances['0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f']
    expect(balance).toEqual({
      id: '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f',
      attributes: {
        address: '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f',
      },
    })
  })
})
