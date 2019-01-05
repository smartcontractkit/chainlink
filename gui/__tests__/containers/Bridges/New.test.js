/* eslint-env jest */
import React from 'react'
import { Provider } from 'react-redux'
import { mount } from 'enzyme'
import { MemoryRouter } from 'react-router'
import { Switch, Route } from 'react-router-dom'
import createStore from 'connectors/redux'
import New from 'containers/Bridges/New'
import syncFetch from 'test-helpers/syncFetch'
import formikFillIn from 'test-helpers/formikFillIn'

const classes = {}

const TestPrompt = () => <div>Shouldn't be rendered</div>

const mountNew = (store, props) => {
  const NewWithProps = () => <New {...props} />
  return mount(
    <Provider store={store}>
      <MemoryRouter initialEntries={['/bridges/new']}>
        <Switch>
          <Route exact path='/bridges/new' component={NewWithProps} classes={classes} />
          <Route exact path='/' component={TestPrompt} classes={classes} />
        </Switch>
      </MemoryRouter>
    </Provider>
  )
}

describe('containers/Bridges/New', () => {
  it('makes sure all needed fields are entered', async () => {
    expect.assertions(3)
    const store = createStore()
    const wrapper = mountNew(store)
    expect(wrapper.find('form button').getDOMNode().disabled).toEqual(true)
    formikFillIn(wrapper, 'input[name="name"]', 'someRandomBridge', 'name')

    await syncFetch(wrapper)
    expect(wrapper.find('form button').getDOMNode().disabled).toEqual(true)

    formikFillIn(wrapper, 'input[name="url"]', 'https://bridges.com', 'url')
    await syncFetch(wrapper)
    expect(wrapper.find('form button').getDOMNode().disabled).toEqual(false)
  })
})
