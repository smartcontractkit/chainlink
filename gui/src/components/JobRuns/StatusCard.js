import React from 'react'
import { withStyles } from '@material-ui/core/styles'
import PaddedCard from 'components/PaddedCard'
import StatusIcon from 'components/JobRuns/StatusIcon'
import Typography from '@material-ui/core/Typography'
import classNames from 'classnames'
import { titleCase } from 'change-case'

const styles = theme => ({
  completed: {
    backgroundColor: theme.palette.success.light
  },
  errored: {
    backgroundColor: theme.palette.error.light
  },
  pending: {
    backgroundColor: theme.palette.warning.light
  },
  statusCard: {
    display: 'flex',
    alignItems: 'center',
    '&:last-child': {
      paddingBottom: theme.spacing.unit * 2
    }
  },
  statusIcon: {
    display: 'inline-block'
  },
  statusText: {
    display: 'inline-block',
    paddingLeft: theme.spacing.unit * 2,
    textTransform: 'capitalize'
  }
})

const StatusCard = ({ classes, children }) => {
  return (
    <PaddedCard
      className={classNames(
        classes.statusCard,
        classes[children] || classes.pending
      )}
    >
      <StatusIcon className={classes.statusIcon} width={80}>
        {children}
      </StatusIcon>
      <Typography className={classes.statusText} variant="h5" color="inherit">
        {titleCase(children)}
      </Typography>
    </PaddedCard>
  )
}

export default withStyles(styles)(StatusCard)
