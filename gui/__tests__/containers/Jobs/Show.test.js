import React from 'react'
import createStore from 'connectors/redux'
import syncFetch from 'test-helpers/syncFetch'
import jsonApiJobSpecFactory from 'factories/jsonApiJobSpec'
import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import { Router } from 'react-static'
import { ConnectedShow as Show } from 'containers/Jobs/Show'
import { ConnectedDefinition as Definition } from 'containers/Jobs/Definition'
import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'

const classes = {}
const mountShow = (props) => (
  mount(
    <Provider store={createStore()}>
      <Router>
        <Show classes={classes} {...props} />
      </Router>
    </Provider>
  )
)
const mountDefinition = (props) => (
  mount(
    <Provider store={createStore()}>
      <Router>
        <Definition classes={classes} {...props} />
      </Router>
    </Provider>
  )
)

describe('containers/Jobs/Show', () => {
  const jobSpecId = 'c60b9927eeae43168ddbe92584937b1b'

  it('renders the details of the job spec and its latest runs', async () => {
    expect.assertions(6)

    const minuteAgo = isoDate(Date.now() - MINUTE_MS)
    const jobSpecResponse = jsonApiJobSpecFactory({
      id: jobSpecId,
      initiators: [{ 'type': 'web' }],
      createdAt: minuteAgo,
      runs: [{ id: 'runA', result: { data: { value: '8400.00' } } }]
    })
    global.fetch.getOnce(`/v2/specs/${jobSpecId}`, jobSpecResponse)

    const props = { match: { params: { jobSpecId: jobSpecId } } }
    const wrapper = mountShow(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('IDc60b9927eeae43168ddbe92584937b1b')
    expect(wrapper.text()).toContain('Initiatorweb')
    expect(wrapper.text()).toContain('Createda minute ago')
    expect(wrapper.text()).toContain('Run Count1')
    expect(wrapper.text()).toContain('{"value":"8400.00"}')
    expect(wrapper.text()).not.toContain('Show More')
  })

  it('displays a show more link if there are more runs than the display count', async () => {
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
    expect(wrapper.text()).toContain('Show More')
  })

  it('displays the definition without altering the keys', async () => {
    const jobSpecResponse = jsonApiJobSpecFactory({
      id: jobSpecId,
      initiators: [{ 'type': 'web' }],
      tasks: [
        { type: 'noop', params: { asset_id_base: 'BTC' } }
      ]
    })
    global.fetch.getOnce(`/v2/specs/${jobSpecId}`, jobSpecResponse)

    const props = { match: { params: { jobSpecId: jobSpecId } } }
    const wrapper = mountDefinition(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('asset_id_base')
  })
})
