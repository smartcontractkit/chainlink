import * as React from 'react'

import { renderWithRouter, screen } from 'support/test-utils'
import { MockedProvider, MockedResponse } from '@apollo/client/testing'

import { DashboardView } from './DashboardView'

const { findByText } = screen

describe('DashboardView', () => {
  it('renders the cards', async () => {
    const mocks: MockedResponse[] = []

    renderWithRouter(
      <MockedProvider mocks={mocks} addTypename={false}>
        <DashboardView />
      </MockedProvider>,
    )

    expect(await findByText('Activity')).toBeInTheDocument()
    expect(await findByText('Account Balance')).toBeInTheDocument()
    expect(await findByText('Recent Jobs')).toBeInTheDocument()
  })
})
