export default class DockerodeMock {
  static containerState: object = {
    Paused: false,
  }
  static setContainerState(newState: object) {
    DockerodeMock.containerState = newState
  }
  static container = {
    pause: jest.fn(),
    unpause: jest.fn(),
    inspect: jest.fn().mockImplementation(async () => ({
      State: DockerodeMock.containerState,
    })),
  }
  static getContainer = jest.fn().mockReturnValue(DockerodeMock.container)

  getContainer = DockerodeMock.getContainer
}
