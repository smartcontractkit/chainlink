import React from 'react'

import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'

import { SearchTextField } from './SearchTextField'

const { getByRole } = screen

describe('SearchTextField', () => {
  it('renders the text field', () => {
    const handleChange = jest.fn()

    render(
      <SearchTextField
        value={''}
        onChange={handleChange}
        placeholder="Search..."
      />,
    )

    const textbox = getByRole('textbox')
    expect(textbox).toHaveAttribute('placeholder', 'Search...')
    expect(textbox).toHaveAttribute('value', '')

    userEvent.paste(textbox, 'foo')

    expect(handleChange).toHaveBeenCalledWith('foo')
  })
})
