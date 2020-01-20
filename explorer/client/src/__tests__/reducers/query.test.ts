import { createStore, applyMiddleware } from 'redux'
import { createQueryMiddleware } from '../../middleware'
import reducer from '../../reducers'
import { UpdateQueryAction } from '../../reducers/actions'

describe('reducers/search', () => {
  it('returns an initial state', () => {
    const middleware = [createQueryMiddleware(location)]
    const store = createStore(reducer, applyMiddleware(...middleware))
    const state = store.getState()

    expect(state.search).toEqual({
      query: undefined,
    })
  })

  it('can parse the search query from a query param', () => {
    const location = {
      toString: () => 'http://localhost/?search=find-me',
    } as Location
    const middleware = [createQueryMiddleware(location)]
    const store = createStore(reducer, applyMiddleware(...middleware))
    const action: UpdateQueryAction = { type: 'QUERY_UPDATED', data: 'BAR' }

    store.dispatch(action)
    const state = store.getState()

    expect(state.search).toEqual({
      query: 'find-me',
    })
  })
})
