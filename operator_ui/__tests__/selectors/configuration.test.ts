import { AppState } from 'connectors/redux/reducers'
import configurationSelector from 'selectors/configuration'

describe('selectors - configs', () => {
  it('returns a tuple per key/value pair', () => {
    const state: Pick<AppState, 'configuration'> = {
      configuration: {
        data: {
          camelCased: 'value',
          key: 'value'
        }
      }
    }

    const expectation = [['CAMEL_CASED', 'value'], ['KEY', 'value']]
    expect(configurationSelector(state)).toEqual(expectation)
  })
})
