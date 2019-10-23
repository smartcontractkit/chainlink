import { Router } from 'express'

const router = Router()

router.delete('/logout', (req, res) => {
  req.session.admin = null
  return res.sendStatus(200)
})

export default router
