import React, { useState } from 'react'
import { connect } from 'react-redux'
import { useParams, Redirect } from 'react-router-dom'
import {
  createStyles,
  Theme,
  WithStyles,
  withStyles,
} from '@material-ui/core/styles'
import { v2 } from 'api'
import { NodeResource } from '../NodesIndex/NodesIndex'
import { localizedTimestamp, TimeAgo } from 'components/TimeAgo'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import Grid from '@material-ui/core/Grid'
import {
  Table,
  TableBody,
  TableCell,
  TableRow,
  Typography,
} from '@material-ui/core'
import Content from 'components/Content'
import ErrorMessage from 'components/Notifications/DefaultError'
import { deleteNode } from 'actionCreators'
import Button from 'components/Button'
import Close from 'components/Icons/Close'
import Dialog from '@material-ui/core/Dialog'

interface RouteParams {
  nodeId: string
}

const styles = (theme: Theme) =>
  createStyles({
    container: {
      backgroundColor: theme.palette.common.white,
      padding: theme.spacing.unit * 5,
      paddingBottom: 0,
    },
    mainRow: {
      marginBottom: theme.spacing.unit * 2,
    },
    actions: {
      textAlign: 'right',
    },
    button: {
      marginLeft: theme.spacing.unit,
      marginRight: theme.spacing.unit,
    },
    tableHeader: {
      fontWeight: 100,
      paddingRight: 0,
      width: '150px',
    },
    dialogPaper: {
      minHeight: '240px',
      maxHeight: '240px',
      minWidth: '670px',
      maxWidth: '670px',
      overflow: 'hidden',
      borderRadius: theme.spacing.unit * 3,
    },
    warningText: {
      fontWeight: 500,
      marginLeft: theme.spacing.unit * 3,
      marginTop: theme.spacing.unit * 3,
      marginBottom: theme.spacing.unit,
    },
    closeButton: {
      marginRight: theme.spacing.unit * 3,
      marginTop: theme.spacing.unit * 3,
    },
    infoText: {
      fontSize: theme.spacing.unit * 2,
      fontWeight: 450,
      marginLeft: theme.spacing.unit * 6,
    },
    modalContent: {
      width: 'inherit',
    },
    deleteButton: {
      marginTop: theme.spacing.unit * 4,
    },
  })

const DeleteSuccessNotification = ({ id }: any) => (
  <React.Fragment>Successfully deleted node {id}</React.Fragment>
)

interface Props extends WithStyles<typeof styles> {
  deleteNode: Function
}

export const NodesShow = ({ classes, deleteNode }: Props) => {
  const { nodeId } = useParams<RouteParams>()
  const [node, setNode] = React.useState<NodeResource>()
  const [modalOpen, setModalOpen] = useState(false)
  const [deleted, setDeleted] = useState(false)

  const handleDelete = (id: string) => {
    deleteNode(id, () => DeleteSuccessNotification({ id }), ErrorMessage)
    setDeleted(true)
  }

  React.useEffect(() => {
    document.title = `${node?.attributes.name}`
  }, [node])

  React.useEffect(() => {
    Promise.all([v2.nodes.getNodes()])
      .then(([v2Nodes]) =>
        v2Nodes.data.find((node: NodeResource) => node.id === nodeId),
      )
      .then(setNode)
  }, [nodeId])

  return (
    <>
      <Dialog
        open={modalOpen}
        classes={{ paper: classes.dialogPaper }}
        onClose={() => setModalOpen(false)}
      >
        <Grid container spacing={0}>
          <Grid item className={classes.modalContent}>
            <Grid container alignItems="baseline" justify="space-between">
              <Grid item>
                <Typography
                  variant="h5"
                  color="secondary"
                  className={classes.warningText}
                >
                  Warning: This Action Cannot Be Undone
                </Typography>
              </Grid>
              <Grid item>
                <Close
                  className={classes.closeButton}
                  onClick={() => setModalOpen(false)}
                />
              </Grid>
            </Grid>
            <Grid container direction="column">
              <Grid item>
                <Grid item>
                  <Typography
                    className={classes.infoText}
                    variant="h5"
                    color="secondary"
                  >
                    - Access to this page will be lost
                  </Typography>
                </Grid>
              </Grid>
              <Grid container spacing={0} alignItems="center" justify="center">
                <Grid item className={classes.deleteButton}>
                  <Button variant="danger" onClick={() => handleDelete(nodeId)}>
                    Delete {node?.attributes.name}
                    {deleted && <Redirect to="/" />}
                  </Button>
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </Dialog>
      <Card className={classes.container}>
        <Grid container spacing={0}>
          <Grid item xs={12}>
            <Grid
              container
              spacing={0}
              alignItems="center"
              className={classes.mainRow}
            >
              <Grid item xs={6}>
                {node && (
                  <Typography variant="h5" color="secondary">
                    {node?.attributes.name}
                  </Typography>
                )}
              </Grid>
              <Grid item xs={6} className={classes.actions}>
                <Button
                  className={classes.button}
                  onClick={() => setModalOpen(true)}
                >
                  Delete
                </Button>
              </Grid>
            </Grid>
          </Grid>
          <Grid item xs={12}>
            {node?.attributes.createdAt && (
              <Typography variant="subtitle2" color="textSecondary">
                Created{' '}
                <TimeAgo tooltip={false}>{node.attributes.createdAt}</TimeAgo> (
                {localizedTimestamp(node.attributes.createdAt)})
              </Typography>
            )}
          </Grid>
        </Grid>
      </Card>
      <Content>
        <Card>
          <CardContent>
            <Table>
              <TableBody>
                <TableRow>
                  <TableCell className={classes.tableHeader}>Node ID</TableCell>
                  <TableCell>{node?.id}</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell className={classes.tableHeader}>
                    HTTP URL
                  </TableCell>
                  <TableCell>{node?.attributes.httpURL}</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell className={classes.tableHeader}>WS URL</TableCell>
                  <TableCell>{node?.attributes.wsURL}</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell className={classes.tableHeader}>
                    EVM Chain ID
                  </TableCell>
                  <TableCell>{node?.attributes.evmChainID}</TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </Content>
    </>
  )
}

export const ConnectedNodesShow = connect(null, {
  deleteNode,
})(NodesShow)

export default withStyles(styles)(ConnectedNodesShow)
