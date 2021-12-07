import * as React from 'react'

import { render, screen, waitForElementToBeRemoved } from 'support/test-utils'

import {
  buildOCRKeyBundle,
  buildOCRKeyBundles,
} from 'support/factories/gql/fetchOCRKeyBundles'
import { OCRKeysCard, Props as OCRKeysCardProps } from './OCRKeysCard'
import userEvent from '@testing-library/user-event'

const { getAllByRole, getByRole, queryByRole, queryByText } = screen

function renderComponent(cardProps: OCRKeysCardProps) {
  render(<OCRKeysCard {...cardProps} />)
}

describe('CSAKeysCard', () => {
  let promise: Promise<any>
  let handleCreate: jest.Mock
  let handleDelete: jest.Mock

  beforeEach(() => {
    promise = Promise.resolve()
    handleCreate = jest.fn()
    handleDelete = jest.fn(() => promise)
  })

  it('renders the key bundles', () => {
    const bundles = buildOCRKeyBundles()

    renderComponent({
      loading: false,
      data: {
        ocrKeyBundles: {
          results: bundles,
        },
      },
      onCreate: handleCreate,
      onDelete: handleDelete,
    })

    expect(getAllByRole('row')).toHaveLength(3)

    expect(queryByText(`Key ID: ${bundles[0].id}`)).toBeInTheDocument()
    expect(queryByText(`Key ID: ${bundles[1].id}`)).toBeInTheDocument()
  })

  it('renders no content', () => {
    renderComponent({
      loading: false,
      data: {
        ocrKeyBundles: {
          results: [],
        },
      },
      onCreate: handleCreate,
      onDelete: handleDelete,
    })

    expect(queryByText('No entries to show')).toBeInTheDocument()
  })

  it('renders a loading spinner', () => {
    renderComponent({
      loading: true,
      onCreate: handleCreate,
      onDelete: handleDelete,
    })

    expect(queryByRole('progressbar')).toBeInTheDocument()
  })

  it('renders an error message', () => {
    renderComponent({
      loading: false,
      errorMsg: 'error message',
      onCreate: handleCreate,
      onDelete: handleDelete,
    })

    expect(queryByText('error message')).toBeInTheDocument()
  })

  it('calls onCreate', () => {
    renderComponent({
      loading: false,
      data: {
        ocrKeyBundles: {
          results: [],
        },
      },
      onCreate: handleCreate,
      onDelete: handleDelete,
    })

    userEvent.click(getByRole('button', { name: /new ocr key/i }))

    expect(handleCreate).toHaveBeenCalled()
  })

  it('calls onDelete', async () => {
    const bundle = buildOCRKeyBundle()
    renderComponent({
      loading: false,
      data: {
        ocrKeyBundles: {
          results: [bundle],
        },
      },
      onCreate: handleCreate,
      onDelete: handleDelete,
    })

    userEvent.click(getByRole('button', { name: /delete/i }))
    expect(queryByText(bundle.id)).toBeInTheDocument()

    userEvent.click(getByRole('button', { name: /confirm/i }))

    await waitForElementToBeRemoved(getByRole('dialog'))

    expect(handleDelete).toHaveBeenCalled()
  })
})
