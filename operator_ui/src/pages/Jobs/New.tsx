import React from 'react'
import Radio from '@material-ui/core/Radio'
import { JobSpecFormats, JobSpecFormat } from 'utils/jobSpec'
import { ApiResponse } from '@chainlink/json-api-client'
import Button from 'components/Button'
import TOML from '@iarna/toml'
import * as api from 'api'
import { useDispatch } from 'react-redux'
import {
  OcrJobSpecRequest,
  OcrJobSpec,
  JobSpecRequest,
} from 'core/store/models'
import { JobSpec } from 'core/store/presenters'
import BaseLink from 'components/BaseLink'
import ErrorMessage from 'components/Notifications/DefaultError'
import { notifySuccess, notifyError } from 'actionCreators'
import * as storage from '@chainlink/local-storage'
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
  Divider,
  CardHeader,
  CircularProgress,
} from '@material-ui/core'
import {
  createStyles,
  withStyles,
  WithStyles,
  Theme,
} from '@material-ui/core/styles'
import { useLocation, useHistory } from 'react-router-dom'
import { setPersistJobSpec, getPersistJobSpec } from 'utils/storage'

const jobSpecFormatList = [JobSpecFormats.JSON, JobSpecFormats.TOML]
export const SELECTED_FORMAT = 'persistSpecFormat'

const styles = (theme: Theme) =>
  createStyles({
    card: {
      padding: theme.spacing.unit,
      marginBottom: theme.spacing.unit * 3,
    },
    loader: {
      position: 'absolute',
    },
  })

const SuccessNotification = ({
  id,
}: {
  id: JobSpec['id'] | OcrJobSpec['id']
}) => (
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
  format: JobSpecFormat
  value: string
}) {
  try {
    if (format === JobSpecFormats.JSON) {
      JSON.parse(value)
    } else if (format === JobSpecFormats.TOML) {
      TOML.parse(value)
    }
    return true
  } catch {
    return false
  }
}

function apiCall({
  format,
  value,
}: {
  format: JobSpecFormat
  value: string
}): Promise<ApiResponse<JobSpec | OcrJobSpec>> {
  if (format === JobSpecFormats.JSON) {
    const definition: JobSpecRequest = JSON.parse(value)
    return api.v2.specs.createJobSpec(definition)
  }

  if (format === JobSpecFormats.TOML) {
    const definition: OcrJobSpecRequest = { toml: value }
    return api.v2.ocrSpecs.createJobSpec(definition)
  }

  return Promise.reject('Invalid format')
}

function initialFormat({ query }: { query: string }): JobSpecFormat {
  const params = new URLSearchParams(query)
  return (
    (params.get('format') as JobSpecFormat) ||
    (storage.get(SELECTED_FORMAT) as JobSpecFormat) ||
    JobSpecFormats.JSON
  )
}

function initialJobSpec(): string {
  const jobSpec = getPersistJobSpec()
  try {
    return JSON.stringify(JSON.parse(jobSpec), null, 4)
  } catch {
    return jobSpec || ''
  }
}

export const New = ({
  classes,
}: {
  classes: WithStyles<typeof styles>['classes']
}) => {
  const location = useLocation()
  const [format, setFormat] = React.useState<JobSpecFormat>(
    initialFormat({
      query: location.search,
    }),
  )
  const [value, setValue] = React.useState<string>(initialJobSpec())
  const [valid, setValid] = React.useState<boolean>(true)
  const [loading, setLoading] = React.useState<boolean>(false)
  const dispatch = useDispatch()
  const history = useHistory()

  React.useEffect(() => {
    setPersistJobSpec(value)
    setValid(true)
  }, [value])

  React.useEffect(() => {
    storage.set(SELECTED_FORMAT, format)
    setValid(true)
    history.push({
      search: `?format=${format}`,
    })
  }, [format, history])

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
        <Grid item xs={12} md={11} lg={9}>
          <Card className={classes.card}>
            <CardHeader title="New Job" />
            <Divider />
            <CardContent>
              <form noValidate onSubmit={handleSubmit}>
                <Grid container>
                  <Grid item xs={12}>
                    <FormControl fullWidth>
                      <FormLabel>Job Spec Format</FormLabel>
                      <RadioGroup
                        name="select-format"
                        value={format}
                        onChange={(event: any) =>
                          setFormat(event.target.value as JobSpecFormat)
                        }
                        row
                      >
                        {jobSpecFormatList.map((format) => (
                          <FormControlLabel
                            key={format}
                            value={format}
                            control={<Radio />}
                            label={format}
                          />
                        ))}
                      </RadioGroup>
                    </FormControl>
                  </Grid>
                  <Grid item xs={12}>
                    <FormLabel>Job Spec</FormLabel>
                    <TextField
                      error={!valid}
                      value={value}
                      onChange={(
                        event: React.ChangeEvent<HTMLTextAreaElement>,
                      ) => setValue(event.target.value)}
                      helperText={!valid && `Invalid ${format}`}
                      autoComplete="off"
                      label={`${format} Blob`}
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
                      disabled={loading}
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
      </Grid>
    </Content>
  )
}

export default withStyles(styles)(New)
