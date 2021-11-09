/* eslint-env jest */

import React from 'react'
import { renderWithRouter, screen, waitForElementToBeRemoved } from 'test-utils'
import userEvent from '@testing-library/user-event'
import globPath from 'test-helpers/globPath'
import bridgesFactory from 'factories/bridges'

import { ConnectedIndex as Index } from 'pages/Bridges/Index'

const { getAllByRole, getAllByText, findByRole, getByRole } = screen

describe('pages/Bridges/Index', () => {
  it('renders the list of bridges', async () => {
    const bridgesResponse = bridgesFactory([
      {
        name: 'reggaeIsntThatGood',
        url: 'butbobistho.com',
      },
    ])

    global.fetch.getOnce(globPath('/v2/bridge_types'), bridgesResponse)

    renderWithRouter(<Index pageSize={1} />)

    expect(await findByRole('heading')).toBeInTheDocument()

    const row = getAllByRole('row')[1]
    expect(row).toHaveTextContent('reggaeIsntThatGood')
    expect(row).toHaveTextContent('butbobistho.com')
  })

  it('can page through the list of bridges', async () => {
    // Page 1
    const pageOneResponse = bridgesFactory(
      [{ name: 'ID-ON-FIRST-PAGE', url: 'bridge.com' }],
      2,
    )
    global.fetch.getOnce(globPath('/v2/bridge_types'), pageOneResponse)

    renderWithRouter(<Index pageSize={1} />)

    expect(await findByRole('heading')).toBeInTheDocument()

    let row = getAllByRole('row')[1]
    expect(row).toHaveTextContent('ID-ON-FIRST-PAGE')
    expect(row).toHaveTextContent('bridge.com')

    // Page 2
    const pageTwoResponse = bridgesFactory(
      [{ name: 'ID-ON-SECOND-PAGE', url: 'bridge.com' }],
      2,
    )
    global.fetch.getOnce(globPath('/v2/bridge_types'), pageTwoResponse)

    userEvent.click(getByRole('button', { name: /Next Page/i }))

    await waitForElementToBeRemoved(() => [
      getAllByText('ID-ON-FIRST-PAGE'),
      getAllByText('bridge.com'),
    ])

    row = getAllByRole('row')[1]
    expect(row).toHaveTextContent('ID-ON-SECOND-PAGE')
    expect(row).toHaveTextContent('bridge.com')

    // Page 1
    global.fetch.getOnce(globPath('/v2/bridge_types'), pageOneResponse)

    userEvent.click(getByRole('button', { name: /Previous Page/i }))

    await waitForElementToBeRemoved(() => [
      getAllByText('ID-ON-SECOND-PAGE'),
      getAllByText('bridge.com'),
    ])

    row = getAllByRole('row')[1]
    expect(row).toHaveTextContent('ID-ON-FIRST-PAGE')
    expect(row).toHaveTextContent('bridge.com')
  })
})
