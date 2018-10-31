import { configure } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import 'mock-local-storage'

// Human readable date & time
import JavascriptTimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'

JavascriptTimeAgo.locale(en)

configure({ adapter: new Adapter() })

global.fetch = require('fetch-mock').sandbox()
global.fetch.config.overwriteRoutes = true
