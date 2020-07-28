import { Router, Response } from 'express'
import { Config } from '../config'

const router = Router()

router.get('/', async (_, res: Response) => {
  return res.json({
    gaId: Config.gaId(),
  })
})

export default router
