import { getLogs, formatAnswer } from './utils'

import AggregatorContract from './AggregatorContract'

export default class AggregatorContractV2 extends AggregatorContract {
  async latestAnswer() {
    const latestAnswer = await this.contract.latestAnswer()
    return formatAnswer(
      latestAnswer,
      this.config.multiply,
      this.config.decimalPlaces,
      this.config.formatDecimalPlaces,
    )
  }

  async latestAnswerTimestamp() {
    const latestTimestamp = await this.contract.latestTimestamp()
    return latestTimestamp.toNumber()
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
        answerFormatted: formatAnswer(
          decodedLog.current,
          this.config.multiply,
          this.config.decimalPlaces,
          this.config.formatDecimalPlaces,
        ),
        answer: Number(decodedLog.current),
        answerId: Number(decodedLog.roundId),
        timestamp: Number(decodedLog.timestamp),
      }),
    )

    return logs
  }
}
