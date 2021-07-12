import React from 'react'
import { Route } from 'react-router-dom'
import { JobsShow } from 'pages/Jobs/Show'
import { jsonApiOcrJobSpec } from 'factories/jsonApiJob'
import { jsonApiOcrJobRuns } from 'factories/jsonApiOcrJobRuns'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { syncFetch } from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'

const JOB_SPEC_ID = '1'

describe('pages/Jobs/Runs', () => {
  it('adds pagination info to the URL and renders job runs', async () => {
    const runs = []
    const RUNS_COUNT = 100

    for (let runId = 100; runId >= 1; runId--) {
      runs.push({ id: String(runId) })
    }

    global.fetch.getOnce(
      globPath(`/v2/jobs/${JOB_SPEC_ID}/runs?page=1&size=10`),
      jsonApiOcrJobRuns(runs.slice(0, 10), RUNS_COUNT),
    )
    global.fetch.getOnce(
      globPath(`/v2/jobs/${JOB_SPEC_ID}`),
      jsonApiOcrJobSpec({
        id: JOB_SPEC_ID,
      }),
    )

    const wrapper = mountWithProviders(
      <Route path="/jobs/:jobSpecId" component={JobsShow} />,
      {
        initialEntries: [`/jobs/${JOB_SPEC_ID}/runs`],
      },
    )
    await syncFetch(wrapper)

    expect(wrapper.find('tr').length).toEqual(10)

    const routerCopmonentProps: any = wrapper.find('Router').props()
    expect(routerCopmonentProps?.history?.location?.search).toEqual(
      '?page=1&size=10',
    )
  })
})
