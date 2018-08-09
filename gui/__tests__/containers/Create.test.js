/* eslint-env jest */
import React from 'react'
import { withoutStyles as Create } from 'containers/Create'
import { Provider } from 'react-redux'
import { mount } from 'enzyme'
import createStore from 'connectors/redux'
import syncFetch from 'test-helpers/syncFetch'
import { Switch, Route } from 'react-static'
import { MemoryRouter } from 'react-router'

const classes = {}

const mountCreatePage = (store, props) =>
  mount(
    <Provider store={(store = createStore())}>
      <MemoryRouter initialEntries={['/create']}>
        <Switch>
          <Route exact path="/create" component={Create} classes={classes} />
          <Route exact path="/create/bridge" component={Create} classes={classes} />
          <Route exact path="/create/job" component={Create} classes={classes} />
        </Switch>
      </MemoryRouter>
    </Provider>
  )

describe('containers/Create', () => {
  it('renders the bridge create page', async () => {
    expect.assertions(3)
    let wrapper = mountCreatePage()

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('Create Bridge')
    expect(wrapper.text()).toContain('Build Bridge')
    expect(wrapper.text()).toContain('Type Confirmations')
  })

  it('renders the job create page', async () => {
    expect.assertions(3)
    let wrapper = mountCreatePage()

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('Create Job')
    expect(wrapper.text()).toContain('Build Job')
    expect(wrapper.text()).toContain('Paste JSON')
  })
})
