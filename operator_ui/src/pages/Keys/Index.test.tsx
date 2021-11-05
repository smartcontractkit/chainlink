import React from 'react'
import { jsonApiOcrKeys, OcrKeyBundle } from 'factories/jsonApiOcrKeys'
import { jsonApiP2PKeys, P2PKeyBundle } from 'factories/jsonApiP2PKeys'
import { jsonApiFeatureFlags } from 'factories/jsonApiFeatures'
import { accountBalances } from 'factories/accountBalance'

import globPath from 'test-helpers/globPath'
import { partialAsFull } from 'support/test-helpers/partialAsFull'
import { ACCOUNT_BALANCES_ENDPOINT } from 'api/v2/user/balances'
import { INDEX_ENDPOINT as FEATURES_INDEX_ENDPOINT } from 'api/v2/features'
import { INDEX_ENDPOINT as OCR_INDEX_ENDPOINT } from 'api/v2/ocrKeys'
import { INDEX_ENDPOINT as P2P_INDEX_ENDPOINT } from 'api/v2/p2pKeys'

import { KeysIndex } from 'pages/Keys/Index'
import {
  renderWithRouter,
  screen,
  waitForElementToBeRemoved,
} from 'support/test-utils'

const { getAllByText, getAllByRole, getByText } = screen

describe('pages/Keys/Index', () => {
  it('renders the OCR and P2P lists of keys', async () => {
    const [expectedOcr, expectedP2P] = [
      partialAsFull<OcrKeyBundle>({
        id: 'OcrKey',
      }),
      partialAsFull<P2PKeyBundle>({
        peerId: 'P2PId',
      }),
    ]

    global.fetch.getOnce(
      globPath(FEATURES_INDEX_ENDPOINT),
      jsonApiFeatureFlags(),
    )
    global.fetch.getOnce(
      globPath(OCR_INDEX_ENDPOINT),
      jsonApiOcrKeys([expectedOcr]),
    )
    global.fetch.getOnce(
      globPath(P2P_INDEX_ENDPOINT),
      jsonApiP2PKeys([expectedP2P]),
    )
    global.fetch.getOnce(
      globPath(ACCOUNT_BALANCES_ENDPOINT),
      accountBalances([]),
    )

    renderWithRouter(<KeysIndex />)

    await waitForElementToBeRemoved(() => getAllByText('Loading...'))

    expect(getAllByRole('row')).toHaveLength(6)

    expect(getByText(/peer id: p2pid/i)).toBeInTheDocument()
    expect(getByText(/key id: ocrkey/i)).toBeInTheDocument()
  })
})
