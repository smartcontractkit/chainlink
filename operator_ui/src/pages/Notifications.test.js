/* eslint-disable react/display-name */
import React from 'react'
import { Provider } from 'react-redux'
import { MemoryRouter } from 'react-router-dom'
import { mount } from 'enzyme'
import configureStore from 'redux-mock-store'
import Notifications from 'pages/Notifications'

const classes = {}
const mockStore = configureStore()

const mountNotifications = (store) =>
  mount(
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
    const store = mockStore(state)
    const wrapper = mountNotifications(store)

    expect(wrapper.text()).toContain('Success 1')
    expect(wrapper.text()).toContain('Error 2')
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
    const wrapper = mountNotifications(store)

    expect(wrapper.text()).toContain('Success Message')
    expect(wrapper.text()).toContain('Error Message')
  })

  it('renders an unhandled error when there is no component', () => {
    const state = {
      notifications: {
        successes: [],
        errors: [{}],
        currentUrl: null,
      },
    }
    const store = mockStore(state)
    const wrapper = mountNotifications(store)

    expect(wrapper.text()).toContain(
      'Unhandled error. Please help us by opening a bug report',
    )
  })
})
