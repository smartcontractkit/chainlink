import React from 'react'
import { jsonApiP2PKeys, P2PKeyBundle } from 'factories/jsonApiP2PKeys'

import { syncFetch } from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { partialAsFull } from 'support/test-helpers/partialAsFull'

import {
  ENDPOINT as P2P_ENDPOINT,
  INDEX_ENDPOINT as P2P_INDEX_ENDPOINT,
} from 'api/v2/p2pKeys'

import { P2PKeys } from './P2PKeys'

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

      const wrapper = mountWithProviders(<P2PKeys />)
      await syncFetch(wrapper)
      expect(wrapper.text()).toContain('Delete')
      expect(wrapper.find('tbody').children().length).toEqual(2)
      expect(wrapper.text()).toContain(expectedKey1.publicKey)
      expect(wrapper.text()).toContain(expectedKey1.peerId)
      expect(wrapper.text()).toContain(expectedKey2.publicKey)
      expect(wrapper.text()).toContain(expectedKey2.peerId)
    })

    it('allows to create a new key bundle', async () => {
      const expectedKey = partialAsFull<P2PKeyBundle>({
        peerId: 'peerId',
      })
      global.fetch.getOnce(globPath(P2P_INDEX_ENDPOINT), [])
      const wrapper = mountWithProviders(<P2PKeys />)
      await syncFetch(wrapper)

      expect(wrapper.find('tbody').children().length).toEqual(0)
      expect(wrapper.text()).not.toContain(expectedKey.peerId)

      global.fetch.getOnce(
        globPath(P2P_INDEX_ENDPOINT),
        jsonApiP2PKeys([expectedKey]),
      )
      global.fetch.postOnce(globPath(P2P_ENDPOINT), {})
      wrapper.find('[data-testid="keys-create"]').first().simulate('click')
      await syncFetch(wrapper)

      expect(wrapper.find('tbody').children().length).toEqual(1)
      expect(wrapper.text()).toContain(expectedKey.peerId)
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
      const wrapper = mountWithProviders(<P2PKeys />)
      await syncFetch(wrapper)

      expect(wrapper.find('tbody').children().length).toEqual(1)
      expect(wrapper.text()).toContain(expectedKey.peerId)

      global.fetch.getOnce(globPath(P2P_INDEX_ENDPOINT), {})
      global.fetch.deleteOnce(
        globPath(`${P2P_INDEX_ENDPOINT}/${expectedKey.id}`),
        {},
      )
      wrapper
        .find('[data-testid="keys-delete-dialog"]')
        .first()
        .simulate('click')
      wrapper
        .find('[data-testid="keys-delete-confirm"]')
        .first()
        .simulate('click')
      await syncFetch(wrapper)

      expect(wrapper.find('tbody').children().length).toEqual(0)
      expect(wrapper.text()).not.toContain(expectedKey.peerId)
    })
  })
})
