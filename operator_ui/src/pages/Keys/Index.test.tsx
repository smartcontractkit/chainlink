import React from 'react'
import { Route } from 'react-router-dom'
import { jsonApiOcrKeys } from 'factories/jsonApiOcrKeys'
import { syncFetch } from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'
import { mountWithProviders } from 'test-helpers/mountWithTheme'

import {
  ENDPOINT as OCR_ENDPOINT,
  INDEX_ENDPOINT as OCR_INDEX_ENDPOINT,
} from 'api/v2/ocrKeys'
import { KeysIndex } from 'pages/Keys/Index'

describe('pages/Keys/Index', () => {
  describe('Off-Chain Reporting keys', () => {
    it('renders the list of keys', async () => {
      const [expectedKey1, expectedKey2] = [
        {
          createdAt: new Date().toISOString(),
          OffChainPublicKey: 'OffChainPublicKey1',
          ConfigPublicKey: 'ConfigPublicKey1',
          OnChainSigningAddress: 'OnChainSigningAddress1',
        },
        {
          createdAt: new Date().toISOString(),
          OffChainPublicKey: 'OffChainPublicKey2',
          ConfigPublicKey: 'ConfigPublicKey2',
          OnChainSigningAddress: 'OnChainSigningAddress2',
        },
      ]

      global.fetch.getOnce(
        globPath(OCR_INDEX_ENDPOINT),
        jsonApiOcrKeys([expectedKey1, expectedKey2]),
      )

      const wrapper = mountWithProviders(<Route component={KeysIndex} />)
      await syncFetch(wrapper)
      expect(wrapper.text()).toContain('just now')
      expect(wrapper.text()).toContain('Delete')
      expect(wrapper.find('tbody').children().length).toEqual(2)
      expect(wrapper.text()).toContain(expectedKey1.OffChainPublicKey)
      expect(wrapper.text()).toContain(expectedKey1.ConfigPublicKey)
      expect(wrapper.text()).toContain(expectedKey1.OnChainSigningAddress)
      expect(wrapper.text()).toContain(expectedKey1.OffChainPublicKey)
      expect(wrapper.text()).toContain(expectedKey2.OffChainPublicKey)
      expect(wrapper.text()).toContain(expectedKey2.ConfigPublicKey)
      expect(wrapper.text()).toContain(expectedKey2.OnChainSigningAddress)
      expect(wrapper.text()).toContain(expectedKey2.OffChainPublicKey)
    })

    it('allows to create a new key bundle', async () => {
      const expectedKey = {
        OffChainPublicKey: 'OffChainPublicKey1',
      }
      global.fetch.getOnce(globPath(OCR_INDEX_ENDPOINT), [])
      const wrapper = mountWithProviders(<Route component={KeysIndex} />)
      await syncFetch(wrapper)

      expect(wrapper.find('tbody').children().length).toEqual(0)
      expect(wrapper.text()).not.toContain(expectedKey.OffChainPublicKey)

      global.fetch.getOnce(
        globPath(OCR_INDEX_ENDPOINT),
        jsonApiOcrKeys([expectedKey]),
      )
      global.fetch.postOnce(globPath(OCR_ENDPOINT), {})
      wrapper.find('[data-testid="keys-ocr-create"]').first().simulate('click')
      await syncFetch(wrapper)

      expect(wrapper.find('tbody').children().length).toEqual(1)
      expect(wrapper.text()).toContain(expectedKey.OffChainPublicKey)
    })

    it('allows to delete a key bundle', async () => {
      const expectedKey = {
        OffChainPublicKey: 'OffChainPublicKey1',
        id: 'keyId',
      }
      global.fetch.getOnce(
        globPath(OCR_INDEX_ENDPOINT),
        jsonApiOcrKeys([expectedKey]),
      )
      const wrapper = mountWithProviders(<Route component={KeysIndex} />)
      await syncFetch(wrapper)

      expect(wrapper.find('tbody').children().length).toEqual(1)
      expect(wrapper.text()).toContain(expectedKey.OffChainPublicKey)

      global.fetch.getOnce(globPath(OCR_INDEX_ENDPOINT), {})
      global.fetch.deleteOnce(
        globPath(`${OCR_INDEX_ENDPOINT}/${expectedKey.id}`),
        {},
      )
      wrapper
        .find('[data-testid="keys-ocr-delete-dialog"]')
        .first()
        .simulate('click')
      wrapper
        .find('[data-testid="keys-ocr-delete-confirm"]')
        .first()
        .simulate('click')
      await syncFetch(wrapper)

      expect(wrapper.find('tbody').children().length).toEqual(0)
      expect(wrapper.text()).not.toContain(expectedKey.OffChainPublicKey)
    })
  })
})
