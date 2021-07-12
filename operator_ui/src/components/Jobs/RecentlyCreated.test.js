import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import { render } from 'enzyme'
import RecentlyCreated from 'components/Jobs/RecentlyCreated'
import isoDate, { MINUTE_MS, TWO_MINUTES_MS } from 'test-helpers/isoDate'

const renderComponent = (jobs) =>
  render(
    <MemoryRouter>
      <RecentlyCreated jobs={jobs} />
    </MemoryRouter>,
  )

describe('components/Jobs/RecentlyCreated', () => {
  it('shows the id and creation date', () => {
    const minuteAgo = isoDate(Date.now() - MINUTE_MS)
    const jobB = {
      id: 'job_b',
      createdAt: minuteAgo,
    }
    const twoMinutesAgo = isoDate(Date.now() - TWO_MINUTES_MS)
    const jobA = {
      id: 'job_a',
      createdAt: twoMinutesAgo,
    }

    const wrapper = renderComponent([jobB, jobA])
    expect(wrapper.text()).toContain('job_bCreated 1 minute ago')
    expect(wrapper.text()).toContain('job_aCreated 2 minutes ago')
  })

  it('shows a loading indicator', () => {
    const wrapper = renderComponent(null)
    expect(wrapper.text()).toContain('...')
  })

  it('shows a message for no jobs', () => {
    const wrapper = renderComponent([])
    expect(wrapper.text()).toContain('No recently created jobs')
  })
})
