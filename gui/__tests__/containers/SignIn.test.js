import React from 'react'
import createStore from 'connectors/redux'
import syncFetch from 'test-helpers/syncFetch'
import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import { Switch, Route } from 'react-router-dom'
import { MemoryRouter } from 'react-router'
import SignIn from 'containers/SignIn'
import fillIn from 'test-helpers/fillIn'

const RedirectApp = () => (
  <div>Behind authentication</div>
)
const mountSignIn = (store, props) => (
  mount(
    <Provider store={store}>
      <MemoryRouter initialEntries={['/signin']}>
        <Switch>
          <Route exact path='/signin' component={SignIn} />
          <Route exact path='/' component={RedirectApp} />
        </Switch>
      </MemoryRouter>
    </Provider>
  )
)

const submitForm = (wrapper) => {
  fillIn(wrapper, 'input#email', 'some@email.net')
  fillIn(wrapper, 'input#password', 'abracadabra')
  expect(wrapper.find('form button').getDOMNode().disabled).toEqual(false)
  wrapper.find('form').simulate('submit')
}

describe('containers/SignIn', () => {
  it('unauthenticated user can input credentials and sign in', async () => {
    const store = createStore()
    global.fetch.postOnce(`/sessions`, { authenticated: true })

    const wrapper = mountSignIn(store)
    submitForm(wrapper)

    await syncFetch(wrapper)
    const newState = store.getState()
    expect(newState.authentication.allowed).toEqual(true)
    expect(wrapper.text()).toContain('Behind authentication')
  })

  it('unauthenticated user inputs wrong credentials', async () => {
    const store = createStore()
    global.fetch.postOnce(
      '/sessions',
      { authenticated: false, errors: [{ detail: 'Invalid email' }] },
      { response: { status: 401 } }
    )

    const wrapper = mountSignIn(store)
    submitForm(wrapper)

    await syncFetch(wrapper)

    const newState = store.getState()
    expect(newState.notifications).toEqual({
      errors: [{ detail: 'Your email or password is incorrect. Please try again' }],
      successes: [],
      currentUrl: '/signin'
    })
    expect(newState.authentication.allowed).toEqual(false)
    expect(wrapper.text()).toContain('Your email or password is incorrect. Please try again')
  })
})
