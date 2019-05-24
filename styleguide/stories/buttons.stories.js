import React from 'react'
import { storiesOf } from '@storybook/react'
import { muiTheme } from 'storybook-addon-material-ui'
import { createMuiTheme } from '@material-ui/core/styles'
import Button from '@material-ui/core/Button'
import theme from '../theme'

const customTheme = createMuiTheme(theme)

storiesOf('Buttons', module)
  .addDecorator(muiTheme([customTheme]))
  .add('Contained', () => (
    <React.Fragment>
      <Button variant="contained">Default Button</Button>
      <Button variant="contained" color="primary">
        Primary Button
      </Button>
      <Button variant="contained" color="secondary">
        Secondary
      </Button>
      <Button variant="contained" color="secondary" disabled>
        Disabled
      </Button>
      <Button variant="contained" href="#contained-buttons">
        Link
      </Button>
      <input accept="image/*" id="contained-button-file" multiple type="file" />
      <label htmlFor="contained-button-file">
        <Button variant="contained" component="span">
          Upload
        </Button>
      </label>
    </React.Fragment>
  ))
  .add('Outlined', () => (
    <React.Fragment>
      <Button variant="outlined">Default Button</Button>
      <Button variant="outlined" color="primary">
        Primary Button
      </Button>
      <Button variant="outlined" color="secondary">
        Secondary
      </Button>
      <Button variant="outlined" color="secondary" disabled>
        Disabled
      </Button>
      <Button variant="outlined" href="#outlined-buttons">
        Link
      </Button>
      <input accept="image/*" id="outlined-button-file" multiple type="file" />
      <label htmlFor="outlined-button-file">
        <Button variant="outlined" component="span">
          Upload
        </Button>
      </label>
    </React.Fragment>
  ))
