import React from 'react'
import NotFoundSVG from 'images/four-oh-four.js'
import { makeStyles } from '@material-ui/styles'

const useStyles = makeStyles(() => ({
  logo: {
    top: '30%',
    left: '50%',
    transform: 'translate(-50%, -30%)',
    position: 'absolute'
  }
}))

const Logo = () => {
  const classes = useStyles()
  return (
    <div className={classes.logo}>
      <NotFoundSVG />
    </div>
  )
}

export default Logo
