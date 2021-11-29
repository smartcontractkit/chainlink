import React from 'react'
import { Route } from 'react-router-dom'
import { JobsShow } from 'pages/Jobs/Show'
import { jsonApiOcrJobSpec } from 'factories/jsonApiJob'
import { jobRunsAPIResponse } from 'factories/jsonApiOcrJobRun'
import globPath from 'test-helpers/globPath'
import { renderWithRouter, screen } from 'support/test-utils'

const { findByText, queryByText } = screen

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

    renderWithRouter(<Route path="/jobs/:jobId" component={JobsShow} />, {
      initialEntries: [`/jobs/${JOB_ID}`],
    })

    expect(await findByText('View more')).toBeInTheDocument()
    expect(await findByText(taskNames[0])).toBeInTheDocument()
    expect(await findByText(taskNames[1])).toBeInTheDocument()
    expect(await findByText(taskNames[2])).toBeInTheDocument()
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

    renderWithRouter(<Route path="/jobs/:jobId" component={JobsShow} />, {
      initialEntries: [`/jobs/${JOB_ID}`],
    })

    expect(await findByText('Recent job runs')).toBeInTheDocument()
    expect(queryByText('Task list')).toBeNull()
  })
})
