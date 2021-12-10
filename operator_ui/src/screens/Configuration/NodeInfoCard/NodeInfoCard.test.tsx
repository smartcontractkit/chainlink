import React from 'react'

import { render, screen } from 'support/test-utils'

import { NodeInfoCard } from './NodeInfoCard'

const { queryByText } = screen

describe('NodeInfoCard', () => {
  const originalEnv = process.env

  beforeEach(() => {
    jest.resetModules()
    process.env = {
      ...originalEnv,
      CHAINLINK_VERSION: '1.0.0@6989a388ef26d981e771fec6710dc65bcc8fb5af',
    }
  })

  afterEach(() => {
    process.env = originalEnv
  })

  it('renders the node info card', () => {
    render(<NodeInfoCard />)

    expect(queryByText(/version/i)).toBeInTheDocument()
    expect(queryByText('1.0.0')).toBeInTheDocument()

    expect(queryByText(/sha/i)).toBeInTheDocument()
    expect(
      queryByText('6989a388ef26d981e771fec6710dc65bcc8fb5af'),
    ).toBeInTheDocument()
  })
})
