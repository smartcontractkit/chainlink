/* eslint-env jest */
import React from 'react'
import CreateJobSpec from 'containers/CreateJobSpec'
import { Provider } from 'react-redux'
import { mount } from 'enzyme'
import createStore from 'connectors/redux'
import syncFetch from 'test-helpers/syncFetch'
import JobForm from 'components/JobForm'
import { MemoryRouter } from 'react-router'
import { Switch, Route } from 'react-static'

const classes = {}
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
})
