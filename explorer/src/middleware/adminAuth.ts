import { Request, Response, NextFunction } from 'express'
import httpStatus from 'http-status-codes'
import { find as findAdmin, isValidPassword } from '../entity/Admin'
import {
  ADMIN_USERNAME_HEADER,
  ADMIN_PASSWORD_HEADER,
  ADMIN_USERNAME_PARAM,
  ADMIN_PASSWORD_PARAM,
} from '../utils/constants'

export default async function(req: Request, res: Response, next: NextFunction) {
  if (!req.session.admin) {
    const username =
      req.body[ADMIN_USERNAME_PARAM] || req.header(ADMIN_USERNAME_HEADER)
    const password =
      req.body[ADMIN_PASSWORD_PARAM] || req.header(ADMIN_PASSWORD_HEADER)

    if (!username || !password) {
      res.sendStatus(httpStatus.UNAUTHORIZED)
      return
    }

    const admin = await findAdmin(username)
    const validPassword = await isValidPassword(password, admin)

    // Express.js middleware is sequential
    if (validPassword) {
      /* eslint-disable-next-line require-atomic-updates */
      req.session.admin = admin
    } else {
      /* eslint-disable-next-line require-atomic-updates */
      req.session.admin = null
      res.sendStatus(httpStatus.UNAUTHORIZED)
      return
    }
  }

  next()
}
