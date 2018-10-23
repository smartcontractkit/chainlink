/* eslint-env jest */
import React from 'react'
import New from 'containers/Jobs/New'
import { Provider } from 'react-redux'
import { mount } from 'enzyme'
import createStore from 'connectors/redux'
import syncFetch from 'test-helpers/syncFetch'
import Form from 'components/Jobs/Form'
import { MemoryRouter } from 'react-router'
import { Switch, Route } from 'react-static'

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
  it('lands correctly', async () => {
    expect.assertions(1)
    let wrapper = mountNew(createStore())

    await syncFetch(wrapper)
    expect(wrapper.contains(<Form />)).toBe(true)
  })
})
