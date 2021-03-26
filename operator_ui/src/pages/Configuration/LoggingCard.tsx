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
import * as models from 'core/store/models'
import { v2 } from 'api'
import { notifyError, notifySuccess } from 'actionCreators'
import ErrorMessage from 'components/Notifications/DefaultError'
import { useErrorHandler } from 'hooks/useErrorHandler'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'

const logLevels = ['debug', 'info', 'warn', 'error']

type FormValues = {
  level: models.LogConfigLevel
  sqlEnabled: boolean
}

const LogConfigurationForm: React.FC<{ initialValues: FormValues }> = ({
  initialValues,
}) => {
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
        helperText={formik.touched.level && formik.errors.level}
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
            <Checkbox
              id="sqlEnabled"
              checked={formik.values.sqlEnabled}
              onChange={formik.handleChange}
            />
          }
          label="Log SQL Statements"
        />
      </FormGroup>

      <Button variant="primary" type="submit" disabled={formik.isSubmitting}>
        Update
      </Button>
    </form>
  )
}

export const LoggingCard = () => {
  const [logConfig, setLogConfig] = useState<models.LogConfig | null>(null)
  const { error, ErrorComponent, setError } = useErrorHandler()
  const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !logConfig)

  useEffect(() => {
    async function fetch() {
      try {
        const res = await v2.logConfig.getLogConfig()

        setLogConfig(res.data.attributes)
      } catch (e) {
        setError(e)
      }
    }

    fetch()
  }, [setError])

  return (
    <Card>
      <CardHeader
        title="Configure Logging"
        subheader="Overrides the LOG_LEVEL and LOG_SQL_STATEMENTS environment variables"
      />
      <CardContent>
        <LoadingPlaceholder />
        <ErrorComponent />

        {logConfig && <LogConfigurationForm initialValues={logConfig} />}
      </CardContent>
    </Card>
  )
}
