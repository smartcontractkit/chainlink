import * as React from 'react'

import { render, screen, waitFor } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { FeedsManagerForm, FormValues } from './FeedsManagerForm'

const { getByRole, getByTestId, getByText } = screen

describe('FeedsManagerForm', () => {
  it('validates input', async () => {
    const handleSubmit = jest.fn()
    const initialValues: FormValues = {
      name: '',
      uri: '',
      publicKey: '',
      jobTypes: [],
      isBootstrapPeer: false,
      bootstrapPeerMultiaddr: undefined,
    }

    render(
      <FeedsManagerForm
        initialValues={initialValues}
        onSubmit={handleSubmit}
      />,
    )

    userEvent.click(
      getByRole('checkbox', {
        name: 'Is this node running as a bootstrap peer?',
      }),
    )

    userEvent.click(getByRole('button', { name: /submit/i }))

    await waitFor(() => {
      expect(getByTestId('name-helper-text')).toHaveTextContent('Required')
      expect(getByTestId('uri-helper-text')).toHaveTextContent('Required')
      expect(getByTestId('publicKey-helper-text')).toHaveTextContent('Required')
      expect(
        getByTestId('bootstrapPeerMultiaddr-helper-text'),
      ).toHaveTextContent('Required')
    })
  })

  it('submits the form', async () => {
    const handleSubmit = jest.fn()
    const initialValues: FormValues = {
      name: '',
      uri: '',
      publicKey: '',
      jobTypes: [],
      isBootstrapPeer: false,
      bootstrapPeerMultiaddr: undefined,
    }

    render(
      <FeedsManagerForm
        initialValues={initialValues}
        onSubmit={handleSubmit}
      />,
    )

    userEvent.type(
      getByRole('textbox', { name: 'Name *' }),
      'Chainlink Feeds Manager',
    )
    userEvent.type(getByRole('textbox', { name: 'URI *' }), 'localhost:8080')
    userEvent.type(getByRole('textbox', { name: 'Public Key *' }), '11111')
    userEvent.click(getByRole('checkbox', { name: 'Flux Monitor' }))

    userEvent.click(getByText(/submit/i))

    await waitFor(() =>
      expect(handleSubmit).toHaveBeenCalledWith(
        {
          name: 'Chainlink Feeds Manager',
          uri: 'localhost:8080',
          publicKey: '11111',
          jobTypes: ['FLUX_MONITOR'],
          isBootstrapPeer: false,
          bootstrapPeerMultiaddr: undefined,
        },
        expect.anything(),
      ),
    )
  })
})
