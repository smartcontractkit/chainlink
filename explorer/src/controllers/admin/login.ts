import { Router, Request, Response } from 'express'
import { Connection } from 'typeorm'
import { Admin } from '../../entity/Admin'
import { getDb } from '../../database'
import { compare as comparePassword } from '../../services/password'

const router = Router()

router.post('/admin/login', async (req: Request, res: Response) => {
  const username = req.header('Explorer-Admin-Username')
  const password = req.header('Explorer-Admin-Password')

  if (username && password) {
    const db = await getDb()
    const admin = await findAdmin(db, username)
    const validPassword = await isValidPassword(password, admin)

    if (validPassword) {
      return res.sendStatus(200)
    }
  }

  return res.sendStatus(401)
})

async function isValidPassword(
  password: string,
  admin?: Admin,
): Promise<boolean> {
  if (!admin) {
    return false
  }

  return comparePassword(password, admin.hashedPassword)
}

function findAdmin(db: Connection, username: string): Promise<Admin> {
  return db.getRepository(Admin).findOne({ username })
}

export default router
