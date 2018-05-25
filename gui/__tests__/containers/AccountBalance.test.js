/* eslint-env jest */
import React from 'react'
import { mount } from 'enzyme'
import { AccountBalance } from 'containers/AccountBalance.js'
import syncFetch from 'test-helpers/syncFetch'
import { decamelizeKeys } from 'humps'

describe('containers/AccountBalance', () => {
  afterEach(global.fetch.reset)

  it('renders the eth & link balance', async () => {
    expect.assertions(2)

    global.fetch.getOnce(
      '/v2/account_balance',
      decamelizeKeys({
        data: {
          attributes: {
            ethBalance: '10120000000000000000000',
            linkBalance: '7460000000000000000000'
          }
        }
      })
    )

    const classes = {}
    const location = {}
    const wrapper = mount(<AccountBalance classes={classes} location={location} />)

    await syncFetch(wrapper).then(() => {
      expect(wrapper.text()).toContain('10.12k')
      expect(wrapper.text()).toContain('7.46k')
    })
  })
})
