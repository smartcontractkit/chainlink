import { createMuiTheme } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import { storiesOf } from '@storybook/react'
import JavascriptTimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
import React from 'react'
import { muiTheme } from 'storybook-addon-material-ui'
import { Logo, theme, TimeAgo } from '../src'

declare global {
  interface Window {
    JavascriptTimeAgo: any
  }
}

const customTheme = createMuiTheme(theme)

window.JavascriptTimeAgo = JavascriptTimeAgo
JavascriptTimeAgo.locale(en)

storiesOf('Custom', module)
  .addDecorator(muiTheme([customTheme]))
  .add('Logo', () => <Logo src="" width={40} height={50} />)
  .add('TimeAgo', () => (
    <Typography>
      <TimeAgo tooltip={false}>2018-11-27T02:26:42.014852Z</TimeAgo>
    </Typography>
  ))
