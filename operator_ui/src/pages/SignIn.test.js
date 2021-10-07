import React from 'react'
import createStore from 'createStore'
import syncFetch from 'test-helpers/syncFetch'
import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import { Switch, Route, MemoryRouter } from 'react-router-dom'
import SignIn from 'pages/SignIn'
import fillIn from 'test-helpers/fillIn'
import globPath from 'test-helpers/globPath'

const RedirectApp = () => <div>Behind authentication</div>
const mountSignIn = (store) =>
  mount(
    <Provider store={store}>
      <MemoryRouter initialEntries={['/signin']}>
        <Switch>
          <Route exact path="/signin" component={SignIn} />
          <Route exact path="/" component={RedirectApp} />
        </Switch>
      </MemoryRouter>
    </Provider>,
  )

const submitForm = (wrapper) => {
  fillIn(wrapper, 'input#email', 'some@email.net')
  fillIn(wrapper, 'input#password', 'abracadabra')
  expect(wrapper.find('form button').getDOMNode().disabled).toEqual(false)
  wrapper.find('form').simulate('submit')
}

const AUTHENTICATED_RESPONSE = {
  data: {
    type: 'session',
    id: 'sessionID',
    attributes: { authenticated: true },
  },
}

describe('pages/SignIn', () => {
  it('unauthenticated user can input credentials and sign in', async () => {
    const store = createStore()
    global.fetch.postOnce(globPath('/sessions'), AUTHENTICATED_RESPONSE)

    const wrapper = mountSignIn(store)
    submitForm(wrapper)

    await syncFetch(wrapper)
    const newState = store.getState()
    expect(newState.authentication.allowed).toEqual(true)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('Behind authentication')
  })

  it('unauthenticated user inputs wrong credentials', async () => {
    const store = createStore()
    global.fetch.postOnce(
      globPath('/sessions'),
      { authenticated: false, errors: [{ detail: 'Invalid email' }] },
      { response: { status: 401 } },
    )

    const wrapper = mountSignIn(store)
    submitForm(wrapper)

    await syncFetch(wrapper)

    // Wait a tiny bit for events to propagate through the UI
    await sleep(1)
    function sleep(ms) {
      return new Promise((resolve) => {
        setTimeout(resolve, ms)
      })
    }

    const newState = store.getState()
    expect(newState.notifications).toEqual({
      errors: ['Your email or password is incorrect. Please try again'],
      successes: [],
      currentUrl: undefined,
    })
    expect(newState.authentication.allowed).toEqual(false)
    expect(wrapper.text()).toContain(
      'Your email or password is incorrect. Please try again',
    )
  })
})
