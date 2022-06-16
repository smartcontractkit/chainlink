import * as React from 'react'

import { render, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { ChainConfigurationForm, FormValues } from './ChainConfigurationForm'

const { getByRole, findByTestId } = screen

describe('ChainConfigurationForm', () => {
  it('validates top level input', async () => {
    const handleSubmit = jest.fn()
    const initialValues: FormValues = {
      chainID: '',
      chainType: '',
      accountAddr: '',
      adminAddr: '',
      fluxMonitorEnabled: false,
      ocr1Enabled: false,
      ocr1IsBootstrap: false,
      ocr1Multiaddr: '',
      ocr1P2PPeerID: '',
      ocr1KeyBundleID: '',
      ocr2Enabled: false,
      ocr2IsBootstrap: false,
      ocr2Multiaddr: '',
      ocr2P2PPeerID: '',
      ocr2KeyBundleID: '',
    }

    render(
      <ChainConfigurationForm
        initialValues={initialValues}
        onSubmit={handleSubmit}
        accounts={[]}
        chainIDs={[]}
        p2pKeys={[]}
        ocrKeys={[]}
        ocr2Keys={[]}
        showSubmit
      />,
    )

    userEvent.click(getByRole('button', { name: /submit/i }))

    expect(await findByTestId('chainID-helper-text')).toHaveTextContent(
      'Required',
    )
    expect(await findByTestId('accountAddr-helper-text')).toHaveTextContent(
      'Required',
    )
    expect(await findByTestId('adminAddr-helper-text')).toHaveTextContent(
      'Required',
    )
  })

  it('validates OCR input', async () => {
    const handleSubmit = jest.fn()
    const initialValues: FormValues = {
      chainID: '',
      chainType: '',
      accountAddr: '',
      adminAddr: '',
      fluxMonitorEnabled: false,
      ocr1Enabled: false,
      ocr1IsBootstrap: false,
      ocr1Multiaddr: '',
      ocr1P2PPeerID: '',
      ocr1KeyBundleID: '',
      ocr2Enabled: false,
      ocr2IsBootstrap: false,
      ocr2Multiaddr: '',
      ocr2P2PPeerID: '',
      ocr2KeyBundleID: '',
    }

    render(
      <ChainConfigurationForm
        initialValues={initialValues}
        onSubmit={handleSubmit}
        accounts={[]}
        chainIDs={[]}
        p2pKeys={[]}
        ocrKeys={[]}
        ocr2Keys={[]}
        showSubmit
      />,
    )

    userEvent.click(getByRole('checkbox', { name: 'OCR' }))
    userEvent.click(getByRole('button', { name: /submit/i }))

    expect(await findByTestId('ocr1P2PPeerID-helper-text')).toHaveTextContent(
      'Required',
    )
    expect(await findByTestId('ocr1KeyBundleID-helper-text')).toHaveTextContent(
      'Required',
    )

    userEvent.click(
      getByRole('checkbox', {
        name: 'Is this node running as a bootstrap peer?',
      }),
    )

    expect(await findByTestId('ocr1Multiaddr-helper-text')).toHaveTextContent(
      'Required',
    )
  })

  it('validates OCR2 input', async () => {
    const handleSubmit = jest.fn()
    const initialValues: FormValues = {
      chainID: '',
      chainType: '',
      accountAddr: '',
      adminAddr: '',
      fluxMonitorEnabled: false,
      ocr1Enabled: false,
      ocr1IsBootstrap: false,
      ocr1Multiaddr: '',
      ocr1P2PPeerID: '',
      ocr1KeyBundleID: '',
      ocr2Enabled: false,
      ocr2IsBootstrap: false,
      ocr2Multiaddr: '',
      ocr2P2PPeerID: '',
      ocr2KeyBundleID: '',
    }

    render(
      <ChainConfigurationForm
        initialValues={initialValues}
        onSubmit={handleSubmit}
        accounts={[]}
        chainIDs={[]}
        p2pKeys={[]}
        ocrKeys={[]}
        ocr2Keys={[]}
        showSubmit
      />,
    )

    userEvent.click(getByRole('checkbox', { name: 'OCR' }))
    userEvent.click(getByRole('button', { name: /submit/i }))

    expect(await findByTestId('ocr2P2PPeerID-helper-text')).toHaveTextContent(
      'Required',
    )
    expect(await findByTestId('ocr2KeyBundleID-helper-text')).toHaveTextContent(
      'Required',
    )

    userEvent.click(
      getByRole('checkbox', {
        name: 'Is this node running as a bootstrap peer?',
      }),
    )

    expect(await findByTestId('ocr2Multiaddr-helper-text')).toHaveTextContent(
      'Required',
    )
  })
})
