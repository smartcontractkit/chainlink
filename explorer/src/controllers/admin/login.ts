import { Router } from 'express'

const router = Router()

router.post('/login', async (_req, res) => {
  return res.sendStatus(200)
})

export default router
