import React, { useEffect, useState } from 'react'
import { useHistory, useParams } from 'react-router-dom'
import { useDispatch } from 'react-redux'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Dialog from '@material-ui/core/Dialog'
import DialogActions from '@material-ui/core/DialogActions'
import DialogContent from '@material-ui/core/DialogContent'
import DialogTitle from '@material-ui/core/DialogTitle'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import {
  createStyles,
  withStyles,
  WithStyles,
  Theme,
} from '@material-ui/core/styles'
import BaseLink from 'components/BaseLink'
import Button from 'components/Button'
import Content from 'components/Content'
import * as models from 'core/store/models'
import { v2 } from 'api'
import { ApiResponse } from 'utils/json-api-client'
import { notifyError, notifySuccess } from 'actionCreators'
import ErrorMessage from 'components/Notifications/DefaultError'

export const styles = (theme: Theme) =>
  createStyles({
    action: {
      marginLeft: theme.spacing.unit,
    },
  })

const bridgeDetailsStyles = (theme: Theme) =>
  createStyles({
    headerRoot: {
      borderBottom: `1px solid ${theme.palette.divider}`,
    },
    tableRow: {
      '& td:first-child': {
        fontWeight: theme.typography.fontWeightMedium,
      },
    },
  })

const Loading = ({ isErrored }: { isErrored: boolean }) => (
  <CardContent>{isErrored ? 'Bridge not found' : 'Loading...'}</CardContent>
)

interface BridgeDetailsProps extends WithStyles<typeof bridgeDetailsStyles> {
  bridge: models.BridgeType
}

const fields: [string, string][] = [
  ['name', 'Name'],
  ['url', 'URL'],
  ['confirmations', 'Confirmations'],
  ['minimumContractPayment', 'Minimum Contract Payment'],
  ['outgoingToken', 'Outgoing Token'],
]

const BridgeDetails = withStyles(bridgeDetailsStyles)(
  ({ bridge, classes }: BridgeDetailsProps) => (
    <Table>
      <TableBody>
        {fields.map(([k, t]) => {
          return (
            <TableRow key={k} className={classes.tableRow}>
              <TableCell>{t}</TableCell>
              <TableCell>{bridge[k as keyof typeof bridge]}</TableCell>
            </TableRow>
          )
        })}
      </TableBody>
    </Table>
  ),
)

interface RouteParams {
  bridgeId: string
}

interface Props extends WithStyles<typeof styles> {}

export const Show = withStyles(styles)(({ classes }: Props) => {
  const dispatch = useDispatch()
  const history = useHistory()
  const params = useParams<RouteParams>()
  const { bridgeId } = params
  const [bridge, setBridge] = useState<ApiResponse<models.BridgeType>>()
  const [isErrored, setIsErrored] = useState(false)
  const [isDialogOpen, setIsDialogOpen] = useState(false)

  useEffect(() => {
    document.title = 'Show Bridge'
  }, [])

  useEffect(() => {
    async function fetchBridge() {
      try {
        const bt = await v2.bridgeTypes.getBridgeSpec(bridgeId)

        setBridge(bt)
      } catch (e) {
        setIsErrored(true)
      }
    }

    fetchBridge()
  }, [bridgeId])

  const handleDelete = async () => {
    try {
      await v2.bridgeTypes.destroyBridge(bridgeId)

      history.push('/bridges')
      dispatch(notifySuccess(() => <>Bridge deleted</>, {}))
    } catch (e) {
      dispatch(notifyError(ErrorMessage, e))
      setIsDialogOpen(false)
    }
  }

  return (
    <Content>
      <Grid container>
        <Grid item xs={12} md={11} lg={9}>
          <Card>
            <CardHeader
              action={
                bridge && (
                  <>
                    <Button
                      variant="danger"
                      onClick={() => setIsDialogOpen(true)}
                    >
                      Delete
                    </Button>

                    <Button
                      variant="secondary"
                      component={BaseLink}
                      href={`/bridges/${bridge.data.id}/edit`}
                    >
                      Edit
                    </Button>
                  </>
                )
              }
              title="Bridge Info"
            />

            {bridge ? (
              <BridgeDetails bridge={bridge.data.attributes} />
            ) : (
              <Loading isErrored={isErrored} />
            )}

            <Dialog
              open={isDialogOpen}
              onClose={() => setIsDialogOpen(false)}
              maxWidth={false}
            >
              <DialogTitle disableTypography>
                <Typography variant="h5">Delete Bridge</Typography>
              </DialogTitle>
              <DialogContent>
                <Typography>
                  Are you sure you want to delete this bridge?
                </Typography>
              </DialogContent>
              <DialogActions
                classes={{
                  action: classes.action,
                }}
              >
                <Button onClick={() => setIsDialogOpen(false)}>No</Button>
                <Button onClick={handleDelete} variant="danger">
                  Confirm
                </Button>
              </DialogActions>
            </Dialog>
          </Card>
        </Grid>
      </Grid>
    </Content>
  )
})

export default Show
