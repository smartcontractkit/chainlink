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

const errors = [
  {
    createdAt: '2020-10-16T13:18:46.519087+01:00',
    description: 'Unable to start job subscription',
    id: '1',
    occurrences: 1,
    updatedAt: '2020-10-16T13:18:46.519087+01:00',
  },
  {
    createdAt: '2020-10-16T13:18:46.519087+01:00',
    description: '2nd error',
    id: '2',
    occurrences: 1,
    updatedAt: '2020-10-16T13:18:46.519087+01:00',
  },
  {
    createdAt: '2020-10-16T13:18:46.519087+01:00',
    description: '3rd error',
    id: '3',
    occurrences: 1,
    updatedAt: '2020-10-16T13:18:46.519087+01:00',
  },
]
describe('pages/Jobs/Errors', () => {
  it('renders the job spec errors', async () => {
    const jobResponse = jsonApiJob(
      fluxMonitorJobResource({
        id: JOB_ID,
        errors,
      }),
    )
    global.fetch.getOnce(globPath(`/v2/jobs/${JOB_ID}`), jobResponse)

    const wrapper = mountWithProviders(
      <Route path="/jobs/:jobId" component={JobsShow} />,
      {
        initialEntries: [`/jobs/${JOB_ID}/errors`],
      },
    )

    await syncFetch(wrapper)

    expect(wrapper.text()).toContain('Unable to start job subscription')
    expect(wrapper.find('tbody').children().length).toEqual(3)
  })
})
