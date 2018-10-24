/* eslint-env jest */
import React from 'react'
import { Provider } from 'react-redux'
import { mount } from 'enzyme'
import configureStore from 'redux-mock-store'
import syncFetch from 'test-helpers/syncFetch'
import Notifications from 'components/Notifications'
import { MemoryRouter } from 'react-router'

const classes = {}
const mockStore = configureStore()

const mountNotifications = (store, props) => {
  return (mount(
    <Provider store={store}>
      <MemoryRouter>
        <Notifications classes={classes} />
      </MemoryRouter>
    </Provider>
  )
  )
}

describe('components/Notifications', () => {
  it('displays errors', async () => {
    const state = {
      notifications: {
        successes: [],
        errors: [{detail: 'Something unexpected happened'}],
        currentUrl: null
      }
    }
    const store = mockStore(state)
    let wrapper = mountNotifications(store)
    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('Something unexpected happened')
  })
})
