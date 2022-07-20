import reducer, { INITIAL_STATE } from '../../src/reducers'
import { RouterActionType, RedirectAction } from '../../src/reducers/actions'

describe('reducers/redirect', () => {
  const redirectAction: RedirectAction = {
    type: RouterActionType.REDIRECT,
    to: '/foo',
  }

  it('REDIRECT sets "to" as the given url', () => {
    const state = reducer(INITIAL_STATE, redirectAction)
    expect(state.redirect.to).toEqual('/foo')
  })

  it('MATCH_ROUTE clears "to"', () => {
    let state = reducer(INITIAL_STATE, redirectAction)
    expect(state.redirect.to).toBeDefined()

    state = reducer(state, {
      type: RouterActionType.MATCH_ROUTE,
      pathname: '/any',
    })
    expect(state.redirect.to).toBeUndefined()
  })
})
