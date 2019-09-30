import { Request, Response, NextFunction } from 'express'
import httpStatus from 'http-status-codes'
import { find as findAdmin, isValidPassword } from '../entity/Admin'
import { getDb } from '../database'
import {
  ADMIN_USERNAME_HEADER,
  ADMIN_PASSWORD_HEADER,
} from '../utils/constants'

export default async function(req: Request, res: Response, next: NextFunction) {
  const username = req.header(ADMIN_USERNAME_HEADER)
  const password = req.header(ADMIN_PASSWORD_HEADER)
  const db = await getDb()
  const admin = await findAdmin(db, username)
  const validPassword = await isValidPassword(password, admin)

  if (!validPassword) {
    res.sendStatus(httpStatus.UNAUTHORIZED)
  }

  next()
}
