import React from 'react'
import { Route } from 'react-router-dom'
import { JobsShow } from 'pages/Jobs/Show'
import { jsonApiOcrJobSpec } from 'factories/jsonApiJob'
import { jobRunsAPIResponse } from 'factories/jsonApiOcrJobRun'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { syncFetch } from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'

const JOB_ID = '1'

describe('pages/Jobs/RecentRuns', () => {
  it('renders job tasks visualisation', async () => {
    const runs = [
      { id: 'runA', jobId: JOB_ID },
      { id: 'runB', jobId: JOB_ID },
      { id: 'runC', jobId: JOB_ID },
      { id: 'runD', jobId: JOB_ID },
      { id: 'runE', jobId: JOB_ID },
    ]

    const taskNames = ['testFetch', 'testParse', 'testMultiply']

    global.fetch.getOnce(
      globPath(`/v2/jobs/${JOB_ID}/runs`),
      jobRunsAPIResponse(runs, 10),
    )
    global.fetch.getOnce(
      globPath(`/v2/jobs/${JOB_ID}`),
      jsonApiOcrJobSpec({
        id: JOB_ID,
        dotDagSource: `   ${taskNames[0]}    [type=http method=POST url="http://localhost:8001" requestData="{\\"hi\\": \\"hello\\"}"];\n    ${taskNames[1]}    [type=jsonparse path="data,result"];\n    ${taskNames[2]} [type=multiply times=100];\n    ${taskNames[0]} -\u003e ${taskNames[1]} -\u003e ${taskNames[2]};\n`,
      }),
    )

    const wrapper = mountWithProviders(
      <Route path="/jobs/:jobId" component={JobsShow} />,
      {
        initialEntries: [`/jobs/${JOB_ID}`],
      },
    )

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('View more')
    expect(wrapper.text()).toContain(taskNames[0])
    expect(wrapper.text()).toContain(taskNames[1])
    expect(wrapper.text()).toContain(taskNames[2])
  })

  it('works with no tasks (bootstrap node)', async () => {
    global.fetch.getOnce(
      globPath(`/v2/jobs/${JOB_ID}/runs`),
      jobRunsAPIResponse(),
    )
    global.fetch.getOnce(
      globPath(`/v2/jobs/${JOB_ID}`),
      jsonApiOcrJobSpec({
        id: JOB_ID,
        dotDagSource: '',
      }),
    )

    const wrapper = mountWithProviders(
      <Route path="/jobs/:jobId" component={JobsShow} />,
      {
        initialEntries: [`/jobs/${JOB_ID}`],
      },
    )

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('Recent job runs')
    expect(wrapper.text()).not.toContain('Task list')
  })
})
