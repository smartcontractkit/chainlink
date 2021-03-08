import { createStyles, Theme } from '@material-ui/core/styles'

export const styles = (theme: Theme) =>
  createStyles({
    cell: {
      paddingTop: theme.spacing.unit * 2,
      paddingBottom: theme.spacing.unit * 2,
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
