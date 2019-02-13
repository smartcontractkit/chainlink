import React from 'react'
import NotFoundSVG from 'images/four-oh-four.js'
import { withStyles } from '@material-ui/core/styles'

const styles = () => ({
  logo: {
    top: '30%',
    left: '50%',
    transform: 'translate(-50%, -30%)',
    position: 'absolute'
  }
})

const Logo = ({ classes }) => (
  <div className={classes.logo}>
    <NotFoundSVG />
  </div>
)

export default withStyles(styles)(Logo)
