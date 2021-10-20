import React from 'react'
import { jsonApiCSAKeys, CSAKey } from 'factories/jsonApiCSAKeys'

import { syncFetch } from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { partialAsFull } from 'support/test-helpers/partialAsFull'

import { ENDPOINT as CSA_INDEX_ENDPOINT } from 'api/v2/csaKeys'

import { CSAKeys } from './CSAKeys'

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

      const wrapper = mountWithProviders(<CSAKeys />)
      await syncFetch(wrapper)
      expect(wrapper.text()).toContain('just now')
      expect(wrapper.find('tbody').children().length).toEqual(1)
      expect(wrapper.text()).toContain(expectedKey.publicKey)
    })

    it('allows to create a new CSA key', async () => {
      const expectedKey = partialAsFull<CSAKey>({
        publicKey: 'publicKey1',
      })
      global.fetch.getOnce(globPath(CSA_INDEX_ENDPOINT), jsonApiCSAKeys([]))
      const wrapper = mountWithProviders(<CSAKeys />)
      await syncFetch(wrapper)

      expect(wrapper.find('tbody').children().length).toEqual(1)
      expect(wrapper.text()).not.toContain(expectedKey.publicKey)
      expect(wrapper.text()).toContain('No entries to show.')

      // The button will only appear when there is no existing CSA Key
      global.fetch.postOnce(globPath(CSA_INDEX_ENDPOINT), {})
      // Refetch the keys
      global.fetch.getOnce(
        globPath(CSA_INDEX_ENDPOINT),
        jsonApiCSAKeys([expectedKey]),
      )

      wrapper.find('[data-testid="keys-create"]').first().simulate('click')
      await syncFetch(wrapper)

      expect(wrapper.find('tbody').children().length).toEqual(1)
      expect(wrapper.text()).toContain(expectedKey.publicKey)
    })
  })
})
