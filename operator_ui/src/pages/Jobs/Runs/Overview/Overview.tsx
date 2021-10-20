import React from 'react'
import { createStyles } from '@material-ui/core'
import Card from '@material-ui/core/Card'
import Grid from '@material-ui/core/Grid'
import { withStyles, WithStyles } from '@material-ui/core/styles'
import { stringify } from 'javascript-stringify'
import capitalize from 'lodash/capitalize'
import StatusItem from './StatusItem'
import { DirectRequestJobRun } from '../../sharedTypes'

const fontStyles = () =>
  createStyles({
    header: {
      fontFamily: 'Roboto Mono',
      fontWeight: 'bold',
      fontSize: '14px',
      color: '#818ea3',
    },
    subHeader: {
      fontFamily: 'Roboto Mono',
      fontSize: '12px',
      color: '#818ea3',
    },
  })

interface ItemProps extends WithStyles<typeof fontStyles> {
  colATitle: string
  colAValue: string
  colBTitle: string
  colBValue?: string
}

const Item = withStyles(fontStyles)(
  ({ colATitle, colAValue, colBTitle, colBValue, classes }: ItemProps) => (
    <Grid container>
      <Grid item sm={2}>
        <p className={classes.header}>{colATitle}</p>
        <p className={classes.subHeader}>{stringify(colAValue)}</p>
      </Grid>
      <Grid item md={10}>
        <p className={classes.header}>{colBTitle}</p>
        <p className={classes.subHeader}>
          {stringify(colBValue) || 'No Value Available'}
        </p>
      </Grid>
    </Grid>
  ),
)

interface InitiatorProps {
  params: object
}

const Initiator = ({ params }: InitiatorProps) => {
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

interface ParamsProps {
  params?: object
}

const Params = ({ params }: ParamsProps) => {
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

interface ResultProps {
  run: DirectRequestJobRun['taskRuns'][0]
}

const Result = ({ run }: ResultProps) => {
  const data = run?.result?.data
  const result = data && 'result' in data ? data.result : ''

  return (
    <Item
      colATitle="Result"
      colAValue="Task Run Data"
      colBTitle="Values"
      colBValue={result}
    />
  )
}

interface Props {
  jobRun: DirectRequestJobRun
}

export const Overview = ({ jobRun }: Props) => {
  const initiator = jobRun.initiator

  return (
    <Card>
      <Grid container spacing={0}>
        <Grid item xs={12}>
          <StatusItem
            summary={capitalize(initiator.type)}
            status={jobRun.status}
            borderTop={false}
            confirmations={0}
            minConfirmations={0}
          >
            <div
              style={{
                display: 'flex',
                flexDirection: 'column',
                width: '100%',
              }}
            >
              <Initiator params={initiator.params || {}} />
            </div>
          </StatusItem>
        </Grid>
        {jobRun.taskRuns.map((taskRun) => {
          return (
            <Grid item xs={12} key={taskRun.id}>
              <StatusItem
                borderTop
                summary={capitalize(taskRun.task.type)}
                status={taskRun.status}
                confirmations={taskRun.confirmations}
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
    </Card>
  )
}
