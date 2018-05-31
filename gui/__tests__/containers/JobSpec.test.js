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
  it('renders the definition and details of the job spec', async () => {
    expect.assertions(3)

    const jobSpecId = 'c60b9927eeae43168ddbe92584937b1b'
    const jobSpecResponse = jobSpecFactory({
      id: jobSpecId,
      initiators: [{'type': 'web'}],
      createdAt: '2018-05-10T00:41:54.531043837Z'
    })
    global.fetch.getOnce(`/v2/specs/${jobSpecId}`, jobSpecResponse)

    const props = {match: {params: {jobSpecId: jobSpecId}}}
    const wrapper = mountJobSpec(props)

    await syncFetch(wrapper).then(() => {
      expect(wrapper.text()).toContain('IDc60b9927eeae43168ddbe92584937b1b')
      expect(wrapper.text()).toContain('Initiatorweb')
      expect(wrapper.text()).toContain('Created2018-05-10T00:41:54.531043837Z')
    })
  })
})
