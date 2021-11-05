import React from 'react'
import { Route } from 'react-router-dom'
import { JobsShow } from 'pages/Jobs/Show'

import { jobRunsAPIResponse } from 'factories/jsonApiOcrJobRun'
import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'
import globPath from 'test-helpers/globPath'

import {
  jsonApiJob,
  fluxMonitorJobResource,
  webJobResource,
} from 'support/factories/jsonApiJobs'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

const { findAllByRole, getByRole, findByText } = screen

const JOB_ID = '200'

describe('pages/Jobs/Show', () => {
  it('renders the details of the job, its latest runs, its task list entries', async () => {
    const minuteAgo = isoDate(Date.now() - MINUTE_MS)

    // Mock the job runs fetch
    const runs = [
      { id: 'runA', jobId: JOB_ID },
      { id: 'runB', jobId: JOB_ID },
    ]
    global.fetch.getOnce(
      globPath(`/v2/jobs/${JOB_ID}/runs`),
      jobRunsAPIResponse(runs, 10),
    )

    // Mock the job fetch
    const jobResponse = jsonApiJob(
      fluxMonitorJobResource({
        id: JOB_ID,
        createdAt: minuteAgo,
      }),
    )
    global.fetch.getOnce(globPath(`/v2/jobs/${JOB_ID}`), jobResponse)

    renderWithRouter(<Route path="/jobs/:jobId" component={JobsShow} />, {
      initialEntries: [`/jobs/${JOB_ID}`],
    })

    expect(await findByText(JOB_ID, { exact: true })).toBeInTheDocument()
    expect(await findByText('fluxmonitor', { exact: true })).toBeInTheDocument()
    expect(
      await findByText('1 minute ago', { exact: true }),
    ).toBeInTheDocument()
    expect(await findByText('runA', { exact: true })).toBeInTheDocument()
    expect(await findByText('runB', { exact: true })).toBeInTheDocument()
  })

  describe('RegionalNav', () => {
    it('clicking on "Run" button triggers a new job and updates the recent job_specs list', async () => {
      // Mock the job runs fetch
      const runs = [{ id: 'runA', jobId: JOB_ID }]
      global.fetch.getOnce(
        globPath(`/v2/jobs/${JOB_ID}/runs`),
        jobRunsAPIResponse(runs, 10),
      )

      // Mock the job fetch
      const jobResponse = jsonApiJob(
        webJobResource({
          id: JOB_ID,
        }),
      )
      global.fetch.getOnce(globPath(`/v2/jobs/${JOB_ID}`), jobResponse)

      renderWithRouter(<Route path="/jobs/:jobId" component={JobsShow} />, {
        initialEntries: [`/jobs/${JOB_ID}`],
      })

      let rows = await findAllByRole('row')
      expect(rows.length).toEqual(1)

      // Prep before "Run" button click
      global.fetch.postOnce(
        globPath(`/v2/jobs/${jobResponse.data.attributes.externalJobID}/runs`),
        { data: { id: '1', attributes: { jobId: JOB_ID } } },
      )
      global.fetch.getOnce(
        globPath(`/v2/jobs/${JOB_ID}/runs`),
        jobRunsAPIResponse(runs.concat([{ id: 'runB', jobId: JOB_ID }])),
      )

      userEvent.click(getByRole('button', { name: 'Run' }))

      // Click the modal confirmation button
      userEvent.click(getByRole('button', { name: 'Run job' }))

      rows = await findAllByRole('row')
      expect(rows.length).toEqual(2)
    })
  })
})
