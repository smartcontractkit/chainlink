import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import { render, screen } from '@testing-library/react'
import BaseLink from '../../src/components/BaseLink'

const { queryByText } = screen

const renderBaseLink = (link: React.ReactNode) =>
  render(<MemoryRouter>{link}</MemoryRouter>)

describe('components/BaseLink', () => {
  it('renders an anchor', () => {
    renderBaseLink(<BaseLink href="/foo">My Link</BaseLink>)

    expect(queryByText('My Link')).toBeInTheDocument()
    expect(queryByText('My Link')?.closest('a')).toHaveAttribute('href', '/foo')
  })

  it('can render an id', () => {
    renderBaseLink(
      <BaseLink id="my-id" href="/foo">
        My Link
      </BaseLink>,
    )

    expect(queryByText('My Link')?.closest('a')).toHaveAttribute('id', 'my-id')
  })

  it('can render a css class', () => {
    renderBaseLink(
      <BaseLink className="my-css-class" href="/foo">
        My Link
      </BaseLink>,
    )

    expect(queryByText('My Link')?.closest('a')).toHaveAttribute(
      'class',
      'my-css-class',
    )
  })
})
