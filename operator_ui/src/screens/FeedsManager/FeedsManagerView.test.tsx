import * as React from 'react'

import globPath from 'test-helpers/globPath'
import { render, screen } from 'support/test-utils'

import { buildFeedsManager } from 'support/factories/feedsManager'
import { FeedsManagerView } from './FeedsManagerView'

const { findByText } = screen

describe('EditFeedsManagerScreen', () => {
  it('renders the feeds manager view', async () => {
    const mgr = buildFeedsManager()

    // Temporary until we switch it out for GQL
    global.fetch.getOnce(globPath('/v2/job_proposals'), { data: [] })

    render(<FeedsManagerView manager={mgr} />)

    expect(await findByText('Feeds Manager')).toBeInTheDocument()
    expect(await findByText('Job Proposals')).toBeInTheDocument()
  })
})
