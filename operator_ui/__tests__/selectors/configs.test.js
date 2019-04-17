import configsSelector from 'selectors/configs'

describe('selectors - configs', () => {
  it('returns a tuple per key/value pair', () => {
    const state = {
      configuration: {
        config: {
          camelCased: 'value',
          key: 'value'
        }
      }
    }

    let expectation = [['CAMEL_CASED', 'value'], ['KEY', 'value']]
    expect(configsSelector(state)).toEqual(expectation)
  })
})
