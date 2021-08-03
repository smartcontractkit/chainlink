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
    expect.assertions(2)

    const accountBalanceResponse = accountBalances([
      {
        ethBalance: '10123456000000000000000',
        linkBalance: '7467870000000000000000',
      },
    ])
    global.fetch.getOnce(globPath('/v2/keys/eth'), accountBalanceResponse)

    const wrapper = mountIndex()

    await syncFetch(wrapper)

    expect(wrapper.text()).toContain('Link Balance7.467870k')
    expect(wrapper.text()).toContain('Ether Balance10.123456k')
  })
})
