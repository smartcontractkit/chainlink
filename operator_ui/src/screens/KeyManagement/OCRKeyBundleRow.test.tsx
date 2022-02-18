import * as React from 'react'

import { render, screen } from 'support/test-utils'

import { buildOCRKeyBundle } from 'support/factories/gql/fetchOCRKeyBundles'
import { OCRKeyBundleRow } from './OCRKeyBundleRow'
import userEvent from '@testing-library/user-event'

const { getByRole, queryByText } = screen

describe('OCRKeyBundleRow', () => {
  let handleDelete: jest.Mock

  beforeEach(() => {
    handleDelete = jest.fn()
  })

  function renderComponent(bundle: OcrKeyBundlesPayload_ResultsFields) {
    render(
      <table>
        <tbody>
          <OCRKeyBundleRow bundle={bundle} onDelete={handleDelete} />
        </tbody>
      </table>,
    )
  }

  it('renders a row', () => {
    const bundle = buildOCRKeyBundle()

    renderComponent(bundle)

    expect(queryByText(`Key ID: ${bundle.id}`)).toBeInTheDocument()
    expect(
      queryByText(`Config Public Key: ${bundle.configPublicKey}`),
    ).toBeInTheDocument()
    expect(
      queryByText(`Signing Address: ${bundle.onChainSigningAddress}`),
    ).toBeInTheDocument()
    expect(
      queryByText(`Off-Chain Public Key: ${bundle.offChainPublicKey}`),
    ).toBeInTheDocument()
  })

  it('calls delete', () => {
    const bundle = buildOCRKeyBundle()

    renderComponent(bundle)

    userEvent.click(getByRole('button', { name: /delete/i }))

    expect(handleDelete).toHaveBeenCalled()
  })
})
