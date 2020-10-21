/* eslint-env jest */
import createStore from 'createStore'
import { JobsIndex as Index } from 'pages/Jobs/Index'
import { mount } from 'enzyme'
import jsonApiJobSpecsFactory from 'factories/jsonApiJobSpecs'
import React from 'react'
import { act } from 'react-dom/test-utils'
import { Provider } from 'react-redux'
import { MemoryRouter, Route } from 'react-router-dom'
import clickNextPage from 'test-helpers/clickNextPage'
import clickPreviousPage from 'test-helpers/clickPreviousPage'
import syncFetch from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'

const mountIndex = (opts: { pageSize?: number } = {}) =>
  mount(
    <Provider store={createStore()}>
      <MemoryRouter>
        <Route
          path="/"
          exact
          render={(props) => <Index {...props} pageSize={opts.pageSize} />}
        />
        <Route
          path="/jobs/page/:pageNumber"
          exact
          render={(props) => <Index {...props} pageSize={opts.pageSize} />}
        />
      </MemoryRouter>
    </Provider>,
  )

describe('pages/Jobs/Index', () => {
  it('renders the list of jobs', async () => {
    const jobSpecsResponse = jsonApiJobSpecsFactory([
      {
        id: 'c60b9927eeae43168ddbe92584937b1b',
        initiators: [{ type: 'web' }],
        createdAt: new Date().toISOString(),
      },
    ])
    global.fetch.getOnce(globPath('/v2/specs'), jobSpecsResponse)

    const wrapper = mountIndex()

    await act(async () => {
      await syncFetch(wrapper)
    })
    expect(wrapper.text()).toContain('c60b9927eeae43168ddbe92584937b1b')
    expect(wrapper.text()).toContain('web')
    expect(wrapper.text()).toContain('just now')
  })

  it('can page through the list of jobs', async () => {
    const pageOneResponse = jsonApiJobSpecsFactory(
      [{ id: 'ID-ON-FIRST-PAGE' }],
      2,
    )
    global.fetch.getOnce(globPath('/v2/specs'), pageOneResponse)

    const wrapper = mountIndex({ pageSize: 1 })

    await act(async () => {
      await syncFetch(wrapper)
    })
    expect(wrapper.text()).toContain('ID-ON-FIRST-PAGE')
    expect(wrapper.text()).not.toContain('ID-ON-SECOND-PAGE')

    const pageTwoResponse = jsonApiJobSpecsFactory(
      [{ id: 'ID-ON-SECOND-PAGE' }],
      2,
    )
    global.fetch.getOnce(globPath('/v2/specs'), pageTwoResponse)
    clickNextPage(wrapper)

    await act(async () => {
      await syncFetch(wrapper)
    })
    wrapper.update()

    expect(wrapper.text()).not.toContain('ID-ON-FIRST-PAGE')
    expect(wrapper.text()).toContain('ID-ON-SECOND-PAGE')

    global.fetch.getOnce(globPath('/v2/specs'), pageOneResponse)
    clickPreviousPage(wrapper)

    await act(async () => {
      await syncFetch(wrapper)
    })
    expect(wrapper.text()).toContain('ID-ON-FIRST-PAGE')
    expect(wrapper.text()).not.toContain('ID-ON-SECOND-PAGE')
  })
})
