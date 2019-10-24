import normalize from 'json-api-normalizer'
import * as api from '../api/index'
import { request } from './helpers'

export const signIn = request(
  'ADMIN_SIGNIN',
  api.v1.adminAuth.signIn,
  normalize,
)

export const signOut = request(
  'ADMIN_SIGNOUT',
  api.v1.adminAuth.signOut,
  normalize,
)
