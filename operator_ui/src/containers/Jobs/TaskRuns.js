import React from 'react'
import Typography from '@material-ui/core/Typography'
import { makeStyles } from '@material-ui/styles'
import StatusIcon from 'components/Icons/TaskStatus'
import titleize from 'utils/titleize'

const useStyles = makeStyles(({ spacing, palette }) => {
  return {
    container: {
      margin: 0,
      marginLeft: spacing(2),
      paddingLeft: 0
    },
    item: {
      borderLeft: 'solid 2px',
      borderColor: palette.grey[200],
      display: 'flex',
      alignItems: 'center',
      listStyle: 'none',
      paddingBottom: spacing(1),
      paddingTop: spacing(1),
      marginLeft: spacing(2)
    },
    status: {
      marginRight: spacing(2),
      marginLeft: -22
    }
  }
})

const TaskRuns = ({ taskRuns }) => {
  const classes = useStyles()
  return (
    <ul className={classes.container}>
      {taskRuns &&
        taskRuns.map(run => {
          return (
            <li key={run.id} className={classes.item}>
              <StatusIcon width={40} className={classes.status}>
                {run.status}
              </StatusIcon>
              <Typography variant="body1" inline>
                {titleize(run.type)}
              </Typography>
            </li>
          )
        })}
    </ul>
  )
}

export default TaskRuns
