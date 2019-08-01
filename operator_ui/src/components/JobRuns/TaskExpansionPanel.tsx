import { createStyles } from '@material-ui/core'
import Grid from '@material-ui/core/Grid'
import { withStyles, WithStyles } from '@material-ui/core/styles'
import capitalize from 'lodash/capitalize'
import React from 'react'
import { IJobRun, ITaskRun } from '../../../@types/operator_ui'
import StatusItem from './StatusItem'
import { stringify } from "javascript-stringify";

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
  colATitle: string
  colAValue: string
  colBTitle: string
  colBValue: string | null
}

const Item = withStyles(fontStyles)(
  ({ colATitle, colAValue, colBTitle, colBValue, classes }: IItemProps) => (
    <Grid container>
      <Grid item sm={2}>
        <p className={classes.header}>{colATitle}</p>
        <p className={classes.subHeader}>{stringify(colAValue)}</p>
      </Grid>
      <Grid item md={10}>
        <p className={classes.header}>{colBTitle}</p>
        <p className={classes.subHeader}>{stringify(colBValue) || 'No Value Available'}</p>
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
          colATitle="Initiator Params"
          colAValue="Value"
          colBTitle="Values"
          colBValue="No input Parameters"
        />
      ) : (
        paramsArr.map((par, idx) => (
          <Item
            key={idx}
            colATitle="Initiator Params"
            colAValue={par[0]}
            colBTitle="Values"
            colBValue={par[1]}
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
          colATitle="Params"
          colAValue={p[0]}
          colBTitle="Values"
          colBValue={p[1]}
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
      colATitle="Result"
      colAValue="Task Run Data"
      colBTitle="Values"
      colBValue={result}
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
      {jobRun.taskRuns.map((taskRun: ITaskRun) => {
        return (
          <Grid item xs={12} key={taskRun.id}>
            <StatusItem
              borderTop
              summary={capitalize(taskRun.task.type)}
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
        )
      })}
    </Grid>
  )
}

export default TaskExpansionPanel
