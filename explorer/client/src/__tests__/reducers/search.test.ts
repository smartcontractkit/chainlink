import reducer, { State } from '../../reducers'
import { Action } from '../../reducers/search'

describe('reducers/search', () => {
  it('returns an initial state', () => {
    const action = {} as Action
    const state = reducer({}, action) as State

    expect(state.search).toEqual({
      query: undefined,
    })
  })

  it('can parse the search query from a query param', () => {
    const location = { toString: () => 'http://localhost/?search=find-me' }
    const action = { location } as Action
    const state = reducer({}, action) as State

    expect(state.search).toEqual({
      query: 'find-me',
    })
  })
})
