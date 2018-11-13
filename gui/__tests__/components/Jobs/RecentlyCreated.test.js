import React from 'react'
import { MemoryRouter } from 'react-router'
import { render } from 'enzyme'
import RecentlyCreated from 'components/Jobs/RecentlyCreated'

const renderComponent = (jobs) => (
  render(
    <MemoryRouter>
      <RecentlyCreated jobs={jobs} />
    </MemoryRouter>
  )
)

const createdAt = i => (new Date(i)).toISOString()

describe('components/Jobs/RecentlyCreated', () => {
  it('shows the id and creation date', () => {
    const minuteAgo = createdAt(Date.now() - 60 * 1000)
    const jobB = {
      id: 'job_b',
      createdAt: minuteAgo
    }
    const twoMinutesAgo = createdAt(Date.now() - 60 * 2 * 1000)
    const jobA = {
      id: 'job_a',
      createdAt: twoMinutesAgo
    }

    let wrapper = renderComponent([jobB, jobA])
    expect(wrapper.text()).toContain('job_bCreated a minute ago')
    expect(wrapper.text()).toContain('job_aCreated 2 minutes ago')
  })

  it('shows a loading indicator', () => {
    let wrapper = renderComponent(null)
    expect(wrapper.text()).toContain('...')
  })

  it('shows a message for no jobs', () => {
    let wrapper = renderComponent([])
    expect(wrapper.text()).toContain('No recently created jobs')
  })
})
