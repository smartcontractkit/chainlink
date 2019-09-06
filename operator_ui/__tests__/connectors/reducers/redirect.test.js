import { RouterActionType } from 'actions'
import reducer from 'connectors/redux/reducers'

describe('connectors/reducers/redirect', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.redirect).toEqual({
      to: null,
    })
  })

  it('REDIRECT sets "to" as the given url', () => {
    const state = reducer(undefined, {
      type: RouterActionType.REDIRECT,
      to: '/foo',
    })

    expect(state.redirect).toEqual({
      to: '/foo',
    })
  })

  it('MATCH_ROUTE sets "to" as null', () => {
    const previousState = {
      redirect: { to: '/foo' },
    }
    const state = reducer(previousState, { type: RouterActionType.MATCH_ROUTE })

    expect(state.redirect).toEqual({
      to: null,
    })
  })
})
