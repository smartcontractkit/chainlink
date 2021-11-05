/* eslint-env jest */
import React from 'react'

import { Route } from 'react-router-dom'
import {
  renderWithRouter,
  screen,
  waitForElementToBeRemoved,
} from 'support/test-utils'
import { accountBalances } from 'factories/accountBalance'
import globPath from 'test-helpers/globPath'

import { ENDPOINT as JOBS_ENDPOINT } from 'api/v2/jobs'
import { jsonApiJobSpecsV2 } from 'support/factories/jsonApiJobs'
import Index from 'pages/Dashboards/Index'

const { getAllByText, getByText } = screen

describe('pages/Dashboards/Index', () => {
  it('renders the recent activity, account balances & recently created jobs', async () => {
    const accountBalanceResponse = accountBalances([
      {
        ethBalance: '10123456000000000000000',
        linkBalance: '7467870000000000000000',
      },
    ])
    global.fetch.getOnce(globPath('/v2/keys/eth'), accountBalanceResponse)

    // Dummy responses to stop console warns. This should be propoerly tested when we
    // move over to the GQL API.
    global.fetch.getOnce(globPath(JOBS_ENDPOINT), jsonApiJobSpecsV2([]))
    global.fetch.getOnce(globPath('/v2/pipeline/runs'), {
      data: [],
      meta: { count: 0 },
    })

    renderWithRouter(
      <Route path="/">
        <Index />
      </Route>,
      { initialEntries: ['/'] },
    )

    await waitForElementToBeRemoved(() => getAllByText('...'))

    expect(getByText('7.467870k')).toBeInTheDocument()
    expect(getByText('10.123456k')).toBeInTheDocument()
  })
})
