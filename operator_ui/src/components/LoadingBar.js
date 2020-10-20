import React from 'react'
import PropTypes from 'prop-types'
import LinearProgress from '@material-ui/core/LinearProgress'
import { withStyles } from '@material-ui/core/styles'

const styles = (theme) => ({
  root: {
    flexGrow: 1,
    height: 3,
    overflow: 'hidden',
  },
  colorPrimary: {
    backgroundColor: theme.palette.background.default,
  },
  barColorPrimary: {
    backgroundColor: theme.palette.primary.main,
  },
})

const LoadingBar = ({ classes, fetchCount }) => {
  const progressClasses = {
    colorPrimary: classes.colorPrimary,
    barColorPrimary: classes.barColorPrimary,
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
  fetchCount: PropTypes.number.isRequired,
}

export default withStyles(styles)(LoadingBar)
