import React from 'react'
import { createStyles, withStyles, WithStyles } from '@material-ui/core/styles'

const styles = createStyles({
  animate: {
    animation: 'spin 4s linear infinite',
  },
  '@keyframes spin': {
    '100%': {
      transform: 'rotate(360deg)',
    },
  },
})

interface Props extends WithStyles<typeof styles> {
  src: string
  width?: number
  height?: number
  spin?: boolean
  alt?: string
}

const UnstyledImage = ({
  src,
  width,
  height,
  alt,
  classes,
  spin = false,
}: Props) => {
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

const Image = withStyles(styles)(UnstyledImage)

interface Props {
  src: string
  width?: number
  height?: number
  alt?: string
}

export const Logo: React.FC<Props> = (props) => {
  return <Image {...props} />
}
