import * as React from 'react'

import { render, screen } from 'support/test-utils'

import { buildFeedsManager } from 'support/factories/gql/fetchFeedsManagersWithProposals'
import { FeedsManagerView } from './FeedsManagerView'

const { findByText } = screen

describe('EditFeedsManagerScreen', () => {
  it('renders the feeds manager view', async () => {
    const mgr = buildFeedsManager()

    render(<FeedsManagerView manager={mgr} />)

    expect(await findByText('Feeds Manager')).toBeInTheDocument()
    expect(await findByText('Job Proposals')).toBeInTheDocument()
  })
})
