import { Router } from 'express'
import { getDb } from '../../database'
import { EthereumHeadRepository } from '../../repositories/EthereumHeadRepository'
import { Head } from '../../entity/Head'
import { getCustomRepository } from 'typeorm'
import { parseParams } from '../../utils/pagination'
import headsSerializer from '../../serializers/headsSerializer'
import headSerializer from '../../serializers/headSerializer'

const router = Router()

router.get('/heads', async (req, res) => {
  const params = parseParams(req.query)
  const db = await getDb()
  const ethereumHeadRepository = getCustomRepository(
    EthereumHeadRepository,
    db.name,
  )

  const heads = await ethereumHeadRepository.all(params)

  const headCount = await ethereumHeadRepository.count()
  const json = headsSerializer(heads, headCount)
  return res.send(json)
})

router.get('/heads/:headId', async (req, res) => {
  const db = await getDb()
  const ethereumHeadRepository = db.getRepository(Head)

  const { id } = req.params
  const head = await ethereumHeadRepository.findOne(id)

  const json = headSerializer(head)
  return res.send(json)
})

export default router
