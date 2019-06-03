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
import TermsOfUse from '../TermsOfUse'
import Logo from '../Logo'
import SearchBox from '../SearchBox'
import SearchForm from '../SearchForm'

const styles = ({ spacing }: Theme) =>
  createStyles({
    container: {
      display: 'flex',
      alignItems: 'center',
      width: '100%'
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
    <div className={classes.container}>
      <Grid container justify="center" alignItems="center">
        <Grid item md={8} lg={6} xl={4}>
          <Grid container>
            <Grid item xs={12}>
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
            <Grid item xs={12}>
              <TermsOfUse />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </div>
  )
}

export default withStyles(styles)(Search)
