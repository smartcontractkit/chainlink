"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const tslib_1 = require("tslib");
const path_1 = require("path");
const shelljs_1 = require("shelljs");
const compile_1 = tslib_1.__importDefault(require("../../../src/commands/compile"));
const TEST_PATH = 'test/functional/truffle/';
const TEST_FS_PATH = path_1.join(TEST_PATH, 'testfs');
const FIXTURES_PATH = path_1.join(TEST_PATH, 'fixtures');
describe('compileAll', () => {
    beforeEach(() => {
        shelljs_1.rm('-r', TEST_FS_PATH);
    });
    it('should produce a truffle contract abstraction via a json artifact produced by sol-compiler', async () => {
        await compile_1.default.run([
            `--config=${path_1.join(FIXTURES_PATH, 'test.config.json')}`,
            'truffle',
        ]);
        expect([...shelljs_1.ls('-R', TEST_FS_PATH)]).toMatchSnapshot();
    });
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidHJ1ZmZsZS50ZXN0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vdGVzdC9mdW5jdGlvbmFsL3RydWZmbGUvdHJ1ZmZsZS50ZXN0LnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7OztBQUFBLCtCQUEyQjtBQUMzQixxQ0FBZ0M7QUFDaEMsb0ZBQW1EO0FBRW5ELE1BQU0sU0FBUyxHQUFHLDBCQUEwQixDQUFBO0FBQzVDLE1BQU0sWUFBWSxHQUFHLFdBQUksQ0FBQyxTQUFTLEVBQUUsUUFBUSxDQUFDLENBQUE7QUFDOUMsTUFBTSxhQUFhLEdBQUcsV0FBSSxDQUFDLFNBQVMsRUFBRSxVQUFVLENBQUMsQ0FBQTtBQUVqRCxRQUFRLENBQUMsWUFBWSxFQUFFLEdBQUcsRUFBRTtJQUMxQixVQUFVLENBQUMsR0FBRyxFQUFFO1FBQ2QsWUFBRSxDQUFDLElBQUksRUFBRSxZQUFZLENBQUMsQ0FBQTtJQUN4QixDQUFDLENBQUMsQ0FBQTtJQUVGLEVBQUUsQ0FBQyw0RkFBNEYsRUFBRSxLQUFLLElBQUksRUFBRTtRQUMxRyxNQUFNLGlCQUFPLENBQUMsR0FBRyxDQUFDO1lBQ2hCLFlBQVksV0FBSSxDQUFDLGFBQWEsRUFBRSxrQkFBa0IsQ0FBQyxFQUFFO1lBQ3JELFNBQVM7U0FDVixDQUFDLENBQUE7UUFFRixNQUFNLENBQUMsQ0FBQyxHQUFHLFlBQUUsQ0FBQyxJQUFJLEVBQUUsWUFBWSxDQUFDLENBQUMsQ0FBQyxDQUFDLGVBQWUsRUFBRSxDQUFBO0lBQ3ZELENBQUMsQ0FBQyxDQUFBO0FBQ0osQ0FBQyxDQUFDLENBQUEiLCJzb3VyY2VzQ29udGVudCI6WyJpbXBvcnQgeyBqb2luIH0gZnJvbSAncGF0aCdcbmltcG9ydCB7IGxzLCBybSB9IGZyb20gJ3NoZWxsanMnXG5pbXBvcnQgY29tcGlsZSBmcm9tICcuLi8uLi8uLi9zcmMvY29tbWFuZHMvY29tcGlsZSdcblxuY29uc3QgVEVTVF9QQVRIID0gJ3Rlc3QvZnVuY3Rpb25hbC90cnVmZmxlLydcbmNvbnN0IFRFU1RfRlNfUEFUSCA9IGpvaW4oVEVTVF9QQVRILCAndGVzdGZzJylcbmNvbnN0IEZJWFRVUkVTX1BBVEggPSBqb2luKFRFU1RfUEFUSCwgJ2ZpeHR1cmVzJylcblxuZGVzY3JpYmUoJ2NvbXBpbGVBbGwnLCAoKSA9PiB7XG4gIGJlZm9yZUVhY2goKCkgPT4ge1xuICAgIHJtKCctcicsIFRFU1RfRlNfUEFUSClcbiAgfSlcblxuICBpdCgnc2hvdWxkIHByb2R1Y2UgYSB0cnVmZmxlIGNvbnRyYWN0IGFic3RyYWN0aW9uIHZpYSBhIGpzb24gYXJ0aWZhY3QgcHJvZHVjZWQgYnkgc29sLWNvbXBpbGVyJywgYXN5bmMgKCkgPT4ge1xuICAgIGF3YWl0IGNvbXBpbGUucnVuKFtcbiAgICAgIGAtLWNvbmZpZz0ke2pvaW4oRklYVFVSRVNfUEFUSCwgJ3Rlc3QuY29uZmlnLmpzb24nKX1gLFxuICAgICAgJ3RydWZmbGUnLFxuICAgIF0pXG5cbiAgICBleHBlY3QoWy4uLmxzKCctUicsIFRFU1RfRlNfUEFUSCldKS50b01hdGNoU25hcHNob3QoKVxuICB9KVxufSlcbiJdfQ==