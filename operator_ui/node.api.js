import React from 'react'
import { ServerStyleSheets } from '@material-ui/styles'

import theme from './src/theme'
import { MuiThemeProvider } from '@material-ui/core/styles';

export default () => ({
  beforeRenderToHtml: (App, { meta }) => {
    meta.muiSheets = new ServerStyleSheets()
    return meta.muiSheets.collect(
      <MuiThemeProvider theme={theme}>{App}</MuiThemeProvider>
    )
  },
  headElements: (elements, { meta }) => [
    ...elements,
    meta.muiSheets.getStyleElement()
  ]
})
