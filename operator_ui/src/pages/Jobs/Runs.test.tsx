import React from 'react'
import { Route } from 'react-router-dom'
import { JobsShow } from 'pages/Jobs/Show'
import { jsonApiOcrJobSpec } from 'factories/jsonApiJob'
import { jobRunsAPIResponse } from 'factories/jsonApiOcrJobRun'
import globPath from 'test-helpers/globPath'
import {
  renderWithRouter,
  screen,
  waitForElementToBeRemoved,
} from 'support/test-utils'

const { findAllByRole, getByText } = screen

const JOB_ID = '1'

describe('pages/Jobs/Runs', () => {
  it('renders job runs', async () => {
    const runs = []
    const RUNS_COUNT = 100

    for (let runId = 100; runId >= 1; runId--) {
      runs.push({ id: String(runId) })
    }

    global.fetch.getOnce(
      globPath(`/v2/jobs/${JOB_ID}/runs?page=1&size=10`),
      jobRunsAPIResponse(runs.slice(0, 10), RUNS_COUNT),
    )
    global.fetch.getOnce(
      globPath(`/v2/jobs/${JOB_ID}`),
      jsonApiOcrJobSpec({
        id: JOB_ID,
      }),
    )

    renderWithRouter(<Route path="/jobs/:jobId" component={JobsShow} />, {
      initialEntries: [`/jobs/${JOB_ID}/runs`],
    })

    await waitForElementToBeRemoved(() =>
      getByText('No jobs have been run yet'),
    )

    const rows = await findAllByRole('row')

    expect(rows).toHaveLength(10)
  })
})
