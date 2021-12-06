import React, { useEffect, useState } from 'react'
import { isToml, getTaskList } from './utils'
import { ApiResponse, BadRequestError } from 'utils/json-api-client'
import Button from 'components/Button'
import * as api from 'api'
import { useDispatch } from 'react-redux'
import { CreateJobRequest, Job } from 'core/store/models'
import BaseLink from 'components/BaseLink'
import ErrorMessage from 'components/Notifications/DefaultError'
import { notifySuccess, notifyError } from 'actionCreators'
import * as storage from 'utils/local-storage'
import Content from 'components/Content'
import {
  TextField,
  Grid,
  Card,
  CardContent,
  FormLabel,
  CardHeader,
  CircularProgress,
  Typography,
} from '@material-ui/core'
import {
  createStyles,
  withStyles,
  WithStyles,
  Theme,
} from '@material-ui/core/styles'
import { useLocation } from 'react-router-dom'
import TaskListDag from './TaskListDag'
import { Stratify } from 'utils/parseDot'

export const SELECTED_FORMAT = 'persistSpec.format'
export const PERSIST_SPEC = 'persistSpec.'

const styles = (theme: Theme) =>
  createStyles({
    loader: {
      position: 'absolute',
    },
    emptyTasks: {
      padding: theme.spacing.unit * 3,
    },
  })

const SuccessNotification = ({ id }: { id: string }) => (
  <>
    Successfully created job{' '}
    <BaseLink id="created-job" href={`/jobs/${id}`}>
      {id}
    </BaseLink>
  </>
)

export function validate({ value }: { value: string }) {
  if (value.trim() === '') {
    return false
  } else if (isToml({ value })) {
    return true
  } else {
    return false
  }
}

function apiCall({ value }: { value: string }): Promise<ApiResponse<Job>> {
  const definition: CreateJobRequest = { toml: value }
  return api.v2.jobs.createJobSpec(definition)
}

function getInitialValues({ query }: { query: string }): { jobSpec: string } {
  const params = new URLSearchParams(query)
  const queryJobSpec = params.get('definition') as string

  if (queryJobSpec) {
    storage.set(PERSIST_SPEC, queryJobSpec)
    return {
      jobSpec: queryJobSpec,
    }
  }

  const lastOpenedJobSpec = storage.get(`${PERSIST_SPEC}`) || ''

  return {
    jobSpec: lastOpenedJobSpec,
  }
}

export const New = ({
  classes,
}: {
  classes: WithStyles<typeof styles>['classes']
}) => {
  const dispatch = useDispatch()
  const location = useLocation()
  const [initialValues] = useState(() =>
    getInitialValues({
      query: location.search,
    }),
  )
  const [value, setValue] = useState<string>(initialValues.jobSpec)
  const [valid, setValid] = useState<boolean>(true)
  const [valueErrorMsg, setValueErrorMsg] = useState<string>('')
  const [serverErrorMsg, setServerErrorMsg] = useState<string>('')
  const [loading, setLoading] = useState<boolean>(false)
  const [tasks, setTasks] = useState(() =>
    getTaskList({ value: initialValues.jobSpec }),
  )

  // Extract the tasks from the job spec to display in the preview
  useEffect(() => {
    const timeout = setTimeout(() => {
      setValueErrorMsg('')
      const taskList = getTaskList({ value })
      if (taskList.error) {
        setValid(false)
        setValueErrorMsg(taskList.error)
      } else {
        setTasks(taskList)
      }
    }, 500)

    return () => clearTimeout(timeout)
  }, [value, setTasks])

  // Update the job spec value
  function handleValueChange(event: React.ChangeEvent<HTMLTextAreaElement>) {
    const noWhiteSpaceValue = event.target.value.replace(
      /[\u200B-\u200D\uFEFF]/g,
      '',
    )
    setValue(noWhiteSpaceValue)
    storage.set(`${PERSIST_SPEC}`, noWhiteSpaceValue)
    setValid(true)
  }

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const isValid = validate({ value })
    setValid(isValid)

    if (isValid) {
      setLoading(true)
      setServerErrorMsg('')

      await apiCall({
        value,
      })
        .then(({ data }) => {
          dispatch(notifySuccess(SuccessNotification, data))
        })
        .catch((error) => {
          dispatch(notifyError(ErrorMessage, error))
          if (error instanceof BadRequestError) {
            setServerErrorMsg('Invalid job spec')
          } else {
            setServerErrorMsg(error.toString())
          }

          setValid(false)
        })
        .finally(() => {
          setLoading(false)
        })
    }
  }

  return (
    <Content>
      <Grid container spacing={40}>
        <Grid item xs={12} lg={8}>
          <Card>
            <CardHeader title="New Job" />
            <CardContent>
              <form noValidate onSubmit={handleSubmit}>
                <Grid container>
                  <Grid item xs={12}>
                    <b>
                      NOTE: Support for adding JSON jobs has been removed. For
                      help writing jobs in TOML format, please see the{' '}
                      <a href="https://docs.chain.link/docs/jobs/">docs</a>.
                    </b>
                  </Grid>
                  <Grid item xs={12}>
                    <FormLabel>Job Spec</FormLabel>
                    <TextField
                      error={!valid}
                      value={value}
                      onChange={handleValueChange}
                      helperText={(!valid && valueErrorMsg) || serverErrorMsg}
                      autoComplete="off"
                      label={'TOML blob'}
                      rows={10}
                      rowsMax={25}
                      placeholder={'Paste TOML'}
                      multiline
                      margin="normal"
                      name="jobSpec"
                      id="jobSpec"
                      variant="outlined"
                      fullWidth
                    />
                  </Grid>
                  <Grid item xs={12}>
                    <Button
                      variant="primary"
                      type="submit"
                      size="large"
                      disabled={loading || Boolean(valueErrorMsg)}
                    >
                      Create Job
                      {loading && (
                        <CircularProgress
                          className={classes.loader}
                          size={30}
                          color="primary"
                        />
                      )}
                    </Button>
                  </Grid>
                </Grid>
              </form>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} lg={4}>
          <Card style={{ overflow: 'visible' }}>
            <CardHeader title="Task list preview" />
            {tasks.list && <TaskListDag stratify={tasks.list as Stratify[]} />}
            {!tasks.list && (
              <Typography
                className={classes.emptyTasks}
                variant="body1"
                color="textSecondary"
              >
                Tasks not found
              </Typography>
            )}
          </Card>
        </Grid>
      </Grid>
    </Content>
  )
}

export default withStyles(styles)(New)
