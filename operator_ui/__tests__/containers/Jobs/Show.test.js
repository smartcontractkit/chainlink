import createStore from 'connectors/redux'
import { ConnectedShow as Show } from 'containers/Jobs/Show'
import jsonApiJobSpecFactory from 'factories/jsonApiJobSpec'
import jsonApiJobSpecRunsFactory from 'factories/jsonApiJobSpecRuns'
import React from 'react'
import { Provider } from 'react-redux'
import { MemoryRouter } from 'react-router-dom'
import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'
import mountWithTheme from 'test-helpers/mountWithTheme'
import syncFetch from 'test-helpers/syncFetch'
import { GWEI_PER_TOKEN, WEI_PER_TOKEN } from '../../../src/utils/constants'

const mountShow = props =>
  mountWithTheme(
    <Provider store={createStore()}>
      <MemoryRouter>
        <Show {...props} />
      </MemoryRouter>
    </Provider>
  )

describe.only('containers/Jobs/Show', () => {
  const jobSpecId = 'c60b9927eeae43168ddbe92584937b1b'
  const jobRunId = 'ad24b72c12f441b99b9877bcf6cb506e'
  it('renders the details of the job spec, its latest runs, its task list entries and its total earnings', async () => {
    expect.assertions(9)

    const minuteAgo = isoDate(Date.now() - MINUTE_MS)
    const jobSpecResponse = jsonApiJobSpecFactory({
      id: jobSpecId,
      initiators: [{ type: 'web' }],
      createdAt: minuteAgo,
      earnings: GWEI_PER_TOKEN,
      minPayment: 100 * WEI_PER_TOKEN
    })
    global.fetch.getOnce(`/v2/specs/${jobSpecId}`, jobSpecResponse)

    const jobRunResponse = jsonApiJobSpecRunsFactory([
      {
        id: jobRunId,
        jobId: jobSpecId,
        status: 'pending'
      }
    ])
    global.fetch.getOnce(`begin:/v2/runs`, jobRunResponse)

    const props = { match: { params: { jobSpecId: jobSpecId } } }
    const wrapper = mountShow(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('c60b9927eeae43168ddbe92584937b1b')
    expect(wrapper.text()).toContain('Initiatorweb')
    expect(wrapper.text()).toContain('Created a minute ago')
    expect(wrapper.text()).toContain('1.000000')
    expect(wrapper.text()).toContain('Httpget')
    expect(wrapper.text()).toContain('Run Count1')
    expect(wrapper.text()).toContain('Minimum Payment100 Link')
    expect(wrapper.text()).toContain('Pending')
    expect(wrapper.text()).not.toContain('View More')
  })

  it('displays a view more link if there are more runs than the display count', async () => {
    const runs = [
      { id: 'runA', jobId: jobSpecId },
      { id: 'runB', jobId: jobSpecId },
      { id: 'runC', jobId: jobSpecId }
    ]

    const jobSpecResponse = jsonApiJobSpecFactory({
      id: jobSpecId
    })
    const jobRunsResponse = jsonApiJobSpecRunsFactory(runs)

    global.fetch.getOnce(`begin:/v2/specs/${jobSpecId}`, jobSpecResponse)
    global.fetch.getOnce(`begin:/v2/runs`, jobRunsResponse)

    const props = { match: { params: { jobSpecId: jobSpecId } } }
    const wrapper = mountShow(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('View More')
  })
})
