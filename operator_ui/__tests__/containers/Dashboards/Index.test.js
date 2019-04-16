/* eslint-env jest */
import React from 'react'
import accountBalanceFactory from 'factories/accountBalance'
import syncFetch from 'test-helpers/syncFetch'
import mountWithTheme from 'test-helpers/mountWithTheme'
import Index from 'containers/Dashboards/Index'

const classes = {}
const mountIndex = (opts = {}) =>
  mountWithTheme(<Index classes={classes} pageSize={opts.pageSize} />)

describe('containers/Dashboards/Index', () => {
  it('renders the recent activity, account balances & recently created jobs', async () => {
    expect.assertions(7)

    const recentJobRuns = {
      data: [
        {
          id: 'runA',
          type: 'runs',
          attributes: {
            id: 'runA',
            createdAt: new Date().toISOString(),
            status: 'completed'
          }
        },
        {
          id: 'runB',
          type: 'runs',
          attributes: {
            id: 'runB',
            createdAt: new Date().toISOString(),
            status: 'completed'
          }
        }
      ]
    }
    global.fetch.getOnce('/v2/runs?sort=-createdAt&size=2', recentJobRuns)

    const recentlyCreatedJobsResponse = {
      data: [
        {
          id: 'job_b',
          type: 'specs',
          attributes: {
            id: 'job_b',
            createdAt: new Date().toISOString()
          }
        },
        {
          id: 'job_a',
          type: 'specs',
          attributes: {
            id: 'job_a',
            createdAt: new Date().toISOString()
          }
        }
      ]
    }
    global.fetch.getOnce(
      '/v2/specs?size=2&sort=-createdAt',
      recentlyCreatedJobsResponse
    )

    const accountBalanceResponse = accountBalanceFactory(
      '10123456000000000000000',
      '7467870000000000000000'
    )
    global.fetch.getOnce('/v2/user/balances', accountBalanceResponse)

    const wrapper = mountIndex()

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('runA')
    expect(wrapper.text()).toContain('runB')

    expect(wrapper.text()).toContain('Link Balance7.467870k')
    expect(wrapper.text()).toContain('Ether Balance10.123456k')

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('Recently Created Jobs')
    expect(wrapper.text()).toContain('job_bCreated just now')
    expect(wrapper.text()).toContain('job_aCreated just now')
  })
})
