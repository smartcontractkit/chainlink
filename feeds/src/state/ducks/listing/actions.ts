export interface SetAnswersAction {
  type: 'listing/SET_ANSWERS'
  payload: any
}

export function setAnswers(payload: any): SetAnswersAction {
  return {
    type: 'listing/SET_ANSWERS',
    payload,
  }
}
