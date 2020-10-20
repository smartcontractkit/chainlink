import React from 'react'
import { act } from 'react-dom/test-utils'
import { JobsErrors } from 'pages/Jobs/Errors'
import jsonApiJobSpecFactory from 'factories/jsonApiJobSpec'
import { MemoryRouter, Route } from 'react-router-dom'
import mountWithTheme from 'test-helpers/mountWithTheme'
import syncFetch from 'test-helpers/syncFetch'
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
    description: 'Another error',
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

const mountShow = (path: string) =>
  mountWithTheme(
    <MemoryRouter initialEntries={[path]}>
      <Route path="/jobs/:jobSpecId/errors" component={JobsErrors} />
    </MemoryRouter>,
  )

describe('pages/Jobs/Errors', () => {
  it('renders the job spec errors', async () => {
    global.fetch.getOnce(
      globPath(`/v2/specs/${JOB_SPEC_ID}`),
      jsonApiJobSpecFactory({
        id: JOB_SPEC_ID,
        errors,
      }),
    )

    const wrapper = mountShow(`/jobs/${JOB_SPEC_ID}/errors`)

    await act(async () => {
      await syncFetch(wrapper)
      wrapper.update()
    })

    expect(wrapper.text()).toContain('Unable to start job subscription')
    expect(wrapper.find('tbody').children().length).toEqual(3)

    // Set up to delete the job spec error
    global.fetch.deleteOnce(globPath(`/v2/job_spec_errors/*`), {})
    global.fetch.getOnce(
      globPath(`/v2/specs/${JOB_SPEC_ID}`),
      jsonApiJobSpecFactory({
        id: JOB_SPEC_ID,
        errors: errors.slice().splice(0, 1),
      }),
    )

    wrapper.find('table button').first().simulate('click')

    // Check that optimistic delete works
    expect(wrapper.find('tbody').children().length).toEqual(2)

    await act(async () => {
      await syncFetch(wrapper)
      wrapper.update()
    })

    // Check that server delete works
    expect(wrapper.find('tbody').children().length).toEqual(2)
  })
})
