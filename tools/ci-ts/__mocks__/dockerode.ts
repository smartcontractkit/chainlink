interface ContainerMock {
  pause: typeof jest.fn
  unpause: typeof jest.fn
  inspect: typeof jest.fn
}

interface DockerMock {
  getContainer: typeof jest.fn
}

interface DockerodeMock extends jest.Mock {
  setContainerState: (newState: object) => void
  container: ContainerMock
  docker: DockerMock
}

let containerState: object = {
  Paused: false,
}

function setContainerState(newState: object) {
  containerState = newState
}

const container = {
  pause: jest.fn(),
  unpause: jest.fn(),
  inspect: jest
    .fn()
    .mockImplementation(async () => ({ State: containerState })),
}

const docker = {
  getContainer: jest.fn().mockReturnValue(container),
}

const generator = () => jest.fn().mockReturnValue(docker)
const mock = (jest.mock('dockerode', generator) as unknown) as DockerodeMock

mock.setContainerState = setContainerState
mock.container = container
mock.docker = docker
export default mock
