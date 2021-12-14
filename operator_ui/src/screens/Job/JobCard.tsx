import React from 'react'

import DeleteIcon from '@material-ui/icons/Delete'
import FileCopyIcon from '@material-ui/icons/FileCopy'
import IconButton from '@material-ui/core/IconButton'
import Grid from '@material-ui/core/Grid'
import ListItemIcon from '@material-ui/core/ListItemIcon'
import ListItemText from '@material-ui/core/ListItemText'
import Menu from '@material-ui/core/Menu'
import MenuItem from '@material-ui/core/MenuItem'
import MoreVertIcon from '@material-ui/icons/MoreVert'

import {
  DetailsCard,
  DetailsCardItemTitle,
  DetailsCardItemValue,
} from 'src/components/Cards/DetailsCard'
import { formatJobSpecType } from 'src/utils/formatJobSpecType'
import { generateJobDefinition } from './generateJobDefinition'
import { TimeAgo } from 'src/components/TimeAgo'
import { MenuItemLink } from 'src/components/MenuItemLink'
import { ConfirmationDialog } from 'src/components/Dialogs/ConfirmationDialog'

interface Props {
  job: JobPayload_Fields
  onDelete: () => void
}

export const JobCard: React.FC<Props> = ({ job, onDelete }) => {
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null)
  const [deleteDialogOpen, setDeleteDialogOpen] = React.useState(false)

  const handleMenuOpen = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleMenuClose = () => {
    setAnchorEl(null)
  }

  const { definition } = generateJobDefinition(job)

  return (
    <>
      <DetailsCard
        actions={
          <div>
            <IconButton onClick={handleMenuOpen} aria-label="open-menu">
              <MoreVertIcon />
            </IconButton>
            <Menu
              anchorEl={anchorEl}
              open={Boolean(anchorEl)}
              onClose={handleMenuClose}
              disableAutoFocusItem
            >
              <MenuItemLink
                to={`/jobs/new?definition=${encodeURIComponent(definition)}`}
              >
                <ListItemIcon>
                  <FileCopyIcon />
                </ListItemIcon>
                <ListItemText>Duplicate</ListItemText>
              </MenuItemLink>
              <MenuItem onClick={() => setDeleteDialogOpen(true)}>
                <ListItemIcon>
                  <DeleteIcon />
                </ListItemIcon>
                <ListItemText>Delete</ListItemText>
              </MenuItem>
            </Menu>
          </div>
        }
      >
        <Grid container>
          <Grid item xs={12} sm={6} md={1}>
            <DetailsCardItemTitle title="ID" />
            <DetailsCardItemValue value={job.id} />
          </Grid>
          <Grid item xs={12} sm={6} md={2}>
            <DetailsCardItemTitle title="Type" />
            <DetailsCardItemValue
              value={formatJobSpecType(job.spec.__typename)}
            />
          </Grid>
          <Grid item xs={12} sm={6} md={5}>
            <DetailsCardItemTitle title="External Job ID" />
            <DetailsCardItemValue
              value={formatJobSpecType(job.externalJobID)}
            />
          </Grid>
          <Grid item xs={12} sm={6} md={4}>
            <DetailsCardItemTitle title="Created" />
            <DetailsCardItemValue>
              <TimeAgo tooltip>{job.createdAt}</TimeAgo>
            </DetailsCardItemValue>
          </Grid>
        </Grid>
      </DetailsCard>

      <ConfirmationDialog
        open={deleteDialogOpen}
        title="Delete Job?"
        body="Warning: This action cannot be undone!"
        confirmButtonText="Confirm"
        onConfirm={() => {
          onDelete()
          setDeleteDialogOpen(false)
        }}
        cancelButtonText="Cancel"
        onCancel={() => setDeleteDialogOpen(false)}
      />
    </>
  )
}
