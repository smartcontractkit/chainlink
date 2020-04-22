import { configure } from 'enzyme'
import 'mock-local-storage'
import Adapter from 'enzyme-adapter-react-16'
import JavascriptTimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'

configure({ adapter: new Adapter() })

JavascriptTimeAgo.locale(en)
