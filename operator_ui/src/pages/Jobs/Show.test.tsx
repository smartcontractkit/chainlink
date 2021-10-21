import React from 'react'
import { Route } from 'react-router-dom'
import { JobsShow } from 'pages/Jobs/Show'

import { jobRunsAPIResponse } from 'factories/jsonApiOcrJobRun'
import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { syncFetch } from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'

import {
  jsonApiJob,
  fluxMonitorJobResource,
  webJobResource,
} from 'support/factories/jsonApiJobs'

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

    const wrapper = mountWithProviders(
      <Route path="/jobs/:jobId" component={JobsShow} />,
      {
        initialEntries: [`/jobs/${JOB_ID}`],
      },
    )
    await syncFetch(wrapper)

    expect(wrapper.text()).toContain(JOB_ID)
    expect(wrapper.text()).toContain('fluxmonitor')
    expect(wrapper.text()).toContain('Created 1 minute ago')
    expect(wrapper.text()).toContain('runA')
    expect(wrapper.text()).toContain('runB')
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

      const wrapper = mountWithProviders(
        <Route path="/jobs/:jobId" component={JobsShow} />,
        {
          initialEntries: [`/jobs/${JOB_ID}`],
        },
      )
      await syncFetch(wrapper)

      expect(wrapper.find('tbody').first().children().length).toEqual(1)

      // Prep before "Run" button click
      global.fetch.postOnce(
        globPath(`/v2/jobs/${jobResponse.data.attributes.externalJobID}/runs`),
        {},
      )
      global.fetch.getOnce(
        globPath(`/v2/jobs/${JOB_ID}/runs`),
        jobRunsAPIResponse(runs.concat([{ id: 'runB', jobId: JOB_ID }])),
      )

      wrapper.find('Button').find({ children: 'Run' }).first().simulate('click')
      wrapper
        .find('RunJobModal')
        .find({ children: 'Run job' })
        .first()
        .simulate('click')
      await syncFetch(wrapper)
      expect(wrapper.find('tbody').first().children().length).toEqual(2)
    })
  })
})
