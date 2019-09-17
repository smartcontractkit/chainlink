import { Router, Request, Response } from 'express'
import { Connection } from 'typeorm'
import { validate } from 'class-validator'
import { getDb } from '../../database'
import { buildChainlinkNode, ChainlinkNode } from '../../entity/ChainlinkNode'
import { PG_UNIQUE_CONSTRAINT_VIOLATION } from '../../utils/constants'

const router = Router()

router.post('/nodes', async (req: Request, res: Response) => {
  const name: string = req.body.name
  const url: string = req.body.url
  const db: Connection = await getDb()
  const [node, secret] = await buildChainlinkNode(db, name, url)
  const errors = await validate(node)

  if (errors.length === 0) {
    try {
      const savedNode = await db.manager.save(node)

      return res.status(201).json({
        id: savedNode.id,
        accessKey: savedNode.accessKey,
        secret: secret,
      })
    } catch (e) {
      if (e.code === PG_UNIQUE_CONSTRAINT_VIOLATION) {
        return res.sendStatus(409)
      }

      console.error(e)
      return res.sendStatus(400)
    }
  }

  const jsonApiErrors = errors.reduce((acc, e) => {
    return { ...acc, [e.property]: e.constraints }
  }, {})

  return res.status(422).send({ errors: jsonApiErrors })
})

router.delete('/nodes/:name', async (req: Request, res: Response) => {
  const db: Connection = await getDb()

  await db.getRepository(ChainlinkNode).delete({ name: req.params.name })

  return res.sendStatus(200)
})

export default router
