import React from 'react'
import RegionalNav from './RegionalNav'
import { renderWithRouter, screen } from 'support/test-utils'

const { queryByText } = screen

describe('pages/Jobs/Runs/RegionalNav', () => {
  it('displays an overview & json tab by default', () => {
    renderWithRouter(<RegionalNav jobRunId="1" jobId="1" />)

    expect(queryByText('Overview')).toBeInTheDocument()
    expect(queryByText('JSON')).toBeInTheDocument()
    expect(queryByText('Error Log')).toBeNull()
  })
})
