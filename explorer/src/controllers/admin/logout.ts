import { Router, Request, Response } from 'express'

const router = Router()

router.delete('/logout', async (req: Request, res: Response) => {
  req.session.admin = null
  return res.sendStatus(200)
})

export default router
