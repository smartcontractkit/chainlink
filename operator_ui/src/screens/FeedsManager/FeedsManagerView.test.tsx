import * as React from 'react'

import { renderWithRouter, screen } from 'support/test-utils'

import { buildFeedsManagerResultFields } from 'support/factories/gql/fetchFeedsManagersWithProposals'
import { FeedsManagerView } from './FeedsManagerView'

const { findByText } = screen

describe('FeedsManagerScreen', () => {
  it('renders the feeds manager view', async () => {
    const mgr = buildFeedsManagerResultFields()

    renderWithRouter(<FeedsManagerView manager={mgr} />)

    expect(await findByText('Feeds Manager')).toBeInTheDocument()
    expect(await findByText('Job Proposals')).toBeInTheDocument()
  })
})
