import React from 'react'
import Grid from '@material-ui/core/Grid'
import StatusItem from 'components/JobRuns/StatusItem'
import PrettyJson from 'components/PrettyJson'
import capitalize from 'lodash/capitalize'
import { IInitiator, ITaskRuns, IJobRun } from '../../../@types/operator_ui'
import { createStyles } from '@material-ui/core'
import { withStyles, WithStyles } from '@material-ui/core/styles'

const renderConfirmations = (
  confirmations: number,
  minimumConfirmations: number
) => {
  if (minimumConfirmations) {
    return (
      <div>
        {confirmations}/{minimumConfirmations} confirmations
      </div>
    )
  }
}

const fontStyles = () =>
  createStyles({
    header: {
      fontFamily: 'Roboto Mono',
      fontWeight: 'bold',
      fontSize: '14px',
      color: '#818ea3'
    },
    subHeader: { fontFamily: 'Roboto Mono', fontSize: '12px', color: '#818ea3' }
  })

interface IItemProps extends WithStyles<typeof fontStyles> {
  keyOne: string
  valOne: string
  keyTwo: string
  valTwo: string
}

const Item = withStyles(fontStyles)(
  ({ keyOne, valOne, keyTwo, valTwo, classes }: IItemProps) => (
    <Grid container>
      <Grid item sm={2}>
        <p className={classes.header}>{keyOne}</p>
        <p className={classes.subHeader}>{valOne}</p>
      </Grid>
      <Grid item md={10}>
        <p className={classes.header}>{keyTwo}</p>
        <p className={classes.subHeader}>{valTwo}</p>
      </Grid>
    </Grid>
  )
)

const renderParams = (params: object) => {
  return (
    <>
      {Object.entries(params).map(par => (
        <Item keyOne="Params" valOne={par[0]} keyTwo="Values" valTwo={par[1]} />
      ))}
    </>
  )
}

const renderResult = (result: string) => (
  <Item keyOne="Result" valOne="Value" keyTwo="Values" valTwo={result} />
)

const TaskExpansionPanel = ({ children }: { children: IJobRun }) => {
  const initiator: IInitiator = children.initiator
  const taskRuns: ITaskRuns = children.taskRuns

  return (
    <Grid container spacing={0}>
      <Grid item xs={12}>
        <StatusItem
          summary={capitalize(initiator.type)}
          status={children.status}
          borderTop={false}
        >
          <PrettyJson object={initiator.params} />
        </StatusItem>
      </Grid>
      {taskRuns.map(taskRun => (
        <Grid item xs={12} key={taskRun.id}>
          <StatusItem
            summary={capitalize(taskRun.task.type)}
            status={taskRun.status}
          >
            {renderConfirmations(
              taskRun.task.confirmations,
              taskRun.minimumConfirmations
            )}
            <Grid container direction="column">
              <Grid item>{renderParams(taskRun.task.params)}</Grid>
              <Grid item>{renderResult(taskRun.result.data.result)}</Grid>
            </Grid>
          </StatusItem>
        </Grid>
      ))}
    </Grid>
  )
}

export default TaskExpansionPanel
