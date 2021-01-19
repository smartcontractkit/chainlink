import { Api } from 'utils/json-api-client'
import { Sessions } from './sessions'
import { V2 } from './v2'

const api = new Api({
  base: process.env.CHAINLINK_BASEURL,
})

export const sessions = new Sessions(api)
export const v2 = new V2(api)
