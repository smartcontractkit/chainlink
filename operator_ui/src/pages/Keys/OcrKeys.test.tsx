import React from 'react'
import { jsonApiOcrKeys, OcrKeyBundle } from 'factories/jsonApiOcrKeys'
import globPath from 'test-helpers/globPath'
import { partialAsFull } from 'support/test-helpers/partialAsFull'
import {
  ENDPOINT as OCR_ENDPOINT,
  INDEX_ENDPOINT as OCR_INDEX_ENDPOINT,
} from 'api/v2/ocrKeys'
import { OcrKeys } from './OcrKeys'
import {
  renderWithRouter,
  screen,
  waitForElementToBeRemoved,
} from 'support/test-utils'
import userEvent from '@testing-library/user-event'

const { findAllByRole, getByText, getAllByRole, getByRole, findByText } = screen

describe('pages/Keys/OcrKeys', () => {
  describe('Off-Chain Reporting keys', () => {
    it('renders the list of keys', async () => {
      const [expectedKey1, expectedKey2] = [
        partialAsFull<OcrKeyBundle>({
          id: 'keyId1',
          offChainPublicKey: 'offChainPublicKey1',
          configPublicKey: 'configPublicKey1',
          onChainSigningAddress: 'onChainSigningAddress1',
        }),
        partialAsFull<OcrKeyBundle>({
          id: 'keyId2',
          offChainPublicKey: 'offChainPublicKey2',
          configPublicKey: 'configPublicKey2',
          onChainSigningAddress: 'onChainSigningAddress2',
        }),
      ]

      global.fetch.getOnce(
        globPath(OCR_INDEX_ENDPOINT),
        jsonApiOcrKeys([expectedKey1, expectedKey2]),
      )

      renderWithRouter(<OcrKeys />)

      await waitForElementToBeRemoved(getByText('Loading...'))

      expect(getAllByRole('button', { name: 'Delete' })).toHaveLength(2)

      const rows = await findAllByRole('row')
      expect(rows).toHaveLength(3)

      expect(rows[1]).toHaveTextContent(expectedKey1.id as string)
      expect(rows[1]).toHaveTextContent(expectedKey1.offChainPublicKey)
      expect(rows[1]).toHaveTextContent(expectedKey1.configPublicKey as string)
      expect(rows[1]).toHaveTextContent(
        expectedKey1.onChainSigningAddress as string,
      )
      expect(rows[1]).toHaveTextContent(
        expectedKey1.offChainPublicKey as string,
      )

      expect(rows[2]).toHaveTextContent(expectedKey2.id as string)
      expect(rows[2]).toHaveTextContent(expectedKey2.offChainPublicKey)
      expect(rows[2]).toHaveTextContent(expectedKey2.configPublicKey as string)
      expect(rows[2]).toHaveTextContent(
        expectedKey2.onChainSigningAddress as string,
      )
      expect(rows[2]).toHaveTextContent(
        expectedKey2.offChainPublicKey as string,
      )
    })

    it('allows to create a new key bundle', async () => {
      const expectedKey = partialAsFull<OcrKeyBundle>({
        id: 'keyId',
      })
      global.fetch.getOnce(globPath(OCR_INDEX_ENDPOINT), jsonApiOcrKeys([]))

      renderWithRouter(<OcrKeys />)

      expect(await findByText('No entries to show.')).toBeInTheDocument()

      global.fetch.postOnce(globPath(OCR_ENDPOINT), {})

      global.fetch.getOnce(
        globPath(OCR_INDEX_ENDPOINT),
        jsonApiOcrKeys([expectedKey]),
      )

      userEvent.click(getByRole('button', { name: 'New OCR Key' }))

      await waitForElementToBeRemoved(getByText('No entries to show.'))

      const rows = await findAllByRole('row')
      expect(rows).toHaveLength(2)

      expect(rows[1]).toHaveTextContent(expectedKey.id as string)
    })

    it('allows to delete a key bundle', async () => {
      const expectedKey = partialAsFull<OcrKeyBundle>({
        id: 'keyId',
      })
      global.fetch.getOnce(
        globPath(OCR_INDEX_ENDPOINT),
        jsonApiOcrKeys([expectedKey]),
      )

      renderWithRouter(<OcrKeys />)

      await waitForElementToBeRemoved(getByText('Loading...'))

      let rows = await findAllByRole('row')
      expect(rows).toHaveLength(2)
      expect(rows[1]).toHaveTextContent(expectedKey.id as string)

      global.fetch.getOnce(globPath(OCR_INDEX_ENDPOINT), {})
      global.fetch.deleteOnce(
        globPath(`${OCR_INDEX_ENDPOINT}/${expectedKey.id}`),
        {},
      )

      userEvent.click(getByRole('button', { name: 'Delete' }))
      userEvent.click(getByRole('button', { name: 'Yes' }))

      rows = await findAllByRole('row')
      expect(rows).toHaveLength(1)
    })
  })
})
