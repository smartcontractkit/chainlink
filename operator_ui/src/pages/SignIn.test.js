import React from 'react'
import { Route } from 'react-router-dom'
import SignIn from 'pages/SignIn'
import globPath from 'test-helpers/globPath'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

const { findAllByText, findByText, getByLabelText, getByRole } = screen

const RedirectApp = () => <div>Behind authentication</div>
const mountSignIn = () =>
  renderWithRouter(
    <>
      <Route exact path="/signin" component={SignIn} />
      <Route exact path="/" component={RedirectApp} />
    </>,
    { initialEntries: ['/signin'] },
  )

const submitForm = () => {
  userEvent.type(getByRole('textbox', { name: /email/i }), 'some@email.net')
  userEvent.type(getByLabelText(/password/i), 'abracadabra')

  userEvent.click(getByRole('button', { name: 'Access Account' }))
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
    global.fetch.postOnce(globPath('/sessions'), AUTHENTICATED_RESPONSE)

    mountSignIn()
    submitForm()

    expect(await findByText('Behind authentication')).toBeInTheDocument()
  })

  it('unauthenticated user inputs wrong credentials', async () => {
    // const store = createStore()
    global.fetch.postOnce(
      globPath('/sessions'),
      { authenticated: false, errors: [{ detail: 'Invalid email' }] },
      { response: { status: 401 } },
    )

    mountSignIn()
    submitForm()

    // Because our notifications are setup in the renderWithRouter and it is also
    // set up in sign in, this renders twice in tests.
    //
    // TODO - Fix this when we refactor to GQL
    expect(
      await findAllByText(
        'Your email or password is incorrect. Please try again',
      ),
    ).toHaveLength(2)
  })
})
