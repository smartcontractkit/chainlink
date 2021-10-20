import React from 'react'
import { JobsShow } from 'pages/Jobs/Show'
import jsonApiJobSpecFactory from 'factories/jsonApiJobSpec'
import { Route } from 'react-router-dom'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { syncFetch } from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'

const JOB_SPEC_ID = 'c60b9927eeae43168ddbe92584937b1b'

const errors = [
  {
    createdAt: '2020-10-16T13:18:46.519087+01:00',
    description: 'Unable to start job subscription',
    id: 1,
    occurrences: 1,
    updatedAt: '2020-10-16T13:18:46.519087+01:00',
  },
  {
    createdAt: '2020-10-16T13:18:46.519087+01:00',
    description: '2nd error',
    id: 2,
    occurrences: 1,
    updatedAt: '2020-10-16T13:18:46.519087+01:00',
  },
  {
    createdAt: '2020-10-16T13:18:46.519087+01:00',
    description: '3rd error',
    id: 3,
    occurrences: 1,
    updatedAt: '2020-10-16T13:18:46.519087+01:00',
  },
]
describe('pages/Jobs/Errors', () => {
  it('renders the job spec errors', async () => {
    global.fetch.getOnce(
      globPath(`/v2/specs/${JOB_SPEC_ID}`),
      jsonApiJobSpecFactory({
        id: JOB_SPEC_ID,
        errors,
      }),
    )

    const wrapper = mountWithProviders(
      <Route path="/jobs/:jobSpecId" component={JobsShow} />,
      {
        initialEntries: [`/jobs/${JOB_SPEC_ID}/errors`],
      },
    )

    await syncFetch(wrapper)

    expect(wrapper.text()).toContain('Unable to start job subscription')
    expect(wrapper.find('tbody').children().length).toEqual(3)

    // Set up to delete the job spec error
    global.fetch.deleteOnce(globPath(`/v2/job_spec_errors/*`), {})
    global.fetch.getOnce(
      globPath(`/v2/specs/${JOB_SPEC_ID}`),
      jsonApiJobSpecFactory({
        id: JOB_SPEC_ID,
        // Intentionally returning 1 result to make sure that we do the correct check
        errors: errors.slice().splice(0, 1),
      }),
    )

    wrapper.find('table button').first().simulate('click')

    // Check that optimistic delete works
    expect(wrapper.find('tbody').children().length).toEqual(2)

    await syncFetch(wrapper)

    expect(wrapper.find('tbody').children().length).toEqual(1)
  })
})
