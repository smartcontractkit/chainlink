import { JobsShow } from 'pages/Jobs/Show'
import jsonApiJobSpecFactory from 'factories/jsonApiJobSpec'
import jsonApiJobSpecRunsFactory from 'factories/jsonApiJobSpecRuns'
import React from 'react'
import { Route } from 'react-router-dom'
import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { syncFetch } from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'
import { GWEI_PER_TOKEN, WEI_PER_TOKEN } from 'utils/constants'

const JOB_SPEC_ID = 'c60b9927eeae43168ddbe92584937b1b'
const JOB_RUN_ID = 'ad24b72c12f441b99b9877bcf6cb506e'

describe('pages/Jobs/Show', () => {
  it('renders the details of the job spec, its latest runs, its task list entries and its total earnings', async () => {
    const minuteAgo = isoDate(Date.now() - MINUTE_MS)
    const jobSpecResponse = jsonApiJobSpecFactory({
      id: JOB_SPEC_ID,
      initiators: [{ type: 'web' }],
      createdAt: minuteAgo,
      earnings: GWEI_PER_TOKEN,
      minPayment: 100 * WEI_PER_TOKEN,
    })
    global.fetch.getOnce(globPath(`/v2/specs/${JOB_SPEC_ID}`), jobSpecResponse)

    const jobRunResponse = jsonApiJobSpecRunsFactory([
      {
        id: JOB_RUN_ID,
        jobId: JOB_SPEC_ID,
        status: 'pending',
      },
    ])
    global.fetch.getOnce(globPath('/v2/runs'), jobRunResponse)

    const wrapper = mountWithProviders(
      <Route path="/jobs/:jobSpecId" component={JobsShow} />,
      {
        initialEntries: [`/jobs/${JOB_SPEC_ID}`],
      },
    )

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('c60b9927eeae43168ddbe92584937b1b')
    expect(wrapper.text()).toContain('Initiatorweb')
    expect(wrapper.text()).toContain('Created 1 minute ago')
    expect(wrapper.text()).toContain('1.000000')
    expect(wrapper.text()).toContain('Httpget')
    expect(wrapper.text()).toContain('Minimum Payment100 Link')
    expect(wrapper.text()).toContain('Pending')
    expect(wrapper.text()).not.toContain('View more')
  })

  describe('RegionalNav', () => {
    it('clicking on "Run" button triggers a new job and updates the recent job_specs list', async () => {
      const runs = [{ id: 'runA', jobId: JOB_SPEC_ID }]

      global.fetch.getOnce(
        globPath(`/v2/specs/${JOB_SPEC_ID}`),
        jsonApiJobSpecFactory({
          id: JOB_SPEC_ID,
          initiators: [{ type: 'web' }],
        }),
      )
      global.fetch.getOnce(
        globPath('/v2/runs'),
        jsonApiJobSpecRunsFactory(runs),
      )

      const wrapper = mountWithProviders(
        <Route path="/jobs/:jobSpecId" component={JobsShow} />,
        {
          initialEntries: [`/jobs/${JOB_SPEC_ID}`],
        },
      )

      await syncFetch(wrapper)

      expect(wrapper.find('tbody').first().children().length).toEqual(1)

      // Prep before "Run" button click
      global.fetch.postOnce(globPath(`/v2/specs/${JOB_SPEC_ID}/runs`), {})
      global.fetch.getOnce(
        globPath('/v2/runs'),
        jsonApiJobSpecRunsFactory(
          runs.concat([{ id: 'runB', jobId: JOB_SPEC_ID }]),
        ),
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
