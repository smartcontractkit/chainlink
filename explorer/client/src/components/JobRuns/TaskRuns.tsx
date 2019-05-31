import React from 'react'
import Typography from '@material-ui/core/Typography'
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
      flexWrap: 'wrap',
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
      alignItems: 'center',
      flexGrow: 1
    },
    etherscan: {
      marginLeft: spacing.unit
    }
  })

interface IProps extends WithStyles<typeof styles> {
  etherscanHost: string
  taskRuns?: ITaskRun[]
}

const TaskRuns = ({ etherscanHost, taskRuns, classes }: IProps) => {
  return (
    <ul className={classes.container}>
      {taskRuns &&
        taskRuns.map((run: ITaskRun) => {
          return (
            <li key={run.id} className={classes.item}>
              <div className={classes.track}>
                <StatusIcon width={40} className={classes.status}>
                  {run.status}
                </StatusIcon>
                <Typography variant="body1">{run.type}</Typography>
              </div>
              {run.transactionHash && (
                <EtherscanLink
                  txHash={run.transactionHash}
                  host={etherscanHost}
                  className={classes.etherscan}
                />
              )}
            </li>
          )
        })}
    </ul>
  )
}

export default withStyles(styles)(TaskRuns)
