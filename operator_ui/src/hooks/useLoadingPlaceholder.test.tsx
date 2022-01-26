import React from 'react'
import { render, screen } from '@testing-library/react'
import { useLoadingPlaceholder } from './useLoadingPlaceholder'

const { queryByText } = screen

describe('useLoadingPlaceholder', () => {
  it('renders "Loading..." text while loading', () => {
    const { LoadingPlaceholder } = useLoadingPlaceholder(true)
    render(<LoadingPlaceholder />)

    expect(queryByText('Loading...')).toBeInTheDocument()
  })

  it('defaults to false and renders an empty component', () => {
    const { LoadingPlaceholder } = useLoadingPlaceholder()
    render(<LoadingPlaceholder />)

    expect(document.documentElement).toHaveTextContent('')
  })

  it('exposes "isLoading" variable', () => {
    const { isLoading } = useLoadingPlaceholder(true)

    expect(isLoading).toBe(true)
  })
})
