import Adapter from '@wojtekmaj/enzyme-adapter-react-17'
import { configure } from 'enzyme'
import JavascriptTimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
import 'mock-local-storage'
import promiseFinally from 'promise.prototype.finally'

promiseFinally.shim(Promise)
JavascriptTimeAgo.locale(en)

configure({ adapter: new Adapter() })

global.fetch = require('fetch-mock').sandbox()
global.fetch.config.overwriteRoutes = true
