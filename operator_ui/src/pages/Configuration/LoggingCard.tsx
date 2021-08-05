import React, { useEffect, useState } from 'react'
import capitalize from 'lodash/capitalize'
import { useDispatch } from 'react-redux'
import { useFormik } from 'formik'
import Button from 'components/Button'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Checkbox from '@material-ui/core/Checkbox'
import FormControlLabel from '@material-ui/core/FormControlLabel'
import FormGroup from '@material-ui/core/FormGroup'
import MenuItem from '@material-ui/core/MenuItem'
import TextField from '@material-ui/core/TextField'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import * as models from 'core/store/models'
import { v2 } from 'api'
import { notifyError, notifySuccess } from 'actionCreators'
import ErrorMessage from 'components/Notifications/DefaultError'
import { useErrorHandler } from 'hooks/useErrorHandler'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'
import { FormHelperText } from '@material-ui/core'

const logLevels = ['debug', 'info', 'warn', 'error']

type FormValues = {
  level: models.LogConfigLevel
  sqlEnabled: boolean
}

const styles = (theme: Theme) => {
  return createStyles({
    actions: {
      display: 'flex',
      justifyContent: 'flex-end',
      marginTop: theme.spacing.unit * 0.5,
    },
    logLevelHelperText: {
      marginTop: -8,
    },
  })
}

interface LogConfigurationFormProps extends WithStyles<typeof styles> {
  initialValues: FormValues
}

const LogConfigurationForm = withStyles(styles)(
  ({ classes, initialValues }: LogConfigurationFormProps) => {
    const dispatch = useDispatch()
    const formik = useFormik({
      initialValues,
      onSubmit: async (values) => {
        try {
          await v2.logConfig.updateLogConfig(values)

          dispatch(notifySuccess(() => <>Logging Configuration Updated</>, {}))
        } catch (e) {
          dispatch(notifyError(ErrorMessage, e))
        }
      },
    })

    return (
      <form onSubmit={formik.handleSubmit}>
        <TextField
          id="select-level"
          name="level"
          fullWidth
          select
          label="Log Level"
          value={formik.values.level}
          onChange={formik.handleChange}
          error={formik.touched.level && Boolean(formik.errors.level)}
          helperText="Override the LOG_LEVEL environment variable"
        >
          {logLevels.map((level) => (
            <MenuItem key={level} value={level}>
              {capitalize(level)}
            </MenuItem>
          ))}
        </TextField>

        <FormGroup>
          <FormControlLabel
            name="sqlEnabled"
            control={
              <>
                <Checkbox
                  id="sqlEnabled"
                  name="sqlEnabled"
                  checked={formik.values.sqlEnabled}
                  onChange={formik.handleChange}
                />
              </>
            }
            label="Log SQL Statements"
          />
          <FormHelperText className={classes.logLevelHelperText}>
            Override the LOG_SQL_STATEMENTS environment variable
          </FormHelperText>
        </FormGroup>

        <br />

        <div className={classes.actions}>
          <Button
            variant="primary"
            type="submit"
            disabled={formik.isSubmitting}
          >
            Update
          </Button>
        </div>
      </form>
    )
  },
)

type LogConfig = {
  level: models.LogConfigLevel
  sqlEnabled: boolean
}

export const LoggingCard = () => {
  const [logConfig, setLogConfig] = useState<LogConfig | null>(null)
  const { error, ErrorComponent, setError } = useErrorHandler()
  const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !logConfig)

  useEffect(() => {
    async function fetch() {
      try {
        const res = await v2.logConfig.getLogConfig()

        // The API interface for getLogConfig is really really bad...
        const globalIdx = res.data.attributes.serviceName.findIndex(
          (name) => name == 'Global',
        )

        const logLevel = res.data.attributes.logLevel[
          globalIdx
        ] as models.LogConfigLevel

        const sqlEnabledIdx = res.data.attributes.serviceName.findIndex(
          (name) => name == 'IsSqlEnabled',
        )

        const sqlEnabled = res.data.attributes.logLevel[sqlEnabledIdx] == 'true'

        const logCfg = {
          level: logLevel,
          sqlEnabled,
        }

        setLogConfig(logCfg)
      } catch (e) {
        setError(e)
      }
    }

    fetch()
  }, [setError])

  return (
    <Card>
      <CardHeader title="Logging" />
      <CardContent>
        <LoadingPlaceholder />
        <ErrorComponent />

        {logConfig && <LogConfigurationForm initialValues={logConfig} />}
      </CardContent>
    </Card>
  )
}
