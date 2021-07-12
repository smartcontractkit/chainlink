import { AppState } from '../../src/reducers'
import configurationSelector from '../../src/selectors/configuration'

describe('selectors - configs', () => {
  it('returns a tuple per key/value pair', () => {
    const state: Pick<AppState, 'configuration'> = {
      configuration: {
        data: {
          KEY_1: 'value',
          KEY_2: 'value',
        },
      },
    }

    const expectation = [
      ['KEY_1', 'value'],
      ['KEY_2', 'value'],
    ]
    expect(configurationSelector(state)).toEqual(expectation)
  })
})
