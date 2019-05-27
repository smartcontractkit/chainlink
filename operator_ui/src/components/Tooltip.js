import React from 'react'
import Tooltip from '@material-ui/core/Tooltip'
import { makeStyles } from '@material-ui/styles'

const useStyles = makeStyles(theme => ({
  lightTooltip: {
    background: theme.palette.primary.contrastText,
    color: theme.palette.text.primary,
    boxShadow: theme.shadows[24],
    ...theme.typography.h6
  }
}))

const StyledTooltip = ({ title, children }) => {
  const classes = useStyles()
  return (
    <Tooltip title={title} classes={{ tooltip: classes.lightTooltip }}>
      {children}
    </Tooltip>
  )
}

export default StyledTooltip
