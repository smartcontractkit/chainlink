import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import { TaskSpec } from 'core/store/models'
import React from 'react'
import ListIcon from '../../components/Icons/ListIcon'
import titleize from '../../utils/titleize'

const styles = ({ spacing, palette }: Theme) =>
  createStyles({
    container: {
      margin: 0,
      marginLeft: spacing.unit * 2,
      paddingLeft: 0,
    },
    item: {
      borderLeft: 'solid 2px',
      borderColor: palette.grey[200],
      display: 'flex',
      alignItems: 'center',
      listStyle: 'none',
      paddingBottom: spacing.unit * 1,
      paddingTop: spacing.unit * 1,
      marginLeft: spacing.unit * 2,
    },
    status: {
      marginRight: spacing.unit * 2,
      marginLeft: -22,
    },
  })

interface Props extends WithStyles<typeof styles> {
  tasks: TaskSpec[]
}

const TaskList = ({ tasks, classes }: Props) => {
  return (
    <ul className={classes.container}>
      {tasks &&
        tasks.map((task, idx) => {
          return (
            <li key={idx} className={classes.item} data-testid="task-list-item">
              <ListIcon width={40} className={classes.status} />
              <Typography variant="body1" inline>
                {titleize(task.type)}
              </Typography>
            </li>
          )
        })}
    </ul>
  )
}

export default withStyles(styles)(TaskList)
