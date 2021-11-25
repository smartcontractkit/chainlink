import React from 'react'

import DeleteIcon from '@material-ui/icons/Delete'
import Grid from '@material-ui/core/Grid'
import IconButton from '@material-ui/core/IconButton'
import ListItemIcon from '@material-ui/core/ListItemIcon'
import ListItemText from '@material-ui/core/ListItemText'
import Menu from '@material-ui/core/Menu'
import MenuItem from '@material-ui/core/MenuItem'
import Typography from '@material-ui/core/Typography'
import MoreVertIcon from '@material-ui/icons/MoreVert'

import { DetailsCard } from 'src/components/Cards/DetailsCard'
import { TimeAgo } from 'src/components/TimeAgo'

interface Props {
  node: NodePayload_Fields
  onDelete: () => void
}

export const NodeCard: React.FC<Props> = ({ node, onDelete }) => {
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null)

  const handleOpen = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleClose = () => {
    setAnchorEl(null)
  }

  const onDeleteClick = () => {
    onDelete()
    setAnchorEl(null)
  }

  return (
    <DetailsCard
      actions={
        <div>
          <IconButton onClick={handleOpen} aria-label="open-menu">
            <MoreVertIcon />
          </IconButton>
          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleClose}
          >
            <MenuItem onClick={onDeleteClick}>
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
        <Grid item xs={12} sm={4} md={3}>
          <Typography variant="subtitle2" gutterBottom>
            ID
          </Typography>
          <Typography variant="body1">{node.id}</Typography>
        </Grid>

        <Grid item xs={12} sm={4} md={3}>
          <Typography variant="subtitle2" gutterBottom>
            EVM Chain ID
          </Typography>
          <Typography variant="body1">{node.chain.id}</Typography>
        </Grid>

        <Grid item xs={12} sm={4} md={2}>
          <Typography variant="subtitle2" gutterBottom>
            Created
          </Typography>
          <Typography variant="body1">
            <TimeAgo tooltip>{node.createdAt}</TimeAgo>
          </Typography>
        </Grid>

        <Grid item xs={false} sm={false} md={4}></Grid>

        <Grid item xs={12} md={6}>
          <Typography variant="subtitle2" gutterBottom>
            HTTP URL
          </Typography>
          <Typography variant="body1" noWrap>
            {node.httpURL !== '' ? node.httpURL : '--'}
          </Typography>
        </Grid>

        <Grid item xs={12} md={6}>
          <Typography variant="subtitle2" gutterBottom>
            WS URL
          </Typography>
          <Typography variant="body1" noWrap>
            {node.wsURL !== '' ? node.wsURL : '--'}
          </Typography>
        </Grid>
      </Grid>
    </DetailsCard>
  )
}
