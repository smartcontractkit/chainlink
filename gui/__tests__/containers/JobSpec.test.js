import React from 'react'
import createStore from 'connectors/redux'
import syncFetch from 'test-helpers/syncFetch'
import jsonApiJobSpecFactory from 'factories/jsonApiJobSpec'
import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import { Router } from 'react-static'
import { ConnectedJobSpec as JobSpec } from 'containers/JobSpec'

const classes = {}
const mountJobSpec = (props) => (
  mount(
    <Provider store={createStore()}>
      <Router>
        <JobSpec classes={classes} {...props} />
      </Router>
    </Provider>
  )
)

describe('containers/JobSpec', () => {
  const jobSpecId = 'c60b9927eeae43168ddbe92584937b1b'

  it('renders the details of the job spec and its latest runs', async () => {
    expect.assertions(6)

    const jobSpecResponse = jsonApiJobSpecFactory({
      id: jobSpecId,
      initiators: [{'type': 'web'}],
      createdAt: '2018-05-10T00:41:54.531043837Z',
      runs: [{id: 'runA', result: {data: {value: '8400.00'}}}]
    })
    global.fetch.getOnce(`/v2/specs/${jobSpecId}`, jobSpecResponse)

    const props = {match: {params: {jobSpecId: jobSpecId}}}
    const wrapper = mountJobSpec(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('IDc60b9927eeae43168ddbe92584937b1b')
    expect(wrapper.text()).toContain('Initiatorweb')
    expect(wrapper.text()).toContain('Created2018-05-10T00:41:54.531043837Z')
    expect(wrapper.text()).toContain('Run Count1')
    expect(wrapper.text()).toContain('{"value":"8400.00"}')
    expect(wrapper.text()).not.toContain('Show More')
  })

  it('displays a show more link if there are more than 5 runs', async () => {
    const runs = [
      {id: 'runA', result: {data: {value: '8400.00'}}},
      {id: 'runB', result: {data: {value: '8400.00'}}},
      {id: 'runC', result: {data: {value: '8400.00'}}},
      {id: 'runD', result: {data: {value: '8400.00'}}},
      {id: 'runE', result: {data: {value: '8400.00'}}},
      {id: 'runE', result: {data: {value: '8400.00'}}}
    ]

    const jobSpecResponse = jsonApiJobSpecFactory({
      id: jobSpecId,
      runs: runs
    })

    global.fetch.getOnce(`/v2/specs/${jobSpecId}`, jobSpecResponse)

    const props = {match: {params: {jobSpecId: jobSpecId}}}
    const wrapper = mountJobSpec(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('Show More')
  })
})
