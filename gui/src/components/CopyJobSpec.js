import React from 'react'
import { withStyles } from '@material-ui/core/styles'
import Copy from 'components/Copy'

const styles = theme => ({
  button: {
    margin: theme.spacing.unit
  }
})

const CopyJobSpec = ({classes, JobSpec}) => {
  return (
    <div className={classes.button}>
      <Copy buttonText='Copy JobSpec' data={JobSpec} />
    </div>
  )
}

export default withStyles(styles)(CopyJobSpec)
