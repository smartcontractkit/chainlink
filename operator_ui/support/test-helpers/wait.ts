import { screen, waitForElementToBeRemoved } from 'support/test-utils'

export const waitForLoading = () => {
  return waitForElementToBeRemoved(() => screen.queryByRole('progressbar'))
}
