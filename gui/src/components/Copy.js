import React from 'react'
import Button from '@material-ui/core/Button'
import { CopyToClipboard } from 'react-copy-to-clipboard'

const Copy = ({data, buttonText}) => {
  return (
    <div>
      <CopyToClipboard text={data}>
        <Button color='secondary' variant='outlined'>{buttonText}</Button>
      </CopyToClipboard>
    </div>
  )
}

export default Copy
