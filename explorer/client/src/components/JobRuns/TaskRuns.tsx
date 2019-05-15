import React from 'react'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import StatusIcon from '../Icons/Status'
import EtherscanLink from './EtherscanLink'

const styles = ({ spacing, palette }: Theme) =>
  createStyles({
    container: {
      margin: 0,
      marginLeft: spacing.unit * 2,
      paddingLeft: 0
    },
    item: {
      borderLeft: 'solid 2px',
      borderColor: palette.grey[200],
      display: 'flex',
      alignItems: 'center',
      listStyle: 'none',
      paddingBottom: spacing.unit * 2,
      '&:last-child': {
        paddingBottom: 0
      }
    },
    status: {
      marginRight: spacing.unit * 2,
      marginLeft: -22
    },
    track: {
      display: 'flex',
      alignItems: 'center'
    }
  })

interface IProps extends WithStyles<typeof styles> {
  taskRuns?: ITaskRun[]
}

const TaskRuns = ({ taskRuns, classes }: IProps) => {
  return (
    <ul className={classes.container}>
      {taskRuns &&
        taskRuns.map((run: ITaskRun) => {
          return (
            <li key={run.id} className={classes.item}>
              <Grid container alignItems="center">
                <Grid item xs={3}>
                  <div className={classes.track}>
                    <StatusIcon width={40} className={classes.status}>
                      {run.status}
                    </StatusIcon>
                    <Typography variant="body1">{run.type}</Typography>
                  </div>
                </Grid>
                <Grid item xs={9}>
                  {run.transactionHash && (
                    <EtherscanLink txHash={run.transactionHash} />
                  )}
                </Grid>
              </Grid>
            </li>
          )
        })}
    </ul>
  )
}

export default withStyles(styles)(TaskRuns)
