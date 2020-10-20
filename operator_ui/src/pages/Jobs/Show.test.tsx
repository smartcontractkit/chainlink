import { act } from 'react-dom/test-utils'
import createStore from 'createStore'
import { JobsShow } from 'pages/Jobs/Show'
import jsonApiJobSpecFactory from 'factories/jsonApiJobSpec'
import jsonApiJobSpecRunsFactory from 'factories/jsonApiJobSpecRuns'
import React from 'react'
import { Provider } from 'react-redux'
import { MemoryRouter, Route } from 'react-router-dom'
import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'
import mountWithTheme from 'test-helpers/mountWithTheme'
import syncFetch from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'
import { GWEI_PER_TOKEN, WEI_PER_TOKEN } from 'utils/constants'

const mountShow = (path: string) =>
  mountWithTheme(
    <Provider store={createStore()}>
      <MemoryRouter initialEntries={[path]}>
        <Route path="/jobs/:jobSpecId" component={JobsShow} />
      </MemoryRouter>
    </Provider>,
  )

describe('pages/Jobs/Show', () => {
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
      minPayment: 100 * WEI_PER_TOKEN,
    })
    global.fetch.getOnce(globPath(`/v2/specs/${jobSpecId}`), jobSpecResponse)

    const jobRunResponse = jsonApiJobSpecRunsFactory([
      {
        id: jobRunId,
        jobId: jobSpecId,
        status: 'pending',
      },
    ])
    global.fetch.getOnce(globPath('/v2/runs'), jobRunResponse)

    const wrapper = mountShow(`/jobs/${jobSpecId}`)

    await act(async () => {
      await syncFetch(wrapper)
    })
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
      { id: 'runC', jobId: jobSpecId },
      { id: 'runD', jobId: jobSpecId },
      { id: 'runE', jobId: jobSpecId },
      { id: 'runF', jobId: jobSpecId },
    ]

    const jobSpecResponse = jsonApiJobSpecFactory({
      id: jobSpecId,
    })
    const jobRunsResponse = jsonApiJobSpecRunsFactory(runs)

    global.fetch.getOnce(globPath(`/v2/specs/${jobSpecId}`), jobSpecResponse)
    global.fetch.getOnce(globPath('/v2/runs'), jobRunsResponse)

    const wrapper = mountShow(`/jobs/${jobSpecId}`)

    await act(async () => {
      await syncFetch(wrapper)
    })
    expect(wrapper.text()).toContain('View More')
  })
})
