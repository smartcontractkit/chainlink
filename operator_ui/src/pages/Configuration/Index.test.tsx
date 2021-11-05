import React from 'react'
import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import globPath from 'test-helpers/globPath'
import configurationFactory from 'factories/configuration'

import Configuration from 'pages/Configuration/Index'

const { findByText, queryByText } = screen

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

    expect(await findByText('Configuration')).toBeInTheDocument()

    expect(queryByText('BAND')).toBeInTheDocument()
    expect(queryByText('Major Lazer')).toBeInTheDocument()
    expect(queryByText('SINGER')).toBeInTheDocument()
    expect(queryByText('Bob Marley')).toBeInTheDocument()
  })
})
