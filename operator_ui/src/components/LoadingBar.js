import React from 'react'
import PropTypes from 'prop-types'
import LinearProgress from '@material-ui/core/LinearProgress'
import { makeStyles } from '@material-ui/styles'

const useStyles = makeStyles(theme => ({
  root: {
    flexGrow: 1,
    height: 3,
    overflow: 'hidden'
  },
  colorPrimary: {
    backgroundColor: theme.palette.background.default
  },
  barColorPrimary: {
    backgroundColor: theme.palette.primary.main
  }
}))
const LoadingBar = ({ fetchCount }) => {
  const classes = useStyles()
  const progressClasses = {
    colorPrimary: classes.colorPrimary,
    barColorPrimary: classes.barColorPrimary
  }

  return (
    <div className={classes.root}>
      {fetchCount > 0 && (
        <LinearProgress variant="indeterminate" classes={progressClasses} />
      )}
    </div>
  )
}

LoadingBar.propTypes = {
  fetchCount: PropTypes.number.isRequired
}

export default LoadingBar
