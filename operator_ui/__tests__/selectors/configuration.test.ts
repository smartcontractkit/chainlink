import configurationSelector from '../../src/selectors/configuration'

describe('selectors - configs', () => {
  it('returns a tuple per key/value pair', () => {
    const state = {
      configuration: {
        data: {
          camelCased: 'value',
          key: 'value'
        }
      }
    }

    let expectation = [['CAMEL_CASED', 'value'], ['KEY', 'value']]
    expect(configurationSelector(state)).toEqual(expectation)
  })
})
