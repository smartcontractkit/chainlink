import React, { useState } from 'react'
import { useDispatch } from 'react-redux'
import moment from 'moment'
import {
  deleteCompletedJobRuns,
  deleteErroredJobRuns,
  notifySuccess,
  notifyError,
} from 'actionCreators'

import Button from '@material-ui/core/Button'
import Card from '@material-ui/core/Card'
import CardHeader from '@material-ui/core/CardHeader'
import Dialog from '@material-ui/core/Dialog'
import DialogActions from '@material-ui/core/DialogActions'
import DialogContent from '@material-ui/core/DialogContent'
import DialogContentText from '@material-ui/core/DialogContentText'
import DialogTitle from '@material-ui/core/DialogTitle'
import IconButton from '@material-ui/core/IconButton'
import DeleteIcon from '@material-ui/icons/Delete'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'

const styles = (theme: Theme) => {
  return createStyles({
    confirmDialog: {
      width: 400,
    },
    deleteCell: {
      '&:last-child': {
        paddingRight: theme.spacing.unit * 1.5,
      },
    },
    deleteIcon: {
      fontSize: 20,
    },
  })
}

const WEEK_MS = 1000 * 60 * 60 * 24 * 7

interface Props extends WithStyles<typeof styles> {}

export const JobRuns = withStyles(styles)(({ classes }: Props) => {
  const dispatch = useDispatch()

  const updatedBefore = new Date(Date.now() - WEEK_MS).toISOString()
  const [showCompletedConfirm, setCompletedConfirm] = useState(false)
  const [showErroredConfirm, setErroredConfirm] = useState(false)

  const confirmCompletedDelete = async () => {
    try {
      await dispatch(deleteCompletedJobRuns(updatedBefore))
      dispatch(notifySuccess(() => <>Deleted completed job runs</>, {}))

      setCompletedConfirm(false)
    } catch (e) {
      dispatch(notifyError(() => <>Something went wrong</>, e))
      console.log('some error occurred')
    }
  }
  const confirmErroredDeleted = async () => {
    try {
      await dispatch(deleteErroredJobRuns(updatedBefore))
      dispatch(notifySuccess(() => <>Deleted completed job runs</>, {}))

      setErroredConfirm(false)
    } catch (e) {
      dispatch(notifyError(() => <>Something went wrong</>, e))
      console.log('some error occurred')
    }
  }

  return (
    <Card>
      <CardHeader
        title="Job Runs"
        subheader="Reduce the database size by deleting old job runs"
      />

      <Table>
        <TableBody>
          <TableRow>
            <TableCell>
              <Typography>Completed Runs</Typography>
              <Typography variant="subtitle2" color="textSecondary">
                Keeps runs from the last week
              </Typography>
            </TableCell>
            <TableCell align="right" className={classes.deleteCell}>
              <IconButton
                onClick={() => setCompletedConfirm(true)}
                data-cy="delete-completed-job-runs"
              >
                <DeleteIcon className={classes.deleteIcon} color="error" />
              </IconButton>
            </TableCell>
          </TableRow>

          <TableRow>
            <TableCell>
              <Typography>Errored Runs</Typography>
              <Typography variant="subtitle2" color="textSecondary">
                Keeps runs from the last week
              </Typography>
            </TableCell>
            <TableCell align="right" className={classes.deleteCell}>
              <IconButton
                onClick={() => setErroredConfirm(true)}
                data-cy="delete-errored-job-runs"
              >
                <DeleteIcon className={classes.deleteIcon} color="error" />
              </IconButton>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>

      <Dialog
        open={showCompletedConfirm}
        onClose={() => setCompletedConfirm(false)}
        PaperProps={{ className: classes.confirmDialog }}
      >
        <DialogTitle>
          <Typography variant="h5"> Delete completed jobs runs</Typography>
        </DialogTitle>
        <DialogContent>
          <DialogContentText>
            Are you sure you want to delete all completed job runs up to{' '}
            {moment(updatedBefore).format('dddd, MMMM Do YYYY, h:mm:ss a')}?
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCompletedConfirm(false)} color="secondary">
            Cancel
          </Button>
          <Button
            variant="contained"
            onClick={confirmCompletedDelete}
            color="primary"
            autoFocus
          >
            Confirm
          </Button>
        </DialogActions>
      </Dialog>

      <Dialog
        open={showErroredConfirm}
        onClose={() => setErroredConfirm(false)}
        PaperProps={{ className: classes.confirmDialog }}
      >
        <DialogTitle>
          <Typography variant="h5"> Delete completed jobs runs</Typography>
        </DialogTitle>
        <DialogContent>
          <DialogContentText>
            Are you sure you want to delete all errored job runs up to{' '}
            {moment(updatedBefore).format('dddd, MMMM Do YYYY, h:mm:ss a')}?
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setErroredConfirm(false)} color="secondary">
            Cancel
          </Button>
          <Button
            variant="contained"
            onClick={confirmErroredDeleted}
            color="primary"
            autoFocus
          >
            Confirm
          </Button>
        </DialogActions>
      </Dialog>
    </Card>
  )
})
