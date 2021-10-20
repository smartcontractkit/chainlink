import React from 'react'
import { Route } from 'react-router-dom'
import { JobsShow } from 'pages/Jobs/Show'
import { jsonApiOcrJobSpec } from 'factories/jsonApiJob'
import { jsonApiOcrJobRuns } from 'factories/jsonApiOcrJobRuns'
import jsonApiJobSpecFactory from 'factories/jsonApiJobSpec'
import jsonApiJobSpecRunsFactory from 'factories/jsonApiJobSpecRuns'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { syncFetch } from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'

const JOB_SPEC_ID = 'c60b9927eeae43168ddbe92584937b1b'
const OCR_JOB_SPEC_ID = '1'

describe('pages/Jobs/RecentRuns', () => {
  it('displays a view more link if there are more runs than the display count', async () => {
    const runs = [
      { id: 'runA', jobId: JOB_SPEC_ID },
      { id: 'runB', jobId: JOB_SPEC_ID },
      { id: 'runC', jobId: JOB_SPEC_ID },
      { id: 'runD', jobId: JOB_SPEC_ID },
      { id: 'runE', jobId: JOB_SPEC_ID },
    ]

    const jobSpecResponse = jsonApiJobSpecFactory({
      id: JOB_SPEC_ID,
    })
    const jobRunsResponse = jsonApiJobSpecRunsFactory(runs, 10)

    global.fetch.getOnce(globPath(`/v2/specs/${JOB_SPEC_ID}`), jobSpecResponse)
    global.fetch.getOnce(globPath('/v2/runs'), jobRunsResponse)

    const wrapper = mountWithProviders(
      <Route path="/jobs/:jobSpecId" component={JobsShow} />,
      {
        initialEntries: [`/jobs/${JOB_SPEC_ID}`],
      },
    )

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('View more')
  })

  it('renders OCR job tasks visualisation', async () => {
    const runs = [
      { id: 'runA', jobId: OCR_JOB_SPEC_ID },
      { id: 'runB', jobId: OCR_JOB_SPEC_ID },
      { id: 'runC', jobId: OCR_JOB_SPEC_ID },
      { id: 'runD', jobId: OCR_JOB_SPEC_ID },
      { id: 'runE', jobId: OCR_JOB_SPEC_ID },
    ]

    const taskNames = ['testFetch', 'testParse', 'testMultiply']

    global.fetch.getOnce(
      globPath(`/v2/jobs/${OCR_JOB_SPEC_ID}/runs`),
      jsonApiOcrJobRuns(runs, 10),
    )
    global.fetch.getOnce(
      globPath(`/v2/jobs/${OCR_JOB_SPEC_ID}`),
      jsonApiOcrJobSpec({
        id: OCR_JOB_SPEC_ID,
        dotDagSource: `   ${taskNames[0]}    [type=http method=POST url="http://localhost:8001" requestData="{\\"hi\\": \\"hello\\"}"];\n    ${taskNames[1]}    [type=jsonparse path="data,result"];\n    ${taskNames[2]} [type=multiply times=100];\n    ${taskNames[0]} -\u003e ${taskNames[1]} -\u003e ${taskNames[2]};\n`,
      }),
    )

    const wrapper = mountWithProviders(
      <Route path="/jobs/:jobSpecId" component={JobsShow} />,
      {
        initialEntries: [`/jobs/${OCR_JOB_SPEC_ID}`],
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
      globPath(`/v2/jobs/${OCR_JOB_SPEC_ID}/runs`),
      jsonApiOcrJobRuns(),
    )
    global.fetch.getOnce(
      globPath(`/v2/jobs/${OCR_JOB_SPEC_ID}`),
      jsonApiOcrJobSpec({
        id: OCR_JOB_SPEC_ID,
        dotDagSource: '',
      }),
    )

    const wrapper = mountWithProviders(
      <Route path="/jobs/:jobSpecId" component={JobsShow} />,
      {
        initialEntries: [`/jobs/${OCR_JOB_SPEC_ID}`],
      },
    )

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('Recent job runs')
    expect(wrapper.text()).not.toContain('Task list')
  })
})
