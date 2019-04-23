import reducer, { IState } from '../../reducers'
import { SearchAction } from '../../reducers/search'

describe('reducers/search', () => {
  it('returns an initial state', () => {
    const action = {} as SearchAction
    const state = reducer({}, action) as IState

    expect(state.search).toEqual({
      query: undefined
    })
  })

  it('can update the search query', () => {
    const action = {
      type: 'UPDATE_SEARCH_QUERY',
      query: 'something'
    } as SearchAction
    const state = reducer({}, action) as IState

    expect(state.search).toEqual({
      query: 'something'
    })
  })
})
