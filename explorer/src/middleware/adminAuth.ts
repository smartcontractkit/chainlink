import { Request, Response, NextFunction } from 'express'
import { Connection } from 'typeorm'
import { Admin, find as findAdmin, isValidPassword } from '../entity/Admin'
import { getDb } from '../database'

export default async function(req: Request, res: Response, next: NextFunction) {
  const username: string = req.header('Explorer-Admin-Username')
  const password: string = req.header('Explorer-Admin-Password')
  const db: Connection = await getDb()
  const admin: Admin = await findAdmin(db, username)
  const validPassword: boolean = await isValidPassword(password, admin)

  if (!validPassword) {
    res.sendStatus(401)
  }

  next()
}
