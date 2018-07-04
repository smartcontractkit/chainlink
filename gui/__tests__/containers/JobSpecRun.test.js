import React from 'react'
import createStore from 'connectors/redux'
import syncFetch from 'test-helpers/syncFetch'
import jsonApiJobSpecRunFactory from 'factories/jsonApiJobSpecRun'
import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import { Router } from 'react-static'
import { ConnectedJobSpecRun as JobSpecRun } from 'containers/JobSpecRun'

const classes = {}
const mountJobSpecRun = (props) => (
  mount(
    <Provider store={createStore()}>
      <Router>
        <JobSpecRun classes={classes} {...props} />
      </Router>
    </Provider>
  )
)

describe('containers/JobSpecRun', () => {
  const jobRunId = 'ad24b72c12f441b99b9877bcf6cb506e'

  it('renders the details of the job spec and its latest runs', async () => {
    expect.assertions(4)

    const jobRunResponse = jsonApiJobSpecRunFactory({
      id: 'ad24b72c12f441b99b9877bcf6cb506e',
      createdAt: '2018-06-19T15:39:53.315919143-07:00',
      result: {
        data: {
          value: '0x05070f7f6a40e4ce43be01fa607577432c68730c2cb89a0f50b665e980d926b5'
        }
      }
    })
    global.fetch.getOnce(`/v2/runs/${jobRunId}`, jobRunResponse)

    const props = {match: {params: {jobRunId: jobRunId}}}
    const wrapper = mountJobSpecRun(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('IDad24b72c12f441b99b9877bcf6cb506e')
    expect(wrapper.text()).toContain('Statuscompleted')
    expect(wrapper.text()).toContain('Created2018-06-19T15:39:53.315919143-07:00')
    expect(wrapper.text()).toContain('Result{"value":"0x05070f7f6a40e4ce43be01fa607577432c68730c2cb89a0f50b665e980d926b5"}')
  })
})
