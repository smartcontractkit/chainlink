/* eslint-env jest */
import React from 'react'
import Create from 'containers/Create'
import { Provider } from 'react-redux'
import { mount } from 'enzyme'
import createStore from 'connectors/redux'
import configureStore from 'redux-mock-store'
import syncFetch from 'test-helpers/syncFetch'
import { Switch, Route } from 'react-static'
import { MemoryRouter } from 'react-router'
import JobForm from 'components/JobForm'
import BridgeForm from 'components/BridgeForm'

const classes = {}
const mockStore = configureStore()

const mountCreatePage = (store, props) => {
  const CreateWithProps = () => <Create {...props} />
  return (mount(
    <Provider store={store}>
      <MemoryRouter initialEntries={['/create']}>
        <Switch>
          <Route exact path='/create' component={CreateWithProps} classes={classes} />
          <Route exact path='/create/bridge' component={CreateWithProps} classes={classes} />
          <Route exact path='/create/job' component={CreateWithProps} classes={classes} />
        </Switch>
      </MemoryRouter>
    </Provider>
  )
  )
}

const formikFillIn = (wrapper, selector, value, name) => {
  wrapper.find(selector).simulate('change', { target: {value: value, name: name} })
}

describe('containers/Create', () => {
  it('lands on the default page (bridge create)', async () => {
    expect.assertions(2)
    let wrapper = mountCreatePage(createStore())

    await syncFetch(wrapper)
    expect(wrapper.contains(<JobForm />)).toBe(false)
    expect(wrapper.contains(<BridgeForm />)).toBe(true)
  })

  it('lands on job create tab', async () => {
    expect.assertions(2)
    const props = {match: {params: {structure: 'job'}}}
    let wrapper = mountCreatePage(createStore(), props)

    await syncFetch(wrapper)
    expect(wrapper.contains(<JobForm />)).toBe(true)
    expect(wrapper.contains(<BridgeForm />)).toBe(false)
  })

  it('lands on bridge create tab', async () => {
    expect.assertions(2)
    const props = {match: {params: {structure: 'bridge'}}}
    let wrapper = mountCreatePage(createStore(), props)

    await syncFetch(wrapper)
    expect(wrapper.contains(<JobForm />)).toBe(false)
    expect(wrapper.contains(<BridgeForm />)).toBe(true)
  })

  it('displays warning notification with expired session', async () => {
    expect.assertions(1)
    let wrapper = mountCreatePage(createStore())
    expect(wrapper.text()).toContain('Session expired. Please sign back in')
  })

  it('displays success notification', async () => {
    const state = {
      authentication: { allowed: true },
      create: {errors: [], successMessage: {name: 'randombridgename'}, networkError: false}
    }
    const store = mockStore(state)
    let wrapper = mountCreatePage(store)
    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('Bridge randombridgename was successfully created')
  })

  it('displays error notification', async () => {
    const state = {
      authentication: { allowed: true },
      create: {errors: ['bridge validation: ', 'not allowed'], successMessage: {}, networkError: false}
    }
    const store = mockStore(state)
    let wrapper = mountCreatePage(store)
    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('bridge validation: not allowed')
  })

  it('displays network error notification', async () => {
    const state = {
      authentication: { allowed: true },
      create: {errors: [], successMessage: {}, networkError: true}
    }
    const store = mockStore(state)
    let wrapper = mountCreatePage(store)
    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('Network Error')
  })

  it('makes sure all needed fields are entered', async () => {
    expect.assertions(3)
    const store = createStore()
    const wrapper = mountCreatePage(store)
    expect(wrapper.find('form button').getDOMNode().disabled).toEqual(true)
    formikFillIn(wrapper, 'input#name', 'someRandomBridge', 'name')

    await syncFetch(wrapper)
    expect(wrapper.find('form button').getDOMNode().disabled).toEqual(true)

    formikFillIn(wrapper, 'input#url', 'https://bridges.com', 'url')
    await syncFetch(wrapper)
    expect(wrapper.find('form button').getDOMNode().disabled).toEqual(false)
  })
})
