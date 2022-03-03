import React from 'react'

import { render, screen } from '@testing-library/react'

import { JobRunStatusIcon } from './JobRunStatusIcon'

const { getByTestId } = screen

describe('JobRunStatusIcon', () => {
  // Some default values for testing
  const dimensions = { width: 20, height: 20 }

  it('renders the completed icon', () => {
    render(<JobRunStatusIcon {...dimensions} status="COMPLETED" />)
    expect(getByTestId('completed')).toBeInTheDocument()
  })

  it('renders the errored icon', () => {
    render(<JobRunStatusIcon {...dimensions} status="ERRORED" />)
    expect(getByTestId('errored')).toBeInTheDocument()
  })

  it('renders the running icon', () => {
    render(<JobRunStatusIcon {...dimensions} status="RUNNING" />)
    expect(getByTestId('running')).toBeInTheDocument()
  })

  it('renders the suspended icon', () => {
    render(<JobRunStatusIcon {...dimensions} status="SUSPENDED" />)
    expect(getByTestId('suspended')).toBeInTheDocument()
  })

  it('renders no icon when unknown', () => {
    const { container } = render(
      <JobRunStatusIcon {...dimensions} status="UNKNOWN" />,
    )

    expect(container).toBeEmptyDOMElement()
  })
})
