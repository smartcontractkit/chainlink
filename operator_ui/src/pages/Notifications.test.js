/* eslint-disable react/display-name */
import React from 'react'
import { Provider } from 'react-redux'
import { MemoryRouter } from 'react-router-dom'
import configureStore from 'redux-mock-store'
import Notifications from 'pages/Notifications'
import { render, screen } from '@testing-library/react'

const { queryByText } = screen

const classes = {}
const mockStore = configureStore()

const mountNotifications = (store) =>
  render(
    <Provider store={store}>
      <MemoryRouter>
        <Notifications classes={classes} />
      </MemoryRouter>
    </Provider>,
  )

describe('pages/Notifications', () => {
  it('renders success and error component notifications', () => {
    const successes = [
      {
        component: ({ msg }) => <span>Success {msg}</span>,
        props: { msg: '1' },
      },
    ]
    const errors = [
      { component: ({ msg }) => <span>Error {msg}</span>, props: { msg: '2' } },
    ]
    const state = {
      notifications: {
        successes,
        errors,
        currentUrl: null,
      },
    }
    const store = mockStore(state)
    mountNotifications(store)

    expect(queryByText('Success 1')).toBeInTheDocument()
    expect(queryByText('Error 2')).toBeInTheDocument()
  })

  it('renders success and error text notifications', () => {
    const errors = ['Error Message']
    const successes = ['Success Message']
    const state = {
      notifications: {
        successes,
        errors,
        currentUrl: null,
      },
    }
    const store = mockStore(state)
    mountNotifications(store)

    expect(queryByText('Success Message')).toBeInTheDocument()
    expect(queryByText('Error Message')).toBeInTheDocument()
  })
})
