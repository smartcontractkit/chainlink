import React from 'react'
import { Route } from 'react-router-dom'
import { renderWithRouter, screen, waitFor } from 'support/test-utils'
import globPath from 'test-helpers/globPath'
import configurationFactory from 'factories/configuration'

import Configuration from 'pages/Configuration/Index'

const { getByText } = screen

describe('pages/Configuration', () => {
  it('renders the list of configuration options', async () => {
    const configurationResponse = configurationFactory({
      BAND: 'Major Lazer',
      SINGER: 'Bob Marley',
    })
    global.fetch.getOnce(globPath('/v2/config'), configurationResponse)

    renderWithRouter(<Route path="/config" component={Configuration} />, {
      initialEntries: ['/config'],
    })

    await waitFor(() => getByText('Configuration'))

    expect(getByText('BAND')).toBeInTheDocument()
    expect(getByText('Major Lazer')).toBeInTheDocument()
    expect(getByText('SINGER')).toBeInTheDocument()
    expect(getByText('Bob Marley')).toBeInTheDocument()
  })
})
