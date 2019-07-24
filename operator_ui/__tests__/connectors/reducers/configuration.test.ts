import reducer from '../../../src/connectors/redux/reducers'

describe('connectors/reducers/configuration', () => {
  it('returns the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.configuration).toEqual({
      data: {}
    })
  })

  it('UPSERT_CONFIGURATION stores the attribute data', () => {
    const previousState = {
      configuration: {}
    }
    const response = {
      configWhitelists: {
        idA: { attributes: { attributeA: 'ValueA' } }
      }
    }
    const action = {
      type: 'UPSERT_CONFIGURATION',
      data: response
    }
    const state = reducer(previousState, action)

    expect(state.configuration.data).toEqual({ attributeA: 'ValueA' })
  })
})
