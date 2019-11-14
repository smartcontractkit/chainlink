import { Router, Request } from 'express'
import { Connection } from 'typeorm'
import { validate } from 'class-validator'
import httpStatus from 'http-status-codes'
import { getDb } from '../../database'
import { buildChainlinkNode, ChainlinkNode } from '../../entity/ChainlinkNode'
import { PostgresErrorCode } from '../../utils/constants'
import { isPostgresError } from '../../utils/errors'
import {
  DEFAULT_PAGE,
  DEFAULT_PAGE_SIZE,
  PaginationParams,
} from '../../utils/pagination'
import { getCustomRepository } from 'typeorm'
import { ChainlinkNodeRepository } from '../../repositories/ChainlinkNodeRepository'
import chainlinkNodesSerializer from '../../serializers/chainlinkNodesSerializer'

const router = Router()

const parseParams = (req: Request): PaginationParams => {
  const page = parseInt(req.query.page, 10) || DEFAULT_PAGE
  const size = parseInt(req.query.size, 10) || DEFAULT_PAGE_SIZE

  return {
    page: page,
    limit: size,
  }
}

router.get('/nodes', async (req, res) => {
  const params = parseParams(req)
  const db = await getDb()
  const chainlinkNodeRepository = getCustomRepository(
    ChainlinkNodeRepository,
    db.name,
  )
  const chainlinkNodes = await chainlinkNodeRepository.all(params)
  const nodeCount = await chainlinkNodeRepository.count()
  const json = chainlinkNodesSerializer(chainlinkNodes, nodeCount)

  return res.send(json)
})

router.post('/nodes', async (req, res) => {
  const name = req.body.name
  const url = req.body.url
  const db = await getDb()
  const [node, secret] = buildChainlinkNode(db, name, url)
  const errors = await validate(node)

  if (errors.length === 0) {
    try {
      const savedNode = await db.manager.save(node)

      return res.status(httpStatus.CREATED).json({
        id: savedNode.id,
        accessKey: savedNode.accessKey,
        secret: secret,
      })
    } catch (e) {
      if (
        isPostgresError(e) &&
        e.code === PostgresErrorCode.UNIQUE_CONSTRAINT_VIOLATION
      ) {
        return res.sendStatus(httpStatus.CONFLICT)
      }

      console.error(e)
      return res.sendStatus(httpStatus.BAD_REQUEST)
    }
  }

  const jsonApiErrors = errors.reduce(
    (acc, e) => ({ ...acc, [e.property]: e.constraints }),
    {},
  )

  return res
    .status(httpStatus.UNPROCESSABLE_ENTITY)
    .send({ errors: jsonApiErrors })
})

router.delete('/nodes/:name', async (req, res) => {
  const db: Connection = await getDb()

  await db.getRepository(ChainlinkNode).delete({ name: req.params.name })

  return res.sendStatus(httpStatus.OK)
})

export default router
