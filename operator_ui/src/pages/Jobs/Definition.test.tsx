import React from 'react'
import { JobsShow } from 'pages/Jobs/Show'
import jsonApiJobSpecFactory from 'factories/jsonApiJobSpec'
import { Route } from 'react-router-dom'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { syncFetch } from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'

const JOB_SPEC_ID = 'c60b9927eeae43168ddbe92584937b1b'

describe('pages/Jobs/Definition', () => {
  it('renders the job definition component', async () => {
    global.fetch.getOnce(
      globPath(`/v2/specs/${JOB_SPEC_ID}`),
      jsonApiJobSpecFactory({
        id: JOB_SPEC_ID,
      }),
    )

    const wrapper = mountWithProviders(
      <Route path="/jobs/:jobSpecId" component={JobsShow} />,
      {
        initialEntries: [`/jobs/${JOB_SPEC_ID}/definition`],
      },
    )

    await syncFetch(wrapper)
    expect(wrapper.find('PrettyJson').text()).toMatchSnapshot()
  })
})
