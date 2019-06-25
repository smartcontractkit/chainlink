import React from 'react'
import Grid from '@material-ui/core/Grid'
import StatusItem from './StatusItem'
import capitalize from 'lodash/capitalize'
import {
  IInitiator,
  ITaskRuns,
  ITaskRun,
  IJobRun
} from '../../../@types/operator_ui'
import { createStyles } from '@material-ui/core'
import { withStyles, WithStyles } from '@material-ui/core/styles'

const fontStyles = () =>
  createStyles({
    header: {
      fontFamily: 'Roboto Mono',
      fontWeight: 'bold',
      fontSize: '14px',
      color: '#818ea3'
    },
    subHeader: {
      fontFamily: 'Roboto Mono',
      fontSize: '12px',
      color: '#818ea3'
    }
  })

interface IItemProps extends WithStyles<typeof fontStyles> {
  keyOne: string
  valOne: string
  keyTwo: string
  valTwo: string | null
}

const Item = withStyles(fontStyles)(
  ({ keyOne, valOne, keyTwo, valTwo, classes }: IItemProps) => (
    <Grid container>
      <Grid item sm={2}>
        <p className={classes.header}>{keyOne}</p>
        <p className={classes.subHeader}>{valOne || 'No Value Available'}</p>
      </Grid>
      <Grid item md={10}>
        <p className={classes.header}>{keyTwo}</p>
        <p className={classes.subHeader}>{valTwo || 'No Value Available'}</p>
      </Grid>
    </Grid>
  )
)

interface IInitiatorProps {
  params: object
}

const Initiator = ({ params }: IInitiatorProps) => {
  const paramsArr = Object.entries(params)

  return (
    <>
      {JSON.stringify(paramsArr) === '[]' ? (
        <Item
          keyOne="Initiator Params"
          valOne="Value"
          keyTwo="Values"
          valTwo="No input Parameters"
        />
      ) : (
        paramsArr.map((par, idx) => (
          <Item
            key={idx}
            keyOne="Initiator Params"
            valOne={par[0]}
            keyTwo="Values"
            valTwo={par[1]}
          />
        ))
      )}
    </>
  )
}

interface IParamsProps {
  params?: object
}

const Params = ({ params }: IParamsProps) => {
  return (
    <div>
      {Object.entries(params || {}).map((p, idx) => (
        <Item
          key={idx}
          keyOne="Params"
          valOne={p[0]}
          keyTwo="Values"
          valTwo={p[1]}
        />
      ))}
    </div>
  )
}

interface IResultProps {
  run: ITaskRun
}

const Result = ({ run }: IResultProps) => {
  const result = run.result && run.result.data && run.result.data.result

  return (
    <Item
      keyOne="Result"
      valOne="Task Run Data"
      keyTwo="Values"
      valTwo={result}
    />
  )
}

interface IProps {
  jobRun: IJobRun
}

const TaskExpansionPanel = ({ jobRun }: IProps) => {
  const initiator = jobRun.initiator

  return (
    <Grid container spacing={0}>
      <Grid item xs={12}>
        <StatusItem
          summary={capitalize(initiator.type)}
          status={jobRun.status}
          borderTop={false}
          confirmations={0}
          minConfirmations={0}
        >
          <Initiator params={initiator.params} />
        </StatusItem>
      </Grid>
      {jobRun.taskRuns.map((taskRun: ITaskRun) => (
        <Grid item xs={12} key={taskRun.id}>
          <StatusItem
            borderTop
            summary={capitalize(taskRun.type)}
            status={taskRun.status}
            confirmations={taskRun.task.confirmations}
            minConfirmations={taskRun.minimumConfirmations}
          >
            <Grid container direction="column">
              <Grid item>
                {taskRun.task && <Params params={taskRun.task.params} />}
              </Grid>
              <Grid item>
                <Result run={taskRun} />
              </Grid>
            </Grid>
          </StatusItem>
        </Grid>
      ))}
    </Grid>
  )
}

export default TaskExpansionPanel
