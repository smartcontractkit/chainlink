import { Router } from 'express'

const router = Router()

router.post('/operators', async (_req, res) => {
  return res.status(200).send([])
})

export default router
