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
  it('renders the list of jobs, account balances & recently created jobs', async () => {
    expect.assertions(8)

    const recentlyCreatedJobsResponse = {
      data: [
        {
          id: 'job_b',
          type: 'specs',
          attributes: {
            id: 'job_b',
            createdAt: (new Date()).toISOString()
          }
        },
        {
          id: 'job_a',
          type: 'specs',
          attributes: {
            id: 'job_a',
            createdAt: (new Date()).toISOString()
          }
        }
      ]
    }
    global.fetch.getOnce('/v2/specs?size=2&sort=-createdAt', recentlyCreatedJobsResponse)

    const jobSpecsResponse = jsonApiJobSpecsFactory([{
      id: 'c60b9927eeae43168ddbe92584937b1b',
      initiators: [{'type': 'web'}],
      createdAt: (new Date()).toISOString()
    }])
    global.fetch.getOnce('/v2/specs?page=1&size=10', jobSpecsResponse)
    const accountBalanceResponse = accountBalanceFactory(
      '10123456000000000000000',
      '7467870000000000000000'
    )
    global.fetch.getOnce('/v2/user/balances', accountBalanceResponse)

    const wrapper = mountIndex()

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('c60b9927eeae43168ddbe92584937b1b')
    expect(wrapper.text()).toContain('web')
    expect(wrapper.text()).toContain('just now')

    expect(wrapper.text()).toContain('Link Balance7.467870k')
    expect(wrapper.text()).toContain('Ether Balance10.123456k')

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('Recently Created Jobs')
    expect(wrapper.text()).toContain('job_bCreated just now')
    expect(wrapper.text()).toContain('job_aCreated just now')
  })

  it('can page through the list of jobs', async () => {
    expect.assertions(6)

    const pageOneResponse = jsonApiJobSpecsFactory([{ id: 'ID-ON-FIRST-PAGE' }], 2)
    global.fetch.getOnce('/v2/specs?page=1&size=1', pageOneResponse)

    const accountBalanceResponse = accountBalanceFactory('0', '0')
    global.fetch.getOnce('/v2/user/balances', accountBalanceResponse)

    const recentlyCreatedJobsResponse = {data: []}
    global.fetch.getOnce('/v2/specs?size=2\u0026sort=-createdAt', recentlyCreatedJobsResponse)

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
})
