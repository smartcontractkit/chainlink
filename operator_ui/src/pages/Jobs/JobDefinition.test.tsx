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

const { getByTestId, getByText } = screen

describe('pages/Jobs/Definition', () => {
  it('renders the job definition component', async () => {
    const jobResponse = jsonApiJob(
      fluxMonitorJobResource({
        id: JOB_ID,
      }),
    )
    global.fetch.getOnce(globPath(`/v2/jobs/${JOB_ID}`), jobResponse)

    renderWithRouter(<Route path="/jobs/:jobId" component={JobsShow} />, {
      initialEntries: [`/jobs/${JOB_ID}/definition`],
    })

    await waitForElementToBeRemoved(() => getByText('Loading...'))

    expect(getByTestId('definition').textContent).toMatchSnapshot()
  })
})
