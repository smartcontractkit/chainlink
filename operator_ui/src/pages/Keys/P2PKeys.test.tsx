import React from 'react'
import { jsonApiP2PKeys, P2PKeyBundle } from 'factories/jsonApiP2PKeys'

import globPath from 'test-helpers/globPath'
import { partialAsFull } from 'support/test-helpers/partialAsFull'

import {
  ENDPOINT as P2P_ENDPOINT,
  INDEX_ENDPOINT as P2P_INDEX_ENDPOINT,
} from 'api/v2/p2pKeys'

import { P2PKeys } from './P2PKeys'
import {
  renderWithRouter,
  screen,
  waitForElementToBeRemoved,
} from 'support/test-utils'
import userEvent from '@testing-library/user-event'

const { getByText, getByRole, getAllByRole, findAllByRole, findByText } = screen

describe('pages/Keys/P2PKeys', () => {
  describe('P2P keys', () => {
    it('renders the list of keys', async () => {
      const [expectedKey1, expectedKey2] = [
        partialAsFull<P2PKeyBundle>({
          publicKey: 'publicKey1',
          peerId: 'peerId1',
        }),
        partialAsFull<P2PKeyBundle>({
          publicKey: 'publicKey2',
          peerId: 'peerId2',
        }),
      ]

      global.fetch.getOnce(
        globPath(P2P_INDEX_ENDPOINT),
        jsonApiP2PKeys([expectedKey1, expectedKey2]),
      )

      renderWithRouter(<P2PKeys />)

      await waitForElementToBeRemoved(getByText('Loading...'))

      expect(getAllByRole('button', { name: 'Delete' })).toHaveLength(2)

      const rows = await findAllByRole('row')
      expect(rows).toHaveLength(3)

      expect(rows[1]).toHaveTextContent(expectedKey1.publicKey)
      expect(rows[1]).toHaveTextContent(expectedKey1.peerId)

      expect(rows[2]).toHaveTextContent(expectedKey2.publicKey)
      expect(rows[2]).toHaveTextContent(expectedKey2.peerId)
    })

    it('allows to create a new key bundle', async () => {
      const expectedKey = partialAsFull<P2PKeyBundle>({
        peerId: 'peerId',
      })
      global.fetch.getOnce(globPath(P2P_INDEX_ENDPOINT), jsonApiP2PKeys([]))

      renderWithRouter(<P2PKeys />)

      expect(await findByText('No entries to show.')).toBeInTheDocument()

      global.fetch.postOnce(globPath(P2P_ENDPOINT), {})
      global.fetch.getOnce(
        globPath(P2P_INDEX_ENDPOINT),
        jsonApiP2PKeys([expectedKey]),
      )

      userEvent.click(getByRole('button', { name: 'New P2P Key' }))

      await waitForElementToBeRemoved(getByText('No entries to show.'))

      const rows = await findAllByRole('row')
      expect(rows).toHaveLength(2)

      expect(rows[1]).toHaveTextContent(expectedKey.peerId)
    })

    it('allows to delete a key bundle', async () => {
      const expectedKey = partialAsFull<P2PKeyBundle>({
        peerId: 'peerId',
        id: 'keyId',
      })
      global.fetch.getOnce(
        globPath(P2P_INDEX_ENDPOINT),
        jsonApiP2PKeys([expectedKey]),
      )

      renderWithRouter(<P2PKeys />)

      await waitForElementToBeRemoved(getByText('Loading...'))

      let rows = await findAllByRole('row')
      expect(rows).toHaveLength(2)
      expect(rows[1]).toHaveTextContent(expectedKey.peerId as string)

      global.fetch.getOnce(globPath(P2P_INDEX_ENDPOINT), {})
      global.fetch.deleteOnce(
        globPath(`${P2P_INDEX_ENDPOINT}/${expectedKey.id}`),
        {},
      )

      userEvent.click(getByRole('button', { name: 'Delete' }))
      userEvent.click(getByRole('button', { name: 'Yes' }))

      rows = await findAllByRole('row')
      expect(rows).toHaveLength(1)
    })
  })
})
