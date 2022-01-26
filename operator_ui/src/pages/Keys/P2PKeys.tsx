import React from 'react'
import Grid from '@material-ui/core/Grid'
import Button from 'components/Button'
import { v2 } from 'api'
import * as jsonapi from 'utils/json-api-client'
import * as models from 'core/store/models'
import { useErrorHandler } from 'hooks/useErrorHandler'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Typography from '@material-ui/core/Typography'
import {
  createStyles,
  withStyles,
  WithStyles,
  Theme,
} from '@material-ui/core/styles'
import { Delete } from './Delete'
import { KeyBundle } from './KeyBundle'
import { useDispatch } from 'react-redux'
import { deleteNotification, createNotification } from './Notifications'
import { CopyIconButton } from 'components/Copy/CopyIconButton'

const styles = (theme: Theme) =>
  createStyles({
    cardContent: {
      padding: 0,
      '&:last-child': {
        padding: 0,
      },
    },
    avatar: {
      backgroundColor: theme.palette.grey[800],
    },
  })

const KEY_TYPE = 'P2P'

export const P2PKeys = withStyles(styles)(
  ({ classes }: WithStyles<typeof styles>) => {
    const [p2pKeys, setP2Keys] =
      React.useState<jsonapi.ApiResponse<models.P2PKey[]>['data']>()
    const { error, ErrorComponent, setError } = useErrorHandler()
    const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !p2pKeys)
    const dispatch = useDispatch()

    const handleFetchIndex = React.useCallback(() => {
      v2.p2pKeys
        .getP2PKeys()
        .then(({ data }) => {
          setP2Keys(data)
        })
        .catch(setError)
    }, [setError])

    React.useEffect(() => {
      handleFetchIndex()
    }, [handleFetchIndex])

    function handleDelete(id: string) {
      v2.p2pKeys
        .destroyP2PKey(id)
        .then(() => {
          handleFetchIndex()
          dispatch(
            deleteNotification({
              keyType: KEY_TYPE,
            }),
          )
        })
        .catch(setError)
    }

    function handleCreate() {
      v2.p2pKeys
        .createP2PKey()
        .then(({ data }) => {
          handleFetchIndex()
          dispatch(
            createNotification({
              keyType: KEY_TYPE,
              keyValue: data.id,
            }),
          )
        })
        .catch(setError)
    }

    return (
      <Grid item xs={12}>
        <ErrorComponent />
        <LoadingPlaceholder />

        <Card>
          <CardHeader
            action={
              <Button
                data-testid="keys-create"
                variant="secondary"
                onClick={() => handleCreate()}
              >
                New P2P Key
              </Button>
            }
            title={`${KEY_TYPE} Keys`}
            subheader={`Manage your ${KEY_TYPE} Key Bundles.`}
          />
          <CardContent className={classes.cardContent}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>
                    <Typography variant="body1" color="textSecondary">
                      Key Bundle
                    </Typography>
                  </TableCell>
                  <TableCell align="right"></TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {p2pKeys?.length === 0 && (
                  <TableRow>
                    <TableCell component="th" scope="row" colSpan={3}>
                      No entries to show.
                    </TableCell>
                  </TableRow>
                )}
                {p2pKeys?.map((key) => (
                  <TableRow hover key={key.id}>
                    <TableCell>
                      <KeyBundle
                        classes={{ avatar: classes.avatar }}
                        primary={
                          <b>
                            Peer ID: {key.attributes.peerId}{' '}
                            <CopyIconButton data={key.attributes.peerId} />
                          </b>
                        }
                        secondary={[
                          <>Public Key: {key.attributes.publicKey}</>,
                        ]}
                      />
                    </TableCell>
                    <TableCell align="right">
                      <Delete
                        keyId={key.id}
                        keyValue={key.attributes.peerId}
                        onDelete={handleDelete}
                      />
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </Grid>
    )
  },
)
