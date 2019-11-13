import { Router } from 'express'

const router = Router()

router.post('/operators', (_req, res) => {
  return res.status(200).send([])
})

export default router
