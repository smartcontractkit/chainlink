import React, { useState } from 'react'
import { ApiResponse, BadRequestError } from 'utils/json-api-client'
import * as api from 'api'
import { useDispatch } from 'react-redux'
import { Chain, CreateChainRequest } from 'core/store/models'
import BaseLink from 'components/BaseLink'
import ErrorMessage from 'components/Notifications/DefaultError'
import { notifyError, notifySuccess } from 'actionCreators'
import Button from 'components/Button'
import Content from 'components/Content'
import {
  Card,
  CardContent,
  CardHeader,
  CircularProgress,
  Grid,
  TextField,
  Typography,
} from '@material-ui/core'
import {
  ChainConfigFields,
  ConfigOverrides,
} from 'pages/Chains/ChainConfigFields'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
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

  const [chainID, setChainID] = useState<string>('')
  const [overrides, setOverrides] = useState<ConfigOverrides>({})
  const [serverErrorMsg, setServerErrorMsg] = useState<string>('')
  const [chainIDErrorMsg, setChainIDErrorMsg] = useState<string>('')
  const [keySpecificOverridesErrorMsg, setKeySpecificOverridesErrorMsg] =
    useState<string>('')
  const [loading, setLoading] = useState<boolean>(false)

  function validate(chainID: string) {
    let valid = true

    if (!(parseInt(chainID, 10) > 0)) {
      setChainIDErrorMsg('Invalid chain ID')
      valid = false
    }

    return valid
  }

  function handleChainIDChange(event: React.ChangeEvent<HTMLTextAreaElement>) {
    setChainID(event.target.value)
    setChainIDErrorMsg('')
  }

  function onConfigChange(config: ConfigOverrides, error: string) {
    if (error) {
      setKeySpecificOverridesErrorMsg(error)
      return
    }

    setOverrides(config)
  }

  function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const isValid = validate(chainID) && !keySpecificOverridesErrorMsg

    if (isValid) {
      setLoading(true)
      setServerErrorMsg('')

      apiCall({
        chainID,
        config: { ...overrides },
      })
        .then(({ data }) => {
          dispatch(notifySuccess(SuccessNotification, data))
        })
        .catch((error) => {
          dispatch(notifyError(ErrorMessage, error))
          if (error instanceof BadRequestError) {
            setServerErrorMsg('Invalid ChainID')
          } else {
            setServerErrorMsg(error.toString())
          }
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
            <CardHeader title="New Chain" />
            <CardContent>
              <form noValidate onSubmit={handleSubmit}>
                <Grid container>
                  {Boolean(serverErrorMsg) && (
                    <Grid item xs={12}>
                      <Typography variant="body1">{serverErrorMsg}</Typography>
                    </Grid>
                  )}

                  <Grid item xs={12}>
                    <TextField
                      error={Boolean(chainIDErrorMsg)}
                      helperText={Boolean(chainIDErrorMsg)}
                      label="Chain ID"
                      name="ID"
                      placeholder="ID"
                      value={chainID}
                      onChange={handleChainIDChange}
                    />
                  </Grid>

                  <Grid item xs={false} md={12}></Grid>

                  <ChainConfigFields onChange={onConfigChange} />

                  <Grid item xs={12}>
                    <Button
                      data-testid="new-chain-config-submit"
                      variant="primary"
                      type="submit"
                      size="large"
                      disabled={
                        loading ||
                        Boolean(keySpecificOverridesErrorMsg) ||
                        Boolean(chainIDErrorMsg)
                      }
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
