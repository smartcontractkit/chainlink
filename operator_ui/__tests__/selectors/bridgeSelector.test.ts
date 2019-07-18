import { IState } from '../../src/connectors/redux/reducers/index'
import bridgeSelector from '../../src/selectors/bridge'

describe('selectors - bridge', () => {
  it('returns the bridge with the given id', () => {
    const state = <IState>{
      bridges: {
        items: {
          a: { attributes: { name: 'A' } },
          b: { attributes: { name: 'B' } }
        }
      }
    }

    const selected = bridgeSelector(state, 'a')
    expect(selected).toEqual({ id: 'a', name: 'A' })
  })
})
