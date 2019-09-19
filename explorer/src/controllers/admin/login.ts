import { Router, Request, Response } from 'express'

const router = Router()

router.post('/login', async (req: Request, res: Response) => {
  return res.sendStatus(200)
})

export default router
