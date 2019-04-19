import { clearDb } from './src/__tests__/testdatabase'

process.env.NODE_ENV = 'test'

beforeEach(async () => clearDb())
afterEach(async () => clearDb())
