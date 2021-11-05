import userEvent from '@testing-library/user-event'
import React from 'react'
import { render, screen } from 'support/test-utils'
import { Delete } from './Delete'

const { getByRole, getByText } = screen

describe('pages/Keys/Delete', () => {
  it('open modal and confirm delete', async () => {
    const expectedKeyId = 'KeyId'
    const expectedKeyValue = 'keyValue'
    const expectedOnDelete = jest.fn()

    render(
      <Delete
        onDelete={expectedOnDelete}
        keyId={expectedKeyId}
        keyValue={expectedKeyValue}
      />,
    )

    userEvent.click(getByRole('button', { name: 'Delete' }))

    expect(getByText(expectedKeyValue)).toBeInTheDocument()

    userEvent.click(getByRole('button', { name: 'Yes' }))

    expect(expectedOnDelete).toBeCalledWith(expectedKeyId)
  })
})
