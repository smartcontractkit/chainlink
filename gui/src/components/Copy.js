import React from 'react'
import Button from '@material-ui/core/Button'
import { CopyToClipboard } from 'react-copy-to-clipboard'

const Copy = ({data, buttonText}) => (
  <CopyToClipboard text={data}>
    <Button color='primary' variant='outlined'>{buttonText}</Button>
  </CopyToClipboard>
)

export default Copy
