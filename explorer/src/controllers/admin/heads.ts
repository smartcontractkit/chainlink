import { Router } from 'express'
import { EthereumHeadRepository } from '../../repositories/EthereumHeadRepository'
import { Head } from '../../entity/Head'
import { getCustomRepository, getRepository } from 'typeorm'
import { parseParams } from '../../utils/pagination'
import headsSerializer from '../../serializers/headsSerializer'
import headSerializer from '../../serializers/headSerializer'

const router = Router()

router.get('/heads', async (req, res) => {
  const params = parseParams(req.query)
  const ethereumHeadRepository = getCustomRepository(EthereumHeadRepository)

  const heads = await ethereumHeadRepository.all(params)

  const headCount = await ethereumHeadRepository.count()
  const json = headsSerializer(heads, headCount)
  return res.send(json)
})

router.get('/heads/:headId', async (req, res) => {
  const ethereumHeadRepository = getRepository(Head)
  const headId = req.params.headId
  const head = await ethereumHeadRepository.findOne(headId)

  const json = headSerializer(head)
  return res.send(json)
})

export default router
