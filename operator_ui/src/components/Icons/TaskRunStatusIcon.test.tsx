import React from 'react'

import { render, screen } from '@testing-library/react'

import { TaskRunStatusIcon } from './TaskRunStatusIcon'

const { getByTestId } = screen

describe('TaskRunStatusIcon', () => {
  // Some default values for testing
  const dimensions = { width: 20, height: 20 }

  it('renders the completed icon', () => {
    render(<TaskRunStatusIcon {...dimensions} status="completed" />)
    expect(getByTestId('completed')).toBeInTheDocument()
  })

  it('renders the errored icon', () => {
    render(<TaskRunStatusIcon {...dimensions} status="errored" />)
    expect(getByTestId('errored')).toBeInTheDocument()
  })
})
