import React from 'react'
import createStore from 'connectors/redux'
import syncFetch from 'test-helpers/syncFetch'
import jsonApiJobSpecRunFactory from 'factories/jsonApiJobSpecRun'
import { Provider } from 'react-redux'
import { MemoryRouter } from 'react-router-dom'
import { ConnectedShow as Show } from 'containers/JobRuns/Show'
import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'
import mountWithTheme from 'test-helpers/mountWithTheme'

const classes = {}
const mountShow = (props) => (
  mountWithTheme(
    <Provider store={createStore()}>
      <MemoryRouter>
        <Show classes={classes} {...props} />
      </MemoryRouter>
    </Provider>
  )
)

describe('containers/JobRuns/Show', () => {
  const jobSpecId = '942e8b218d414e10a053-000455fdd470'
  const jobRunId = 'ad24b72c12f441b99b9877bcf6cb506e'

  it('renders the details of the job spec and its latest runs', async () => {
    expect.assertions(3)

    const minuteAgo = isoDate(Date.now() - MINUTE_MS)
    const jobRunResponse = jsonApiJobSpecRunFactory({
      id: jobRunId,
      createdAt: minuteAgo,
      jobId: jobSpecId,
      initiator: {
        type: 'web',
        params: {}
      },
      taskRuns: [
        { id: 'taskRunA', status: 'completed', task: { type: 'noop', params: {} } }
      ],
      result: {
        data: {
          value: '0x05070f7f6a40e4ce43be01fa607577432c68730c2cb89a0f50b665e980d926b5'
        }
      }
    })
    global.fetch.getOnce(`/v2/runs/${jobRunId}`, jobRunResponse)

    const props = { match: { params: { jobSpecId: jobSpecId, jobRunId: jobRunId } } }
    const wrapper = mountShow(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('Web')
    expect(wrapper.text()).toContain('Noop')
    expect(wrapper.text()).toContain('completed')
  })
})
