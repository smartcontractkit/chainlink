/* eslint-env jest */
import React from 'react'
import bridgesFactory from 'factories/bridges'
import syncFetch from 'test-helpers/syncFetch'
import clickNextPage from 'test-helpers/clickNextPage'
import clickPreviousPage from 'test-helpers/clickPreviousPage'
import createStore from 'connectors/redux'
import { mount } from 'enzyme'
import { MemoryRouter } from 'react-router-dom'
import { Provider } from 'react-redux'
import { ConnectedIndex as Index } from 'containers/Bridges/Index'

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

describe('containers/Bridges/Index', () => {
  it('renders the list of bridges', async () => {
    expect.assertions(2)

    const bridgesResponse = bridgesFactory([{
      name: 'reggaeIsntThatGood',
      url: 'butbobistho.com'
    }])
    global.fetch.getOnce('/v2/bridge_types?page=1&size=10', bridgesResponse)

    const wrapper = mountIndex()

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('reggaeIsntThatGood')
    expect(wrapper.text()).toContain('butbobistho.com')
  })

  it('can page through the list of bridges', async () => {
    expect.assertions(6)

    const pageOneResponse = bridgesFactory([
      { name: 'ID-ON-FIRST-PAGE', url: 'bridge.com' }
    ], 2)
    global.fetch.getOnce('/v2/bridge_types?page=1&size=1', pageOneResponse)

    const wrapper = mountIndex({ pageSize: 1 })

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('ID-ON-FIRST-PAGE')
    expect(wrapper.text()).not.toContain('ID-ON-SECOND-PAGE')

    const pageTwoResponse = bridgesFactory([
      { name: 'ID-ON-SECOND-PAGE', url: 'bridge.com' }
    ], 2)
    global.fetch.getOnce('/v2/bridge_types?page=2&size=1', pageTwoResponse)
    clickNextPage(wrapper)

    await syncFetch(wrapper)
    expect(wrapper.text()).not.toContain('ID-ON-FIRST-PAGE')
    expect(wrapper.text()).toContain('ID-ON-SECOND-PAGE')

    global.fetch.getOnce('/v2/bridge_types?page=1&size=1', pageOneResponse)
    clickPreviousPage(wrapper)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('ID-ON-FIRST-PAGE')
    expect(wrapper.text()).not.toContain('ID-ON-SECOND-PAGE')
  })

  it('displays an error message when the network requests fail', async () => {
    expect.assertions(1)

    global.fetch.catch(() => { throw new TypeError('Failed to fetch') })

    const wrapper = mountIndex()

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain(
      'There was an error fetching the bridges. Please reload the page.'
    )
  })
})
