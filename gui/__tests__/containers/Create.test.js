/* eslint-env jest */
import React from 'react'
import createStore from 'connectors/redux'
import syncFetch from 'test-helpers/syncFetch'
import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import { Router } from 'react-static'
import Create from 'containers/Create'

const classes = {}
const mountCreatePage = () => (
  mount(
    <Provider store={createStore()}>
      <Router>
        <Create classes={classes} />
      </Router>
    </Provider>
  )
)

describe('containers/Create', () => {

  it('renders the create page focused on bridge creation', async () => {
    expect.assertions(3)
    const props = {location: {state: {tab: 0}}}
    const wrapper = mountCreatePage(props)
    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('Create Bridge')
    expect(wrapper.text()).toContain('Build Bridge')
    expect(wrapper.text()).toContain('Type Confirmations')
  })

  it('renders the create page focused on job creation', async () => {
    expect.assertions(3)
    const props = {location: {state: {tab: 1}}}
    const wrapper = mountCreatePage(props)
    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('Create Job')
    expect(wrapper.text()).toContain('Build Job')
    expect(wrapper.text()).toContain('Paste Json')
  })

})
