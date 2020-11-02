/* eslint-env jest */
import createStore from 'createStore'
import { ConnectedIndex as Index } from 'pages/Bridges/Index'
import { mount } from 'enzyme'
import bridgesFactory from 'factories/bridges'
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

describe('pages/Bridges/Index', () => {
  it('renders the list of bridges', async () => {
    expect.assertions(2)

    const bridgesResponse = bridgesFactory([
      {
        name: 'reggaeIsntThatGood',
        url: 'butbobistho.com',
      },
    ])
    global.fetch.getOnce(globPath('/v2/bridge_types'), bridgesResponse)

    const wrapper = mountIndex()

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('reggaeIsntThatGood')
    expect(wrapper.text()).toContain('butbobistho.com')
  })

  it('can page through the list of bridges', async () => {
    expect.assertions(6)

    const pageOneResponse = bridgesFactory(
      [{ name: 'ID-ON-FIRST-PAGE', url: 'bridge.com' }],
      2,
    )
    global.fetch.getOnce(globPath('/v2/bridge_types'), pageOneResponse)

    const wrapper = mountIndex({ pageSize: 1 })

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('ID-ON-FIRST-PAGE')
    expect(wrapper.text()).not.toContain('ID-ON-SECOND-PAGE')

    const pageTwoResponse = bridgesFactory(
      [{ name: 'ID-ON-SECOND-PAGE', url: 'bridge.com' }],
      2,
    )
    global.fetch.getOnce(globPath('/v2/bridge_types'), pageTwoResponse)
    clickNextPage(wrapper)

    await syncFetch(wrapper)
    expect(wrapper.text()).not.toContain('ID-ON-FIRST-PAGE')
    expect(wrapper.text()).toContain('ID-ON-SECOND-PAGE')

    global.fetch.getOnce(globPath('/v2/bridge_types'), pageOneResponse)
    clickPreviousPage(wrapper)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('ID-ON-FIRST-PAGE')
    expect(wrapper.text()).not.toContain('ID-ON-SECOND-PAGE')
  })
})
