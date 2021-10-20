import React, { useEffect, useState } from 'react'
import Radio from '@material-ui/core/Radio'
import {
  JobSpecFormats,
  JobSpecFormat,
  getJobSpecFormat,
  isJson,
  isToml,
  getTaskList,
} from './utils'
import { ApiResponse, BadRequestError } from 'utils/json-api-client'
import Button from 'components/Button'
import * as api from 'api'
import { useDispatch } from 'react-redux'
import { JobSpecV2Request, JobSpecV2, JobSpecRequest } from 'core/store/models'
import { JobSpec } from 'core/store/presenters'
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
  FormControlLabel,
  FormControl,
  FormLabel,
  RadioGroup,
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
import { useLocation, useHistory } from 'react-router-dom'
import { TaskSpec } from 'core/store/models'
import TaskListDag from './TaskListDag'
import TaskList from 'components/Jobs/TaskList'
import { Stratify } from './parseDot'

const jobSpecFormatList = [JobSpecFormats.JSON, JobSpecFormats.TOML]

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

export function validate({
  format,
  value,
}: {
  format: JobSpecFormats
  value: string
}) {
  if (value.trim() === '') {
    return false
  } else if (format === JobSpecFormats.JSON && isJson({ value })) {
    return true
  } else if (format === JobSpecFormats.TOML && isToml({ value })) {
    return true
  } else {
    return false
  }
}

function apiCall({
  format,
  value,
}: {
  format: JobSpecFormats
  value: string
}): Promise<ApiResponse<JobSpec | JobSpecV2>> {
  if (format === JobSpecFormats.JSON) {
    const definition: JobSpecRequest = JSON.parse(value)
    return api.v2.specs.createJobSpec(definition)
  }

  if (format === JobSpecFormats.TOML) {
    const definition: JobSpecV2Request = { toml: value }
    return api.v2.jobs.createJobSpec(definition)
  }

  return Promise.reject('Invalid format')
}

function getInitialValues({
  query,
}: {
  query: string
}): { jobSpec: string; format: JobSpecFormats } {
  const params = new URLSearchParams(query)
  const queryJobSpec = params.get('definition') as string
  const queryJobSpecFormat =
    getJobSpecFormat({
      value: queryJobSpec,
    }) || JobSpecFormats.JSON

  if (queryJobSpec) {
    storage.set(`${PERSIST_SPEC}${queryJobSpecFormat}`, queryJobSpec)
    return {
      jobSpec: queryJobSpec,
      format: queryJobSpecFormat,
    }
  }

  const lastOpenedFormat =
    JobSpecFormats[params.get('format')?.toUpperCase() as JobSpecFormat] ||
    storage.get(SELECTED_FORMAT) ||
    JobSpecFormats.JSON

  const lastOpenedJobSpec =
    storage.get(`${PERSIST_SPEC}${lastOpenedFormat}`) || ''

  return {
    jobSpec: lastOpenedJobSpec,
    format: lastOpenedFormat,
  }
}

export const New = ({
  classes,
}: {
  classes: WithStyles<typeof styles>['classes']
}) => {
  const dispatch = useDispatch()
  const history = useHistory()
  const location = useLocation()
  const [initialValues] = useState(() =>
    getInitialValues({
      query: location.search,
    }),
  )
  const [format, setFormat] = useState<JobSpecFormats>(initialValues.format)
  const [value, setValue] = useState<string>(initialValues.jobSpec)
  const [valid, setValid] = useState<boolean>(true)
  const [valueErrorMsg, setValueErrorMsg] = useState<string>('')
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

  // Change the form to use either JSON or TOML format
  function handleFormatChange(_event: React.ChangeEvent<{}>, format: string) {
    setValue(storage.get(`${PERSIST_SPEC}${format}`) || '')
    setFormat(format as JobSpecFormats)
    storage.set(SELECTED_FORMAT, format)
    setValid(true)
    history.replace({
      search: `?format=${format}`,
    })
  }

  // Update the job spec value
  function handleValueChange(event: React.ChangeEvent<HTMLTextAreaElement>) {
    setValue(event.target.value)
    storage.set(`${PERSIST_SPEC}${format}`, event.target.value)
    setValid(true)
  }

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const isValid = validate({ format, value })
    setValid(isValid)

    if (isValid) {
      setLoading(true)

      apiCall({
        format,
        value,
      })
        .then(({ data }) => {
          dispatch(notifySuccess(SuccessNotification, data))
        })
        .catch((error) => {
          dispatch(notifyError(ErrorMessage, error))
          if (error instanceof BadRequestError) {
            setValueErrorMsg('Invalid JSON')
          } else {
            setValueErrorMsg(error.toString())
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
                    <FormControl fullWidth>
                      <FormLabel>Job Spec Format</FormLabel>
                      <RadioGroup
                        name="select-format"
                        value={format}
                        onChange={handleFormatChange}
                        row
                      >
                        {jobSpecFormatList.map((format) => (
                          <FormControlLabel
                            key={format}
                            value={format}
                            control={<Radio />}
                            label={format}
                            checked={format === 'toml'}
                          />
                        ))}
                      </RadioGroup>
                    </FormControl>
                    <b>
                      NOTE: Support for JSON jobs has been deprecated. These
                      legacy job types will be disabled entirely in an upcoming
                      release. For help migrating existing jobs to the TOML
                      format, please see the{' '}
                      <a href="https://docs.chain.link/docs/jobs/">docs</a>.
                    </b>
                  </Grid>
                  <Grid item xs={12}>
                    <FormLabel>Job Spec</FormLabel>
                    <TextField
                      error={!valid}
                      value={value}
                      onChange={handleValueChange}
                      helperText={!valid && valueErrorMsg}
                      autoComplete="off"
                      label={`${format} blob`}
                      rows={10}
                      rowsMax={25}
                      placeholder={`Paste ${format}`}
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
                      data-testid="new-job-spec-submit"
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
            {tasks.format === JobSpecFormats.JSON && tasks.list && (
              <TaskList tasks={tasks.list as TaskSpec[]} />
            )}
            {tasks.format === JobSpecFormats.TOML && tasks.list && (
              <TaskListDag stratify={tasks.list as Stratify[]} />
            )}
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
