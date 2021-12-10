import React from 'react'

import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'

import { TaskRunStatusIcon } from 'src/components/Icons/TaskRunStatusIcon'

const styles = (theme: Theme) =>
  createStyles({
    root: {
      display: 'flex',
      borderBottom: '1px solid rgba(0, 0, 0, 0.12);',
      padding: theme.spacing.unit * 2,
      '&:last-child': {
        borderBottom: 'none',
      },
    },
    content: {
      flex: 1,
      overflow: 'auto',
      paddingLeft: theme.spacing.unit * 2,
    },
    heading: {
      lineHeight: `${theme.spacing.unit * 5}px`,
    },
    subheading: {
      color: theme.palette.grey[500],
    },
    text: {
      marginBottom: theme.spacing.unit,
      marginTop: theme.spacing.unit,
    },
    attributeKey: {
      fontWeight: theme.typography.fontWeightLight,
    },
  })

export interface Props
  extends WithStyles<typeof styles>,
    Pick<JobRunPayload_TaskRunsFields, 'dotID' | 'output' | 'error' | 'type'> {
  attrs?: object
}

export const TaskRunItem = withStyles(styles)(
  ({ attrs, classes, dotID, output, error, type }: Props) => {
    const status = error ? 'errored' : 'completed'

    return (
      <div className={classes.root}>
        <div>
          <TaskRunStatusIcon status={status} width={40} height={40} />
        </div>

        <div className={classes.content}>
          <Typography className={classes.heading} variant="h5">
            {dotID} <small className={classes.subheading}>{type}</small>
          </Typography>

          {['completed', 'errored'].includes(status) && (
            <Typography className={classes.text} variant="body1">
              {status === 'completed' ? output : error}
            </Typography>
          )}

          {attrs &&
            Object.entries(attrs).map(([key, value]) => {
              if (key === 'type') {
                return null
              }

              return (
                <Typography key={key} variant="body1">
                  <span className={classes.attributeKey}>{key}</span>: {value}
                </Typography>
              )
            })}
        </div>
      </div>
    )
  },
)
