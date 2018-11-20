import React from 'react'
import { withStyles } from '@material-ui/core/styles'
import PaddedCard from 'components/PaddedCard'
import StatusIcon from 'components/JobRuns/StatusIcon'
import Typography from '@material-ui/core/Typography'

const styles = theme => ({
  statusCard: {
    display: 'flex',
    alignItems: 'center'
  },
  statusIcon: {
    display: 'inline-block'
  },
  statusText: {
    display: 'inline-block',
    paddingLeft: theme.spacing.unit * 2
  }
})

const StatusCard = ({classes, children}) => {
  return (
    <PaddedCard className={classes.statusCard}>
      <StatusIcon className={classes.statusIcon}>
        {children}
      </StatusIcon>
      <Typography className={classes.statusText} variant='body1' color='inherit'>
        {children}
      </Typography>
    </PaddedCard>
  )
}

export default withStyles(styles)(StatusCard)
