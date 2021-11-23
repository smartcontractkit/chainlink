import React from 'react'

import { createStyles, withStyles, WithStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'

const styles = () => {
  return createStyles({
    root: {
      fontSize: 24,
    },
  })
}

interface Props extends WithStyles<typeof styles> {}

export const Heading2 = withStyles(styles)(
  ({ children, classes }: React.PropsWithChildren<Props>) => (
    <Typography variant="h2" className={classes.root}>
      {children}
    </Typography>
  ),
)
