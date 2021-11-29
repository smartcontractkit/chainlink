import React from 'react'
import globPath from 'test-helpers/globPath'
import { accountBalances } from 'factories/accountBalance'
import { ACCOUNT_BALANCES_ENDPOINT } from 'api/v2/user/balances'
import { AccountAddresses } from './AccountAddresses'
import {
  renderWithRouter,
  screen,
  waitForElementToBeRemoved,
} from 'support/test-utils'

const { findAllByRole, getByText } = screen

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

      renderWithRouter(<AccountAddresses />)

      await waitForElementToBeRemoved(getByText('Loading...'))

      const rows = await findAllByRole('row')
      expect(rows).toHaveLength(3)

      expect(rows[1]).toHaveTextContent(address1.address)
      expect(rows[1]).toHaveTextContent(address1.ethBalance)
      expect(rows[1]).toHaveTextContent(address1.linkBalance)
      expect(rows[1]).toHaveTextContent('Emergency funding')

      expect(rows[2]).toHaveTextContent(address2.address)
      expect(rows[2]).toHaveTextContent(address2.ethBalance)
      expect(rows[2]).toHaveTextContent(address2.linkBalance)
      expect(rows[2]).toHaveTextContent('Regular')
    })
  })
})
