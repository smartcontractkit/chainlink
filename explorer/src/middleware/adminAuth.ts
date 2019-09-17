import { Request, Response, NextFunction } from 'express'
import { Connection } from 'typeorm'
import { Admin, find as findAdmin, isValidPassword } from '../entity/Admin'
import { getDb } from '../database'
import {
  ADMIN_USERNAME_HEADER,
  ADMIN_PASSWORD_HEADER,
} from '../utils/constants'

export default async function(req: Request, res: Response, next: NextFunction) {
  const username: string = req.header(ADMIN_USERNAME_HEADER)
  const password: string = req.header(ADMIN_PASSWORD_HEADER)
  const db: Connection = await getDb()
  const admin: Admin = await findAdmin(db, username)
  const validPassword: boolean = await isValidPassword(password, admin)

  if (!validPassword) {
    res.sendStatus(401)
  }

  next()
}
