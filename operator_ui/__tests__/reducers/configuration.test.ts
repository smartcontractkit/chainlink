import reducer, { INITIAL_STATE } from '../../src/reducers'
import {
  UpsertConfigurationAction,
  ResourceActionType,
} from '../../src/reducers/actions'

describe('reducers/configuration', () => {
  it('UPSERT_CONFIGURATION sets the config attributes', () => {
    const data = {
      configPrinters: {
        idA: { attributes: { attributeA: 'ValueA' } },
      },
    }
    const action: UpsertConfigurationAction = {
      type: ResourceActionType.UPSERT_CONFIGURATION,
      data,
    }
    const state = reducer(INITIAL_STATE, action)

    expect(state.configuration.data).toEqual({ attributeA: 'ValueA' })
  })
})
