import React from 'react'
import Grid from '@material-ui/core/Grid'
import Card from '@material-ui/core/Card'
import Button from '@material-ui/core/Button'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import Logo from '../Logo'
import SearchBox from '../SearchBox'
import SearchForm from '../SearchForm'

const styles = ({ spacing }: Theme) =>
  createStyles({
    container: {
      height: '100vh'
    },
    card: {
      paddingTop: spacing.unit * 5,
      paddingBottom: spacing.unit * 5,
      paddingLeft: spacing.unit * 8,
      paddingRight: spacing.unit * 8
    },
    logo: {
      display: 'flex'
    }
  })

interface IProps extends WithStyles<typeof styles> {
  path: string
}

const Search = ({ classes }: IProps) => {
  return (
    <Grid
      container
      justify="center"
      alignItems="center"
      className={classes.container}>
      <Grid item md={8} lg={6} xl={4}>
        <Card className={classes.card}>
          <SearchForm>
            <Grid container justify="center">
              <Grid item>
                <Logo className={classes.logo} width={300} height={80} />
              </Grid>
              <Grid item xs={12}>
                <SearchBox />
              </Grid>
              <Grid item>
                <Button variant="contained" color="primary" type="submit">
                  Search
                </Button>
              </Grid>
            </Grid>
          </SearchForm>
        </Card>
      </Grid>
    </Grid>
  )
}

export default withStyles(styles)(Search)
