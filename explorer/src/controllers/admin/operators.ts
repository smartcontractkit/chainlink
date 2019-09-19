import { Router, Request, Response } from 'express'

const router = Router()

router.post('/operators', async (req: Request, res: Response) => {
  return res.status(200).send([])
})

export default router
