import React from 'react'
import { Route } from 'react-router-dom'
import syncFetch from 'test-helpers/syncFetch'
import jsonApiJobSpecRunFactory from 'factories/jsonApiJobSpecRun'
import { Show } from './Overview'
import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import globPath from 'test-helpers/globPath'

describe('pages/JobRuns/Show/Overview', () => {
  const JOB_SPEC_ID = '942e8b218d414e10a053-000455fdd470'
  const JOB_RUN_ID = 'ad24b72c12f441b99b9877bcf6cb506e'

  it('renders the details of the job spec and its latest runs', async () => {
    expect.assertions(3)

    const minuteAgo = isoDate(Date.now() - MINUTE_MS)
    const jobRunResponse = jsonApiJobSpecRunFactory({
      id: JOB_RUN_ID,
      createdAt: minuteAgo,
      jobId: JOB_SPEC_ID,
      initiator: {
        type: 'web',
        params: {},
      },
      taskRuns: [
        {
          id: 'taskRunA',
          status: 'completed',
          task: { type: 'noop', params: {} },
        },
      ],
      result: {
        data: {
          value:
            '0x05070f7f6a40e4ce43be01fa607577432c68730c2cb89a0f50b665e980d926b5',
        },
      },
    })
    global.fetch.getOnce(globPath(`/v2/runs/${JOB_RUN_ID}`), jobRunResponse)

    const wrapper = mountWithProviders(<Route component={Show} />, {
      initialEntries: [`/jobs/${JOB_SPEC_ID}/runs/${JOB_RUN_ID}`],
    })

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('Web')
    expect(wrapper.text()).toContain('Noop')
    expect(wrapper.text()).toContain('Completed')
  })
})
