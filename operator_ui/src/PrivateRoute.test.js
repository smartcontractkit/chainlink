/* eslint-env jest */
import React from 'react'
import configureStore from 'redux-mock-store'
import { Provider } from 'react-redux'
import { Route, Switch } from 'react-router-dom'
import { MemoryRouter } from 'react-router-dom'
import PrivateRoute from 'PrivateRoute'
import { render, screen } from '@testing-library/react'

const { getByText } = screen

const mockStore = configureStore()
const PrivatePage = () => <div>Behind authentication</div>
const AuthenticatedApp = () => (
  <Switch>
    <PrivateRoute exact path="/" component={PrivatePage} />
    <Route path="/signin">Redirect Success</Route>
  </Switch>
)
const mountAuthenticatedApp = (store) =>
  render(
    <Provider store={store}>
      <MemoryRouter initialEntries={['/']}>
        <Route component={AuthenticatedApp} />
      </MemoryRouter>
    </Provider>,
  )

describe('PrivateRoute', () => {
  it('redirects when user is NOT autheticated', async () => {
    const state = { authentication: { allowed: false } }
    mountAuthenticatedApp(mockStore(state))

    expect(getByText('Redirect Success')).toBeInTheDocument()
  })

  it('goes to destination when user is autheticated', async () => {
    const state = { authentication: { allowed: true } }
    mountAuthenticatedApp(mockStore(state))

    expect(getByText('Behind authentication')).toBeInTheDocument()
  })
})
