import normalize from 'json-api-normalizer'
import { request } from './helpers'
import * as api from '../api/index'

export const fetchAdminHeads = request(
  'ADMIN_HEADS',
  api.v1.adminHeads.getHeads,
  json => normalize(json, { endpoint: 'currentPageHeads' }),
)

export const fetchAdminHead = request(
  'ADMIN_HEAD',
  api.v1.adminHeads.getHead,
  json => normalize(json, { endpoint: 'head' }),
)
