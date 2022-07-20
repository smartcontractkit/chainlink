import * as React from 'react'

import { renderWithRouter, screen } from 'support/test-utils'

import { buildFeedsManagerResultFields } from 'support/factories/gql/fetchFeedsManagersWithProposals'
import { FeedsManagerView } from './FeedsManagerView'
import { MockedProvider } from '@apollo/client/testing'

const { findByText } = screen

describe('FeedsManagerView', () => {
  it('renders the feeds manager view', async () => {
    const mgr = buildFeedsManagerResultFields()

    renderWithRouter(
      <MockedProvider addTypename={false}>
        <FeedsManagerView manager={mgr} />
      </MockedProvider>,
    )

    expect(await findByText('Feeds Manager')).toBeInTheDocument()
    expect(await findByText('Job Proposals')).toBeInTheDocument()
  })
})
