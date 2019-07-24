import React from 'react'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'

const styles = (theme: Theme) => createStyles({
    animate: {
      animation: 'spin 4s linear infinite'
    },
    '@keyframes spin': {
      '100%': {
        transform: 'rotate(360deg)'
      }
    }
})

interface IProps extends WithStyles<typeof styles> {
  src: string
  width?: number
  height?: number
  spin?: boolean
  alt?: string
}

const Image = ({ src, width, height, alt, classes , spin = false}: IProps) => {
  return (
    <img
      src={src}
      className={spin ? classes.animate : ''}
      alt={alt}
      width={width}
      height={height}
    />
  )
}

export default withStyles(styles)(Image)
