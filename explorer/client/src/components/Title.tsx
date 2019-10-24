import React from 'react'
import { withStyles, WithStyles, Theme } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import classNames from 'classnames'

const styles = ({ spacing }: Theme) => ({
  title: {
    marginBottom: spacing.unit * 5,
  },
})

interface Props extends WithStyles<typeof styles> {
  className?: string
}

const Title: React.FC<Props> = ({ children, classes, className }) => {
  return (
    <Typography
      variant="h4"
      color="inherit"
      className={classNames(className, classes.title)}
    >
      {children}
    </Typography>
  )
}

export default withStyles(styles)(Title)
