import React from 'react'
import createStore from 'connectors/redux'
import syncFetch from 'test-helpers/syncFetch'
import jsonApiJobSpecFactory from 'factories/jsonApiJobSpec'
import mountWithTheme from 'test-helpers/mountWithTheme'
import { Provider } from 'react-redux'
import { MemoryRouter } from 'react-router-dom'
import { ConnectedShow as Show } from 'containers/Jobs/Show'
import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'

const mountShow = props =>
  mountWithTheme(
    <Provider store={createStore()}>
      <MemoryRouter>
        <Show {...props} />
      </MemoryRouter>
    </Provider>
  )

describe('containers/Jobs/Show', () => {
  const jobSpecId = 'c60b9927eeae43168ddbe92584937b1b'

  it('renders the details of the job spec and its latest runs', async () => {
    expect.assertions(6)

    const minuteAgo = isoDate(Date.now() - MINUTE_MS)
    const jobSpecResponse = jsonApiJobSpecFactory({
      id: jobSpecId,
      initiators: [{ type: 'web' }],
      createdAt: minuteAgo,
      runs: [
        {
          id: 'runA',
          status: 'pending',
          result: { data: { value: '8400.00' } }
        }
      ]
    })
    global.fetch.getOnce(`/v2/specs/${jobSpecId}`, jobSpecResponse)

    const props = { match: { params: { jobSpecId: jobSpecId } } }
    const wrapper = mountShow(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('c60b9927eeae43168ddbe92584937b1b')
    expect(wrapper.text()).toContain('Initiatorweb')
    expect(wrapper.text()).toContain('Created a minute ago')
    expect(wrapper.text()).toContain('Run Count1')
    expect(wrapper.text()).toContain('Pending')
    expect(wrapper.text()).not.toContain('View More')
  })

  it('displays a view more link if there are more runs than the display count', async () => {
    const runs = [
      { id: 'runA', jobId: jobSpecId },
      { id: 'runB', jobId: jobSpecId },
      { id: 'runC', jobId: jobSpecId }
    ]

    const jobSpecResponse = jsonApiJobSpecFactory({
      id: jobSpecId,
      runs: runs
    })

    global.fetch.getOnce(`/v2/specs/${jobSpecId}`, jobSpecResponse)

    const props = { match: { params: { jobSpecId: jobSpecId } } }
    const wrapper = mountShow(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('View More')
  })
})
