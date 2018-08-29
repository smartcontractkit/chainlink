/* eslint-env jest */
import React from 'react'
import CreateJobSpec from 'containers/CreateJobSpec'
import { Provider } from 'react-redux'
import { mount } from 'enzyme'
import createStore from 'connectors/redux'
import configureStore from 'redux-mock-store'
import syncFetch from 'test-helpers/syncFetch'
import JobForm from 'components/JobForm'
import { MemoryRouter } from 'react-router'
import { Switch, Route } from 'react-static'

const classes = {}
const mockStore = configureStore()
const TestPrompt = () => <div>Shouldn't be rendered</div>

const mountCreatePage = (store, props) => {
  const CreateWithProps = () => <CreateJobSpec {...props} />
  return (mount(
    <Provider store={store}>
      <MemoryRouter initialEntries={['/create/job']}>
        <Switch>
          <Route exact path='/create/job' component={CreateWithProps} classes={classes} />
          <Route exact path='/' component={TestPrompt} classes={classes} />
        </Switch>
      </MemoryRouter>
    </Provider>
  )
  )
}

describe('containers/CreateJobSpec', () => {
  it('lands correctly', async () => {
    expect.assertions(1)
    let wrapper = mountCreatePage(createStore())

    await syncFetch(wrapper)
    expect(wrapper.contains(<JobForm />)).toBe(true)
  })

  it('displays success notification', async () => {
    const state = {
      authentication: { allowed: true },
      create: {successMessage: {id: '83bd80df93f249ef9c8dc8d5a20b34c3'}, networkError: false},
      fetching: {count: 0}
    }
    const store = mockStore(state)
    let wrapper = mountCreatePage(store)
    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('83bd80df93f249ef9c8dc8d5a20b34c3 was successfully created')
  })
})
