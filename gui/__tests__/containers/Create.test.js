/* eslint-env jest */
import React from 'react'
import { withoutStyles as Create } from 'containers/Create'
import { Provider } from 'react-redux'
import { mount } from 'enzyme'
import createStore from 'connectors/redux'
import syncFetch from 'test-helpers/syncFetch'

const classes = {}

const mountCreatePage = (store, props) =>
  mount(
    <Provider store={(store = createStore())}>
      <Create classes={classes} />
    </Provider>
  )

describe('containers/Create', () => {
  it('renders the create page focused on bridge creation', async () => {
    expect.assertions(3)
    const props = { location: { state: { tab: 0 } } }
    let wrapper = mountCreatePage(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('Create Bridge')
    expect(wrapper.text()).toContain('Build Bridge')
    expect(wrapper.text()).toContain('Type Confirmations')
  })

  it('renders the create page focused on job creation', async () => {
    expect.assertions(3)
    const props = { location: { state: { tab: 1 } } }
    let wrapper = mountCreatePage(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('Create Job')
    expect(wrapper.text()).toContain('Build Job')
    expect(wrapper.text()).toContain('Paste JSON')
  })
})
