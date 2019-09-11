import { Router, Request, Response } from 'express'
import { Connection } from 'typeorm'
import { Admin } from '../../entity/Admin'
import { getDb } from '../../database'
import { compare as comparePassword } from '../../services/password'

const router = Router()

router.post('/admin/login', async (req: Request, res: Response) => {
  const username: string = req.header('Explorer-Admin-Username')
  const password: string = req.header('Explorer-Admin-Password')
  const db: Connection = await getDb()
  const admin: Admin = await findAdmin(db, username)
  const validPassword: boolean = await isValidPassword(password, admin)

  if (validPassword) {
    return res.sendStatus(200)
  }

  return res.sendStatus(401)
})

async function isValidPassword(
  password: string,
  admin?: Admin,
): Promise<boolean> {
  if (!admin) {
    return new Promise(resolve => resolve(false))
  }

  return comparePassword(password, admin.hashedPassword)
}

function findAdmin(db: Connection, username: string): Promise<Admin> {
  return db.getRepository(Admin).findOne({ username: username })
}

export default router
