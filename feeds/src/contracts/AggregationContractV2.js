import { getLogs, formatAnswer } from './utils'

import AggregationContract from './AggregationContract'

export default class AggregationContractV2 extends AggregationContract {
  async currentAnswer() {
    const latestAnswer = await this.contract.latestAnswer()
    return formatAnswer(
      latestAnswer,
      this.options.multiply,
      this.options.decimalPlaces,
    )
  }

  async updateHeight() {
    const latestTimestamp = await this.contract.latestTimestamp()
    return Number(latestTimestamp)
  }

  async latestCompletedAnswer() {
    const latestRound = await this.contract.latestRound()
    return Number(latestRound)
  }

  async answerUpdatedLogs({ fromBlock }) {
    const answerUpdatedFilter = {
      ...this.contract.filters.AnswerUpdated(null, null),
      fromBlock,
      toBlock: 'latest',
    }

    const logs = await getLogs(
      {
        provider: this.provider,
        filter: answerUpdatedFilter,
        eventInterface: this.contract.interface.events.AnswerUpdated,
      },
      decodedLog => ({
        responseFormatted: formatAnswer(
          decodedLog.current,
          this.options.multiply,
          this.options.decimalPlaces,
        ),
        response: Number(decodedLog.current),
        answerId: Number(decodedLog.roundId),
        timestamp: Number(decodedLog.timestamp),
      }),
    )

    return logs
  }
}
