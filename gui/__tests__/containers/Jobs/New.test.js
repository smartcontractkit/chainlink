/* eslint-env jest */
import React from 'react'
import { MemoryRouter } from 'react-router'
import { Switch, Route } from 'react-router-dom'
import { Provider } from 'react-redux'
import { mount } from 'enzyme'
import createStore from 'connectors/redux'
import New from 'containers/Jobs/New'
import syncFetch from 'test-helpers/syncFetch'
import formikFillIn from 'test-helpers/formikFillIn'

const classes = {}
const TestPrompt = () => <div>Shouldn't be rendered</div>

const mountNew = (store, props) => {
  const NewWithProps = () => <New {...props} />
  return (mount(
    <Provider store={store}>
      <MemoryRouter initialEntries={['/jobs/new']}>
        <Switch>
          <Route exact path='/jobs/new' component={NewWithProps} classes={classes} />
          <Route exact path='/' component={TestPrompt} classes={classes} />
        </Switch>
      </MemoryRouter>
    </Provider>
  )
  )
}

describe('containers/Jobs/New', () => {
  it('enables the create button when a value is provided', async () => {
    expect.assertions(2)
    const store = createStore()
    const wrapper = mountNew(store)

    await syncFetch(wrapper)
    expect(wrapper.find('form button').getDOMNode().disabled).toEqual(true)

    formikFillIn(wrapper, 'textarea[name="json"]', '{}', 'json')
    await syncFetch(wrapper)
    expect(wrapper.find('form button').getDOMNode().disabled).toEqual(false)
  })
})
