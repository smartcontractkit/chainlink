import React from 'react'
import {
  createStyles,
  withStyles,
  WithStyles,
  Theme,
} from '@material-ui/core/styles'
import { Grid } from '@material-ui/core'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import Typography from '@material-ui/core/Typography'

const styles = ({ spacing, palette }: Theme) =>
  createStyles({
    container: {
      height: '100%',
    },
    cardContent: {
      paddingTop: spacing.unit * 6,
      paddingLeft: spacing.unit * 4,
      paddingRight: spacing.unit * 4,
      '&:last-child': {
        paddingBottom: spacing.unit * 6,
      },
    },
    headerRow: {
      textAlign: 'center',
    },
    error: {
      backgroundColor: palette.error.light,
      marginTop: spacing.unit * 2,
    },
    errorText: {
      color: palette.error.main,
    },
  })

interface Props extends WithStyles<typeof styles> {
  onSubmit: () => void
}

export const SignIn = withStyles(styles)(({ classes, onSubmit }: Props) => {
  return (
    <Grid
      container
      justify="center"
      alignItems="center"
      className={classes.container}
      spacing={0}
    >
      <Grid item xs={10} sm={6} md={4} lg={3} xl={2}>
        <Card>
          <CardContent className={classes.cardContent}>
            <form noValidate onSubmit={onSubmit}>
              <Grid container spacing={8}>
                <Grid item xs={12}>
                  <Grid container spacing={0}>
                    <Grid item xs={12} className={classes.headerRow}>
                      <HexagonLogo width={50} />
                    </Grid>
                    <Grid item xs={12} className={classes.headerRow}>
                      <Typography variant="h5">Operator</Typography>
                    </Grid>
                  </Grid>
                </Grid>

                {errors.length > 0 &&
                  errors.map(({ props }, idx) => {
                    return (
                      <Grid item xs={12} key={idx}>
                        <Card raised={false} className={classes.error}>
                          <CardContent>
                            <Typography
                              variant="body1"
                              className={classes.errorText}
                            >
                              {props.msg}
                            </Typography>
                          </CardContent>
                        </Card>
                      </Grid>
                    )
                  })}

                <Grid item xs={12}>
                  <TextField
                    id="email"
                    label="Email"
                    margin="normal"
                    value={email}
                    onChange={handleChange('email')}
                    error={errors.length > 0}
                    variant="outlined"
                    fullWidth
                  />
                </Grid>
                <Grid item xs={12}>
                  <TextField
                    id="password"
                    label="Password"
                    type="password"
                    autoComplete="password"
                    margin="normal"
                    value={password}
                    onChange={handleChange('password')}
                    error={errors.length > 0}
                    variant="outlined"
                    fullWidth
                  />
                </Grid>
                <Grid item xs={12}>
                  <Grid container spacing={0} justify="center">
                    <Grid item>
                      <Button type="submit" variant="primary">
                        Access Account
                      </Button>
                    </Grid>
                  </Grid>
                </Grid>
                {fetching && (
                  <Typography variant="body1" color="textSecondary">
                    Signing in...
                  </Typography>
                )}
              </Grid>
            </form>
          </CardContent>
        </Card>
      </Grid>
    </Grid>
  )
})

export default SignIn
