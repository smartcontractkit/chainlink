import React from 'react'
import { withStyles, WithStyles } from '@material-ui/core/styles'
import MuiTooltip from '@material-ui/core/Tooltip'

const styles = ({ palette, shadows, typography }) => ({
  lightTooltip: {
    background: palette.primary.contrastText,
    color: palette.text.primary,
    boxShadow: shadows[24],
    ...typography.h6
  }
})

const Tooltip = ({ title, children, classes }) => {
  return (
    <MuiTooltip title={title} classes={{ tooltip: classes.lightTooltip }}>
      {children}
    </MuiTooltip>
  )
}

export default withStyles(styles)(Tooltip)
