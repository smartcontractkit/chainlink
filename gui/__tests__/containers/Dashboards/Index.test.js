/* eslint-env jest */
import React from 'react'
import jsonApiJobSpecsFactory from 'factories/jsonApiJobSpecs'
import accountBalanceFactory from 'factories/accountBalance'
import syncFetch from 'test-helpers/syncFetch'
import clickNextPage from 'test-helpers/clickNextPage'
import clickPreviousPage from 'test-helpers/clickPreviousPage'
import createStore from 'connectors/redux'
import { mount } from 'enzyme'
import { Router } from 'react-static'
import { Provider } from 'react-redux'
import { ConnectedIndex as Index } from 'containers/Dashboards/Index'

const classes = {}
const mountIndex = (opts = {}) => (
  mount(
    <Provider store={createStore()}>
      <Router>
        <Index classes={classes} pageSize={opts.pageSize} />
      </Router>
    </Provider>
  )
)

describe('containers/Dashboards/Index', () => {
  it('renders the list of jobs and account balance', async () => {
    expect.assertions(6)

    const jobSpecsResponse = jsonApiJobSpecsFactory([{
      id: 'c60b9927eeae43168ddbe92584937b1b',
      initiators: [{'type': 'web'}],
      createdAt: '2018-05-10T00:41:54.531043837Z'
    }])
    global.fetch.getOnce('/v2/specs?page=1&size=10', jobSpecsResponse)
    const accountBalanceResponse = accountBalanceFactory(
      '10120000000000000000000',
      '7460000000000000000000'
    )
    global.fetch.getOnce('/v2/user/balances', accountBalanceResponse)

    const wrapper = mountIndex()

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('c60b9927eeae43168ddbe92584937b1b')
    expect(wrapper.text()).toContain('web')
    expect(wrapper.text()).toContain('2018-05-10T00:41:54.531043837Z')

    expect(wrapper.text()).toContain('Link Balance7.46k')
    expect(wrapper.text()).toContain('Ether Balance10.12k')

    expect(wrapper.text()).toContain('Jobs1')
  })

  it('can page through the list of jobs', async () => {
    expect.assertions(6)

    const accountBalanceResponse = accountBalanceFactory('0', '0')
    global.fetch.getOnce('/v2/user/balances', accountBalanceResponse)

    const pageOneResponse = jsonApiJobSpecsFactory([{ id: 'ID-ON-FIRST-PAGE' }], 2)
    global.fetch.getOnce('/v2/specs?page=1&size=1', pageOneResponse)

    const wrapper = mountIndex({pageSize: 1})

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

  it('displays an error message when the network requests fail', async () => {
    expect.assertions(3)

    global.fetch.catch(() => { throw new TypeError('Failed to fetch') })

    const wrapper = mountIndex()

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain(
      'There was an error fetching the jobs. Please reload the page.'
    )
    expect(wrapper.text()).toContain(
      'Ether Balanceerror fetching balance'
    )
    expect(wrapper.text()).toContain(
      'Link Balanceerror fetching balance'
    )
  })
})
