import React from 'react'
import CreateLayout from 'images/create.js'
import { withStyles } from '@material-ui/core/styles'

const styles = () => {
  return {
    logo: {
      top: '30%',
      left: '50%',
      transform: 'translate(-50%, -30%)',
      position: 'absolute'
    }
  }
}

const Create = ({ classes }) => (
  <div className={classes.logo}>
    <CreateLayout />
  </div>
)

export default withStyles(styles)(Create)
