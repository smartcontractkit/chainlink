import React from 'react'
import createStore from 'connectors/redux'
import syncFetch from 'test-helpers/syncFetch'
import jsonApiJobSpecRunFactory from 'factories/jsonApiJobSpecRuns'
import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import { Router } from 'react-static'
import { ConnectedJobSpecRuns as JobSpecRuns } from 'containers/JobSpecRuns'

const classes = {}
const mountJobSpecRuns = (props) => (
  mount(
    <Provider store={createStore()}>
      <Router>
        <JobSpecRuns classes={classes} {...props} />
      </Router>
    </Provider>
  )
)

describe('containers/JobSpecRuns', () => {
  const jobSpecId = 'c60b9927eeae43168ddbe92584937b1b'

  it('renders the runs for the job spec', async () => {
    expect.assertions(4)

    const runsResponse = jsonApiJobSpecRunFactory([{
      jobId: jobSpecId
    }])
    global.fetch.getOnce(`/v2/specs/${jobSpecId}/runs`, runsResponse)

    const props = {match: {params: {jobSpecId: jobSpecId}}}
    const wrapper = mountJobSpecRuns(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain(jobSpecId)
    expect(wrapper.text()).toContain(runsResponse.data[0].id)
    expect(wrapper.text()).toContain('completed')
    expect(wrapper.text()).toContain('{"result":"value"}')
  })
})
