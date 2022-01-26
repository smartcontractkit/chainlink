import React from 'react'
import { jsonApiCSAKeys, CSAKey } from 'factories/jsonApiCSAKeys'

import globPath from 'test-helpers/globPath'
import { partialAsFull } from 'support/test-helpers/partialAsFull'

import { ENDPOINT as CSA_INDEX_ENDPOINT } from 'api/v2/csaKeys'

import { CSAKeys } from './CSAKeys'
import {
  renderWithRouter,
  screen,
  waitForElementToBeRemoved,
} from 'support/test-utils'
import userEvent from '@testing-library/user-event'

const { findAllByRole, getByRole, getByText, queryByText } = screen

describe('pages/Keys/CSAKeys', () => {
  describe('CSA keys', () => {
    it('renders the CSA key', async () => {
      const expectedKey = partialAsFull<CSAKey>({
        publicKey: 'publicKey1',
      })

      global.fetch.getOnce(
        globPath(CSA_INDEX_ENDPOINT),
        jsonApiCSAKeys([expectedKey]),
      )

      renderWithRouter(<CSAKeys />)

      await waitForElementToBeRemoved(getByText('Loading...'))

      const rows = await findAllByRole('row')
      expect(rows).toHaveLength(2)

      expect(rows[1]).toHaveTextContent(expectedKey.publicKey)
    })

    it('allows to create a new CSA key', async () => {
      const expectedKey = partialAsFull<CSAKey>({
        publicKey: 'publicKey1',
      })
      global.fetch.getOnce(globPath(CSA_INDEX_ENDPOINT), jsonApiCSAKeys([]))
      renderWithRouter(<CSAKeys />)

      await waitForElementToBeRemoved(getByText('Loading...'))

      expect(queryByText('No entries to show.')).toBeInTheDocument()

      // The button will only appear when there is no existing CSA Key
      global.fetch.postOnce(globPath(CSA_INDEX_ENDPOINT), {})
      // Refetch the keys
      global.fetch.getOnce(
        globPath(CSA_INDEX_ENDPOINT),
        jsonApiCSAKeys([expectedKey]),
      )

      userEvent.click(getByRole('button', { name: 'New CSA Key' }))

      await waitForElementToBeRemoved(getByText('No entries to show.'))

      const rows = await findAllByRole('row')
      expect(rows).toHaveLength(2)

      expect(rows[1]).toHaveTextContent(expectedKey.publicKey)
    })
  })
})
