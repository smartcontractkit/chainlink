import { createStyles, Theme } from '@material-ui/core/styles'

// Contains styles to make a table row linkable
//
// This has been placed in the components folder because we ideally want to
// wrap the link in a reusable components, however using Material UI 3
// 'withStyles' leads to small freeze when navigating to the link, due to
// the slowness of JSS. When we migrate to MUI 4/5 we can implement this with
// makeStyles which doesn't have those problems
export const tableStyles = (theme: Theme) =>
  createStyles({
    cell: {
      paddingTop: theme.spacing.unit * 2,
      paddingBottom: theme.spacing.unit * 2,
    },
    row: {
      transform: 'scale(1)',
    },
    link: {
      '&::before': {
        content: "''",
        position: 'absolute',
        top: 0,
        left: 0,
        width: '100%',
        height: '100%',
      },
    },
  })
