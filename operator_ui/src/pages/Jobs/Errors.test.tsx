import React from 'react'
import { JobsShow } from 'pages/Jobs/Show'
import { Route } from 'react-router-dom'
import globPath from 'test-helpers/globPath'
import {
  jsonApiJob,
  fluxMonitorJobResource,
} from 'support/factories/jsonApiJobs'

import {
  renderWithRouter,
  screen,
  waitForElementToBeRemoved,
} from 'support/test-utils'

const JOB_ID = '200'

const { getAllByRole, getByText } = screen

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

    renderWithRouter(<Route path="/jobs/:jobId" component={JobsShow} />, {
      initialEntries: [`/jobs/${JOB_ID}/errors`],
    })

    await waitForElementToBeRemoved(() => getByText('Loading...'))

    const rows = getAllByRole('row')
    expect(rows).toHaveLength(4) // Includes the header row

    expect(rows[1]).toHaveTextContent('Unable to start job subscription')
    expect(rows[2]).toHaveTextContent('2nd error')
    expect(rows[3]).toHaveTextContent('3rd error')
  })
})
