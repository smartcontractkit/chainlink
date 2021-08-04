import React from 'react'
import { JobsShow } from 'pages/Jobs/Show'
import { Route } from 'react-router-dom'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { syncFetch } from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'

import {
  jsonApiJob,
  fluxMonitorJobResource,
} from 'support/factories/jsonApiJobs'

const JOB_ID = '200'

describe('pages/Jobs/Definition', () => {
  it('renders the job definition component', async () => {
    // Mock the job fetch
    const jobResponse = jsonApiJob(
      fluxMonitorJobResource({
        id: JOB_ID,
      }),
    )
    global.fetch.getOnce(globPath(`/v2/jobs/${JOB_ID}`), jobResponse)

    const wrapper = mountWithProviders(
      <Route path="/jobs/:jobId" component={JobsShow} />,
      {
        initialEntries: [`/jobs/${JOB_ID}/definition`],
      },
    )

    await syncFetch(wrapper)

    expect(wrapper.find('code').text()).toMatchSnapshot()
  })
})
