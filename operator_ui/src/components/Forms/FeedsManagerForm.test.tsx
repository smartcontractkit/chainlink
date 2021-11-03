import '@testing-library/jest-dom'

import * as React from 'react'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'

import { FeedsManagerForm, FormValues } from './FeedsManagerForm'

test('submits the form', async () => {
  const { getByRole, getByTestId } = screen
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
    <FeedsManagerForm initialValues={initialValues} onSubmit={handleSubmit} />,
  )

  userEvent.type(
    getByRole('textbox', { name: 'Name *' }),
    'Chainlink Feeds Manager',
  )
  userEvent.type(getByRole('textbox', { name: 'URI *' }), 'localhost:8080')
  userEvent.type(getByRole('textbox', { name: 'Public Key *' }), '11111')
  userEvent.click(getByRole('checkbox', { name: 'Flux Monitor' }))

  userEvent.click(getByTestId('create-submit'))

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

test('validates input', async () => {
  const { getByRole, getByTestId } = screen
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
    <FeedsManagerForm initialValues={initialValues} onSubmit={handleSubmit} />,
  )

  userEvent.click(
    getByRole('checkbox', {
      name: 'Is this node running as a bootstrap peer?',
    }),
  )

  userEvent.click(getByTestId('create-submit'))

  await waitFor(() => {
    expect(getByTestId('name-helper-text')).toHaveTextContent('Required')
    expect(getByTestId('uri-helper-text')).toHaveTextContent('Required')
    expect(getByTestId('publicKey-helper-text')).toHaveTextContent('Required')
    expect(getByTestId('bootstrapPeerMultiaddr-helper-text')).toHaveTextContent(
      'Required',
    )
  })
})
