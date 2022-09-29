import React, { useState } from 'react'
import { useDispatch } from 'react-redux'
import moment from 'moment'
import {
  deleteErroredJobRuns,
  notifySuccess,
  notifyErrorMsg,
} from 'actionCreators'

import Card from '@material-ui/core/Card'
import CardHeader from '@material-ui/core/CardHeader'
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
import { ConfirmationDialog } from 'src/components/Dialogs/ConfirmationDialog'

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
  const [showErroredConfirm, setErroredConfirm] = useState(false)

  const confirmErroredDeleted = async () => {
    try {
      await dispatch(deleteErroredJobRuns(updatedBefore))
      dispatch(notifySuccess(() => <>Deleted completed job runs</>, {}))

      setErroredConfirm(false)
    } catch (e) {
      dispatch(notifyErrorMsg('Something went wrong'))
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

      <ConfirmationDialog
        open={showErroredConfirm}
        title="Delete errored jobs runs"
        body={`Are you sure you want to delete all errored job runs up to
                ${moment(updatedBefore).format(
                  'dddd, MMMM Do YYYY, h:mm:ss a',
                )}?`}
        confirmButtonText="Confirm"
        onConfirm={confirmErroredDeleted}
        cancelButtonText="Cancel"
        onCancel={() => setErroredConfirm(false)}
      />
    </Card>
  )
})
