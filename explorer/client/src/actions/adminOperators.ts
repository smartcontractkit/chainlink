import normalize from 'json-api-normalizer'
import { request } from './helpers'
import * as api from '../api/index'

export const fetchAdminOperators = request(
  'ADMIN_OPERATORS',
  api.v1.adminOperators.getOperators,
  json => normalize(json, { endpoint: 'currentPageOperators' }),
)
