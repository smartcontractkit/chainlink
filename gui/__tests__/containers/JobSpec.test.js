import React from 'react'
import createStore from 'connectors/redux'
import syncFetch from 'test-helpers/syncFetch'
import jobSpecFactory from 'factories/jobSpec'
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
  it('renders the details of the job spec and its latest runs', async () => {
    expect.assertions(5)

    const jobSpecId = 'c60b9927eeae43168ddbe92584937b1b'
    const jobSpecResponse = jobSpecFactory({
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
  })
})
