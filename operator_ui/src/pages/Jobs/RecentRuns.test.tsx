import React from 'react'
import { Route } from 'react-router-dom'
import { JobsShow } from 'pages/Jobs/Show'
import jsonApiJobSpecFactory from 'factories/jsonApiJobSpec'
import jsonApiJobSpecRunsFactory from 'factories/jsonApiJobSpecRuns'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { syncFetch } from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'

const JOB_SPEC_ID = 'c60b9927eeae43168ddbe92584937b1b'

describe('pages/Jobs/RecentRuns', () => {
  it('displays a view more link if there are more runs than the display count', async () => {
    const runs = [
      { id: 'runA', jobId: JOB_SPEC_ID },
      { id: 'runB', jobId: JOB_SPEC_ID },
      { id: 'runC', jobId: JOB_SPEC_ID },
      { id: 'runD', jobId: JOB_SPEC_ID },
      { id: 'runE', jobId: JOB_SPEC_ID },
    ]

    const jobSpecResponse = jsonApiJobSpecFactory({
      id: JOB_SPEC_ID,
    })
    const jobRunsResponse = jsonApiJobSpecRunsFactory(runs, 10)

    global.fetch.getOnce(globPath(`/v2/specs/${JOB_SPEC_ID}`), jobSpecResponse)
    global.fetch.getOnce(globPath('/v2/runs'), jobRunsResponse)

    const wrapper = mountWithProviders(
      <Route path="/jobs/:jobSpecId" component={JobsShow} />,
      {
        initialEntries: [`/jobs/${JOB_SPEC_ID}`],
      },
    )

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('View More')
  })
})
