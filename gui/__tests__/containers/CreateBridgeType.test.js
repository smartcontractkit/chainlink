/* eslint-env jest */
import React from 'react'
import CreateBridgeType from 'containers/CreateBridgeType'
import { Provider } from 'react-redux'
import { mount } from 'enzyme'
import createStore from 'connectors/redux'
import syncFetch from 'test-helpers/syncFetch'
import BridgeForm from 'components/BridgeForm'
import { MemoryRouter } from 'react-router'
import { Switch, Route } from 'react-static'

const classes = {}

const TestPrompt = () => <div>Shouldn't be rendered</div>

const mountCreatePage = (store, props) => {
  const CreateWithProps = () => <CreateBridgeType {...props} />
  return mount(
    <Provider store={store}>
      <MemoryRouter initialEntries={['/create/bridge']}>
        <Switch>
          <Route exact path='/create/bridge' component={CreateWithProps} classes={classes} />
          <Route exact path='/' component={TestPrompt} classes={classes} />
        </Switch>
      </MemoryRouter>
    </Provider>
  )
}

const formikFillIn = (wrapper, selector, value, name) => {
  wrapper.find(selector).simulate('change', { target: { value: value, name: name } })
}

describe('containers/CreateBridgeType', () => {
  it('lands correctly', async () => {
    expect.assertions(1)
    let wrapper = mountCreatePage(createStore())

    await syncFetch(wrapper)
    expect(wrapper.contains(<BridgeForm />)).toBe(true)
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
