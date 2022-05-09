import * as React from 'react'

import { GraphQLError } from 'graphql'
import {
  renderWithRouter,
  screen,
  waitForElementToBeRemoved,
} from 'support/test-utils'
import { MockedProvider, MockedResponse } from '@apollo/client/testing'
import userEvent from '@testing-library/user-event'

import {
  OCRKeys,
  CREATE_OCR_KEY_BUNDLE_MUTATION,
  DELETE_OCR_KEY_BUNDLE_MUTATION,
} from './OCRKeys'
import {
  buildOCRKeyBundle,
  buildOCRKeyBundles,
} from 'support/factories/gql/fetchOCRKeyBundles'
import Notifications from 'pages/Notifications'
import { OCR_KEY_BUNDLES_QUERY } from 'src/hooks/queries/useOCRKeysQuery'
import { waitForLoading } from 'support/test-helpers/wait'

const { findByText, getByRole, queryByText } = screen

function renderComponent(mocks: MockedResponse[]) {
  renderWithRouter(
    <>
      <Notifications />
      <MockedProvider mocks={mocks} addTypename={false}>
        <OCRKeys />
      </MockedProvider>
    </>,
  )
}

function fetchOCRKeyBundlesQuery(
  bundles: ReadonlyArray<OcrKeyBundlesPayload_ResultsFields>,
) {
  return {
    request: {
      query: OCR_KEY_BUNDLES_QUERY,
    },
    result: {
      data: {
        ocrKeyBundles: {
          results: bundles,
        },
      },
    },
  }
}

describe('OCRKeys', () => {
  it('renders the page', async () => {
    const payload = buildOCRKeyBundles()
    const mocks: MockedResponse[] = [fetchOCRKeyBundlesQuery(payload)]

    renderComponent(mocks)

    await waitForLoading()

    expect(await findByText(`Key ID: ${payload[0].id}`)).toBeInTheDocument()
  })

  it('renders GQL query errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: OCR_KEY_BUNDLES_QUERY,
        },
        result: {
          errors: [new GraphQLError('Error!')],
        },
      },
    ]

    renderComponent(mocks)

    expect(await findByText('Error!')).toBeInTheDocument()
  })

  it('creates an OCR Key Bundle', async () => {
    const payload = buildOCRKeyBundle()

    const mocks: MockedResponse[] = [
      fetchOCRKeyBundlesQuery([]),
      {
        request: {
          query: CREATE_OCR_KEY_BUNDLE_MUTATION,
        },
        result: {
          data: {
            createOCRKeyBundle: {
              __typename: 'CreateOCRKeyBundleSuccess',
              bundle: payload,
            },
          },
        },
      },
      fetchOCRKeyBundlesQuery([payload]),
    ]

    renderComponent(mocks)

    await waitForLoading()

    userEvent.click(getByRole('button', { name: /new ocr key/i }))

    expect(
      await findByText(
        `Successfully created Off-ChainReporting Key Bundle: ${payload.id}`,
      ),
    ).toBeInTheDocument()
    expect(await findByText(`Key ID: ${payload.id}`)).toBeInTheDocument()
  })

  it('errors on create', async () => {
    const mocks: MockedResponse[] = [
      fetchOCRKeyBundlesQuery([]),
      {
        request: {
          query: CREATE_OCR_KEY_BUNDLE_MUTATION,
        },
        result: {
          errors: [new GraphQLError('Error!')],
        },
      },
    ]

    renderComponent(mocks)

    await waitForLoading()

    userEvent.click(getByRole('button', { name: /new ocr key/i }))

    expect(await findByText('Error!')).toBeInTheDocument()
  })

  it('deletes an OCR Key Bundle', async () => {
    const payload = buildOCRKeyBundle()

    const mocks: MockedResponse[] = [
      fetchOCRKeyBundlesQuery([payload]),
      {
        request: {
          query: DELETE_OCR_KEY_BUNDLE_MUTATION,
          variables: { id: payload.id },
        },
        result: {
          data: {
            deleteOCRKeyBundle: {
              __typename: 'DeleteOCRKeyBundleSuccess',
              bundle: payload,
            },
          },
        },
      },
      fetchOCRKeyBundlesQuery([]),
    ]

    renderComponent(mocks)

    expect(await findByText(`Key ID: ${payload.id}`)).toBeInTheDocument()

    userEvent.click(getByRole('button', { name: /delete/i }))
    userEvent.click(getByRole('button', { name: /confirm/i }))

    await waitForElementToBeRemoved(getByRole('dialog'))

    expect(
      await findByText(
        'Successfully deleted Off-ChainReporting Key Bundle Key',
      ),
    ).toBeInTheDocument()

    expect(queryByText(`Key ID: ${payload.id}`)).toBeNull()
  })
})
