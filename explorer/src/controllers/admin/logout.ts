import { Router } from 'express'
import httpStatus from 'http-status-codes'

const router = Router()

router.delete('/logout', (req, res) => {
  req.session.admin = null
  return res.sendStatus(httpStatus.NO_CONTENT)
})

export default router
