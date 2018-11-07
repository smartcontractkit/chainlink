import React from 'react'
import createStore from 'connectors/redux'
import syncFetch from 'test-helpers/syncFetch'
import jsonApiJobSpecRunFactory from 'factories/jsonApiJobSpecRun'
import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import { Router } from 'react-static'
import { ConnectedShow as Show } from 'containers/JobRuns/Show'
import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'

const classes = {}
const mountShow = (props) => (
  mount(
    <Provider store={createStore()}>
      <Router>
        <Show classes={classes} {...props} />
      </Router>
    </Provider>
  )
)

describe('containers/JobRuns/Show', () => {
  const jobSpecId = '942e8b218d414e10a053-000455fdd470'
  const jobRunId = 'ad24b72c12f441b99b9877bcf6cb506e'

  it('renders the details of the job spec and its latest runs', async () => {
    expect.assertions(4)

    const minuteAgo = isoDate(Date.now() - MINUTE_MS)
    const jobRunResponse = jsonApiJobSpecRunFactory({
      id: jobRunId,
      createdAt: minuteAgo,
      jobId: jobSpecId,
      result: {
        data: {
          value: '0x05070f7f6a40e4ce43be01fa607577432c68730c2cb89a0f50b665e980d926b5'
        }
      }
    })
    global.fetch.getOnce(`/v2/runs/${jobRunId}`, jobRunResponse)

    const props = {match: {params: {jobSpecId: jobSpecId, jobRunId: jobRunId}}}
    const wrapper = mountShow(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('IDad24b72c12f441b99b9877bcf6cb506e')
    expect(wrapper.text()).toContain('Statuscompleted')
    expect(wrapper.text()).toContain('Createda minute ago')
    expect(wrapper.text()).toContain('Result{"value":"0x05070f7f6a40e4ce43be01fa607577432c68730c2cb89a0f50b665e980d926b5"}')
  })
})
