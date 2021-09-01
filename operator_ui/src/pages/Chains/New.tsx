import React, { useState } from 'react'
import { ApiResponse, BadRequestError } from 'utils/json-api-client'
import Button from 'components/Button'
import * as api from 'api'
import { useDispatch } from 'react-redux'
import { CreateChainRequest, Chain } from 'core/store/models'
import BaseLink from 'components/BaseLink'
import ErrorMessage from 'components/Notifications/DefaultError'
import { notifySuccess, notifyError } from 'actionCreators'
import Content from 'components/Content'
import {
  TextField,
  Grid,
  Card,
  CardContent,
  FormLabel,
  CardHeader,
  CircularProgress,
} from '@material-ui/core'
import {
  createStyles,
  withStyles,
  WithStyles,
  Theme,
} from '@material-ui/core/styles'

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
    Successfully created chain{' '}
    <BaseLink id="created-chain" href={`/chains`}>
      {id}
    </BaseLink>
  </>
)

const validate = ({ overrides }: { overrides: string }) => {
  try {
    JSON.parse(overrides)
  } catch (e) {
    return false
  }
  return true
}

function apiCall({
  chainID,
  config,
}: {
  chainID: string
  config: Record<string, JSONPrimitive>
}): Promise<ApiResponse<Chain>> {
  const definition: CreateChainRequest = { chainID, config }
  return api.v2.chains.createChain(definition)
}

export const New = ({
  classes,
}: {
  classes: WithStyles<typeof styles>['classes']
}) => {
  const dispatch = useDispatch()
  const [overrides, setOverrides] = useState<string>('{}')
  const [chainID, setChainID] = useState<string>('')
  const [valid, setValid] = useState<boolean>(true)
  const [valueErrorMsg, setValueErrorMsg] = useState<string>('')
  const [loading, setLoading] = useState<boolean>(false)

  function handleOverrideChange(event: React.ChangeEvent<HTMLTextAreaElement>) {
    setOverrides(event.target.value)
    setValid(true)
    setValueErrorMsg('')
  }

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const isValid = validate({ overrides })
    setValid(isValid)

    if (isValid) {
      setLoading(true)
      apiCall({
        chainID,
        config: JSON.parse(overrides),
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
        <Grid item xs={12}>
          <Card>
            <CardHeader title="New Job" />
            <CardContent>
              <form noValidate onSubmit={handleSubmit}>
                <Grid container>
                  <Grid item xs={12}>
                    <TextField
                      label="Chain ID"
                      name="ID"
                      placeholder="ID"
                      value={chainID}
                      onChange={(e) => setChainID(e.target.value)}
                    />
                  </Grid>
                  <Grid item xs={12}>
                    <FormLabel>Config Overrides</FormLabel>
                    <TextField
                      error={!valid}
                      value={overrides}
                      onChange={handleOverrideChange}
                      helperText={!valid && valueErrorMsg}
                      autoComplete="off"
                      label={'JSON'}
                      rows={10}
                      rowsMax={25}
                      placeholder={'Paste JSON'}
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
                      Create Chain
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
