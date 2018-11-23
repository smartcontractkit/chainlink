import React from 'react'
import { withStyles } from '@material-ui/core/styles'
import Tooltip from '@material-ui/core/Tooltip'

const styles = theme => ({
  lightTooltip: {
    background: theme.palette.primary.contrastText,
    color: theme.palette.text.primary,
    boxShadow: theme.shadows[24],
    ...theme.typography.h6
  }
})

const StyledTooltip = ({ title, children, classes }) => {
  return (
    <div>
      <Tooltip title={title} classes={{ tooltip: classes.lightTooltip }}>
        {children}
      </Tooltip>
    </div>
  )
}

export default withStyles(styles)(StyledTooltip)
