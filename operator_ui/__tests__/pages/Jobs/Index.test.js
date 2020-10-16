/* eslint-env jest */
import createStore from 'createStore'
import { ConnectedIndex as Index } from 'pages/Jobs/Index'
import { mount } from 'enzyme'
import jsonApiJobSpecsFactory from 'factories/jsonApiJobSpecs'
import React from 'react'
import { Provider } from 'react-redux'
import { MemoryRouter } from 'react-router-dom'
import clickNextPage from 'test-helpers/clickNextPage'
import clickPreviousPage from 'test-helpers/clickPreviousPage'
import syncFetch from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'

const classes = {}
const mountIndex = (opts = {}) =>
  mount(
    <Provider store={createStore()}>
      <MemoryRouter>
        <Index classes={classes} pageSize={opts.pageSize} />
      </MemoryRouter>
    </Provider>,
  )

describe('pages/Jobs/Index', () => {
  it('renders the list of jobs', async () => {
    expect.assertions(3)

    const jobSpecsResponse = jsonApiJobSpecsFactory([
      {
        id: 'c60b9927eeae43168ddbe92584937b1b',
        initiators: [{ type: 'web' }],
        createdAt: new Date().toISOString(),
      },
    ])
    global.fetch.getOnce(globPath('/v2/specs'), jobSpecsResponse)

    const wrapper = mountIndex()

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('c60b9927eeae43168ddbe92584937b1b')
    expect(wrapper.text()).toContain('web')
    expect(wrapper.text()).toContain('just now')
  })

  it('can page through the list of jobs', async () => {
    expect.assertions(6)

    const pageOneResponse = jsonApiJobSpecsFactory(
      [{ id: 'ID-ON-FIRST-PAGE' }],
      2,
    )
    global.fetch.getOnce(globPath('/v2/specs'), pageOneResponse)

    const wrapper = mountIndex({ pageSize: 1 })

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('ID-ON-FIRST-PAGE')
    expect(wrapper.text()).not.toContain('ID-ON-SECOND-PAGE')

    const pageTwoResponse = jsonApiJobSpecsFactory(
      [{ id: 'ID-ON-SECOND-PAGE' }],
      2,
    )
    global.fetch.getOnce(globPath('/v2/specs'), pageTwoResponse)
    clickNextPage(wrapper)

    await syncFetch(wrapper)
    expect(wrapper.text()).not.toContain('ID-ON-FIRST-PAGE')
    expect(wrapper.text()).toContain('ID-ON-SECOND-PAGE')

    global.fetch.getOnce(globPath('/v2/specs'), pageOneResponse)
    clickPreviousPage(wrapper)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('ID-ON-FIRST-PAGE')
    expect(wrapper.text()).not.toContain('ID-ON-SECOND-PAGE')
  })
})
