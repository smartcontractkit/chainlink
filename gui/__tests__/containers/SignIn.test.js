import React from 'react'
import createStore from 'connectors/redux'
import syncFetch from 'test-helpers/syncFetch'
import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import { Switch, Route } from 'react-static'
import { MemoryRouter } from 'react-router'
import SignIn from 'containers/SignIn'

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
  wrapper.find('input#email').simulate('change', { target: { value: 'some@email.net' } })
  wrapper.find('input#password').simulate('change', { target: { value: 'abracadabra' } })
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
    expect(newState.session.authenticated).toEqual(true)
    expect(wrapper.text()).toContain('Behind authentication')
  })

  it('unauthenticated user inputs wrong credentials', async () => {
    const store = createStore()
    global.fetch.postOnce(`/sessions`, { authenticated: false, errors: ['Invalid email'] })

    const wrapper = mountSignIn(store)
    submitForm(wrapper)

    await syncFetch(wrapper)

    expect(wrapper.text()).toContain('Invalid email')
    const newState = store.getState()
    expect(newState.session.authenticated).toEqual(false)
    expect(newState.session.errors).toEqual(['Invalid email'])
  })
})
