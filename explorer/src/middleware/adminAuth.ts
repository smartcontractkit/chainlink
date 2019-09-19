import { Request, Response, NextFunction } from 'express'
import httpStatus from 'http-status-codes'
import { find as findAdmin, isValidPassword } from '../entity/Admin'
import { getDb } from '../database'
import {
  ADMIN_USERNAME_HEADER,
  ADMIN_PASSWORD_HEADER,
} from '../utils/constants'

export default async function(req: Request, res: Response, next: NextFunction) {
  if (!req.session.admin) {
    const username: string = req.header(ADMIN_USERNAME_HEADER)
    const password: string = req.header(ADMIN_PASSWORD_HEADER)

    if (!username || !password) {
      res.sendStatus(httpStatus.UNAUTHORIZED)
      return
    }

    const db = await getDb()
    const admin = await findAdmin(db, username)
    const validPassword: boolean = await isValidPassword(password, admin)

    // Express.js middleware is sequential
    if (validPassword) {
      /* eslint-disable-next-line require-atomic-updates */
      req.session.admin = admin
    } else {
      /* eslint-disable-next-line require-atomic-updates */
      req.session.admin = null
      res.sendStatus(httpStatus.UNAUTHORIZED)
    }
  }

  next()
}
