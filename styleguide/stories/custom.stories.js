import React from 'react'
import { storiesOf } from '@storybook/react'
import { muiTheme } from 'storybook-addon-material-ui'
import { createMuiTheme } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import JavascriptTimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
import TimeAgo from '../components/TimeAgo'
import Logo from '../components/Logo'
import theme from '../theme'

const customTheme = createMuiTheme(theme)

window.JavascriptTimeAgo = JavascriptTimeAgo
JavascriptTimeAgo.locale(en)

storiesOf('Custom', module)
  .addDecorator(muiTheme([customTheme]))
  .add('Logo', () => <Logo width={40} height={50} />)
  .add('TimeAgo', () => (
    <Typography>
      <TimeAgo>2018-11-27T02:26:42.014852Z</TimeAgo>
    </Typography>
  ))
