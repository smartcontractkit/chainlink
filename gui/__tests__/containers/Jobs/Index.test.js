/* eslint-env jest */
import React from 'react'
import jsonApiJobSpecsFactory from 'factories/jsonApiJobSpecs'
import syncFetch from 'test-helpers/syncFetch'
import clickNextPage from 'test-helpers/clickNextPage'
import clickPreviousPage from 'test-helpers/clickPreviousPage'
import createStore from 'connectors/redux'
import { mount } from 'enzyme'
import { MemoryRouter } from 'react-router-dom'
import { Provider } from 'react-redux'
import { ConnectedIndex as Index } from 'containers/Jobs/Index'

const classes = {}
const mountIndex = (opts = {}) => (
  mount(
    <Provider store={createStore()}>
      <MemoryRouter>
        <Index classes={classes} pageSize={opts.pageSize} />
      </MemoryRouter>
    </Provider>
  )
)

describe('containers/Jobs/Index', () => {
  it('renders the list of jobs', async () => {
    expect.assertions(3)

    const jobSpecsResponse = jsonApiJobSpecsFactory([{
      id: 'c60b9927eeae43168ddbe92584937b1b',
      initiators: [{ 'type': 'web' }],
      createdAt: (new Date()).toISOString()
    }])
    global.fetch.getOnce('/v2/specs?page=1&size=10', jobSpecsResponse)

    const wrapper = mountIndex()

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('c60b9927eeae43168ddbe92584937b1b')
    expect(wrapper.text()).toContain('web')
    expect(wrapper.text()).toContain('just now')
  })

  it('can page through the list of jobs', async () => {
    expect.assertions(6)

    const pageOneResponse = jsonApiJobSpecsFactory([{ id: 'ID-ON-FIRST-PAGE' }], 2)
    global.fetch.getOnce('/v2/specs?page=1&size=1', pageOneResponse)

    const wrapper = mountIndex({ pageSize: 1 })

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('ID-ON-FIRST-PAGE')
    expect(wrapper.text()).not.toContain('ID-ON-SECOND-PAGE')

    const pageTwoResponse = jsonApiJobSpecsFactory([{ id: 'ID-ON-SECOND-PAGE' }], 2)
    global.fetch.getOnce('/v2/specs?page=2&size=1', pageTwoResponse)
    clickNextPage(wrapper)

    await syncFetch(wrapper)
    expect(wrapper.text()).not.toContain('ID-ON-FIRST-PAGE')
    expect(wrapper.text()).toContain('ID-ON-SECOND-PAGE')

    global.fetch.getOnce('/v2/specs?page=1&size=1', pageOneResponse)
    clickPreviousPage(wrapper)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('ID-ON-FIRST-PAGE')
    expect(wrapper.text()).not.toContain('ID-ON-SECOND-PAGE')
  })
})
