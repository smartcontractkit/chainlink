import React from 'react'
import globPath from 'test-helpers/globPath'
import { logConfigFactory } from 'factories/jsonApiLogConfig'
import { LoggingCard } from './LoggingCard'
import {
  renderWithRouter,
  screen,
  waitForElementToBeRemoved,
} from 'support/test-utils'
import userEvent from '@testing-library/user-event'
import Notifications from 'pages/Notifications'

const { findByText, getByRole, getByText, getByTestId } = screen

describe('pages/Configuration/LoggingCard', () => {
  it('renders the logging configuration card', async () => {
    const logConfig = logConfigFactory({
      serviceName: ['Global', 'IsSqlEnabled'],
      logLevel: ['info', 'false'],
    })

    global.fetch.getOnce(globPath('/v2/log'), logConfig)

    renderWithRouter(<LoggingCard />)

    await waitForElementToBeRemoved(getByText('Loading...'))

    expect(getByTestId('logging-form')).toHaveFormValues({
      level: 'info',
      sqlEnabled: false,
    })
  })

  it('updates the logging configuration', async () => {
    const logConfig = logConfigFactory({
      serviceName: ['Global', 'IsSqlEnabled'],
      logLevel: ['debug', 'false'],
      defaultLogLevel: 'warn',
    })

    global.fetch.getOnce(globPath('/v2/log'), logConfig)
    global.fetch.patchOnce(globPath('/v2/log'), logConfig)

    renderWithRouter(
      <>
        <Notifications />
        <LoggingCard />
      </>,
    )

    await waitForElementToBeRemoved(getByText('Loading...'))

    userEvent.click(
      getByRole('checkbox', { name: /log sql statements \(debug only\)/i }),
    )

    userEvent.click(getByRole('button', { name: 'Update' }))

    expect(
      await findByText('Logging Configuration Updated'),
    ).toBeInTheDocument()
  })
})
