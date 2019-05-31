import React from 'react'
import Typography from '@material-ui/core/Typography'
import { createStyles, withStyles } from '@material-ui/core/styles'
import titleize from 'utils/titleize'
import ListIcon from '../../components/Icons/ListIcon'

const styles = ({ spacing, palette }) =>
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
      paddingBottom: spacing.unit * 1,
      paddingTop: spacing.unit * 1,
      marginLeft: spacing.unit * 2
    },
    status: {
      marginRight: spacing.unit * 2,
      marginLeft: -22
    }
  })

const TaskRuns = ({ taskRuns, classes }) => {
  return (
    <ul className={classes.container}>
      {taskRuns &&
        taskRuns.map((run, idx) => {
          return (
            <li key={idx} className={classes.item}>
              <ListIcon width={40} className={classes.status}>
                {run.status}
              </ListIcon>
              <Typography variant="body1" inline>
                {titleize(run.type)}
              </Typography>
            </li>
          )
        })}
    </ul>
  )
}

export default withStyles(styles)(TaskRuns)
