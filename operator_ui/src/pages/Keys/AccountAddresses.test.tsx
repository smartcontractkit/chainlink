import React from 'react'
import { syncFetch } from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { accountBalances } from 'factories/accountBalance'
import { ACCOUNT_BALANCES_ENDPOINT } from 'api/v2/user/balances'
import { AccountAddresses } from './AccountAddresses'

describe('pages/Keys/OcrKeys', () => {
  describe('Off-Chain Reporting keys', () => {
    it('renders the list of keys', async () => {
      const address1 = {
        ethBalance: '1',
        linkBalance: '2',
        address: '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f',
        isFunding: true,
      }

      const address2 = {
        ethBalance: '158',
        linkBalance: '2913',
        address: '0xABCDd2D5E04012C9Ed24C0e513C9bfAa4A2dD77f',
        isFunding: false,
      }
      global.fetch.getOnce(
        globPath(ACCOUNT_BALANCES_ENDPOINT),
        accountBalances([address1, address2]),
      )

      const wrapper = mountWithProviders(<AccountAddresses />)
      await syncFetch(wrapper)
      expect(wrapper.text()).toContain('just now')
      expect(wrapper.find('tbody').children().length).toEqual(2)

      expect(wrapper.text()).toContain(address1.address)
      expect(wrapper.text()).toContain(address1.ethBalance)
      expect(wrapper.text()).toContain(address1.linkBalance)
      expect(wrapper.text()).toContain('Emergency funding')

      expect(wrapper.text()).toContain(address2.address)
      expect(wrapper.text()).toContain(address2.ethBalance)
      expect(wrapper.text()).toContain(address2.linkBalance)
      expect(wrapper.text()).toContain('Regular')
    })
  })
})
