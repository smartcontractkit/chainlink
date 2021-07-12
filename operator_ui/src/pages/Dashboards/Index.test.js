/* eslint-env jest */
import Index from 'pages/Dashboards/Index'
import { accountBalances } from 'factories/accountBalance'
import React from 'react'
import mountWithTheme from 'test-helpers/mountWithTheme'
import syncFetch from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'

const classes = {}
const mountIndex = () => mountWithTheme(<Index classes={classes} />)

describe('pages/Dashboards/Index', () => {
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
            status: 'completed',
          },
        },
        {
          id: 'runB',
          type: 'runs',
          attributes: {
            id: 'runB',
            createdAt: new Date().toISOString(),
            status: 'completed',
          },
        },
      ],
      meta: { count: 2 },
    }
    global.fetch.getOnce(globPath('/v2/runs'), recentJobRuns)

    const recentlyCreatedJobsResponse = {
      data: [
        {
          id: 'job_b',
          type: 'specs',
          attributes: {
            id: 'job_b',
            createdAt: new Date().toISOString(),
          },
        },
        {
          id: 'job_a',
          type: 'specs',
          attributes: {
            id: 'job_a',
            createdAt: new Date().toISOString(),
          },
        },
      ],
    }
    global.fetch.getOnce(globPath('/v2/specs'), recentlyCreatedJobsResponse)

    const accountBalanceResponse = accountBalances([
      {
        ethBalance: '10123456000000000000000',
        linkBalance: '7467870000000000000000',
      },
    ])
    global.fetch.getOnce(globPath('/v2/keys/eth'), accountBalanceResponse)

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
