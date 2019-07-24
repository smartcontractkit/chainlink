import { IState } from '../../src/connectors/redux/reducers/index'
import bridgesSelector from '../../src/selectors/bridges'

describe('selectors - bridges', () => {
  it('returns the current page of bridges', () => {
    const state = <IState>{
      bridges: {
        items: {
          a: { attributes: { name: 'A' } },
          b: { attributes: { name: 'B' } },
          c: { attributes: { name: 'C' } }
        },
        currentPage: ['c', 'a']
      }
    }

    const selected = bridgesSelector(state)
    expect(selected).toEqual([{ id: 'c', name: 'C' }, { id: 'a', name: 'A' }])
  })

  it('does not return items that cannot be found', () => {
    const state = <IState>{
      bridges: {
        items: {
          a: { attributes: { name: 'A' } },
          b: { attributes: { name: 'B' } },
          c: { attributes: { name: 'C' } }
        },
        currentPage: ['C', 'A', 'b']
      }
    }

    const selected = bridgesSelector(state)
    expect(selected).toEqual([{ id: 'b', name: 'B' }])
  })
})
