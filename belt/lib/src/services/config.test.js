"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const tslib_1 = require("tslib");
const mock_fs_1 = tslib_1.__importDefault(require("mock-fs"));
const config_1 = require("./config");
describe('config.load', () => {
    const strify = JSON.stringify;
    function getDefaultConf() {
        return {
            contractsDir: 'src',
            artifactsDir: 'abi',
            contractAbstractionDir: '.',
            useDockerisedSolc: true,
            compilerSettings: {
                versions: {
                    'v0.4': '0.4.24',
                    'v0.5': '0.5.0',
                    'v0.6': '0.6.2',
                },
            },
            publicVersions: ['0.4.24', '0.5.0'],
        };
    }
    afterEach(() => {
        mock_fs_1.default.restore();
    });
    it('should throw on a missing config', () => {
        mock_fs_1.default({
            src: '',
        });
        expect(() => config_1.load('./doesnte')).toThrowError('Could not load config');
    });
    it('should throw on a missing contracts directory', () => {
        mock_fs_1.default({
            conf: strify(getDefaultConf()),
        });
        expect(() => config_1.load('./conf')).toThrowError('Expected value of config.contractsDir to be a directory');
    });
    it('should throw on a non-string artifacts directory value', () => {
        mock_fs_1.default({
            src: {},
            conf: strify({ ...getDefaultConf(), artifactsDir: 5 }),
        });
        expect(() => config_1.load('./conf')).toThrowError('Expected value of config.artifactsDir to be a string');
    });
    it('should throw on a non-boolean useDockerisedSolc value', () => {
        mock_fs_1.default({
            src: {},
            conf: strify({ ...getDefaultConf(), useDockerisedSolc: '' }),
        });
        expect(() => config_1.load('./conf')).toThrowError('Expected value of config.useDockerisedSolc to be a boolean');
    });
    it('should throw on an invalid contractAbstractionDir value', () => {
        mock_fs_1.default({
            src: {},
            conf: strify({ ...getDefaultConf(), contractAbstractionDir: 5 }),
        });
        expect(() => config_1.load('./conf')).toThrowError('Expected value of config.contractAbstractionDir to be a string');
    });
    it('should throw on a non-valid compilerSettings value', () => {
        mock_fs_1.default({
            src: {},
            conf: strify({ ...getDefaultConf(), compilerSettings: undefined }),
        });
        expect(() => config_1.load('./conf')).toThrowError('Expected value of config.compilerSettings to be an object');
    });
    it('should throw on a non-valid compilerSettings value', () => {
        mock_fs_1.default({
            src: {},
            conf: strify({ ...getDefaultConf(), compilerSettings: {} }),
        });
        expect(() => config_1.load('./conf')).toThrowError('Expected value of config.compilerSettings.versions to be a dictionary');
    });
    it('should load a config correctly', () => {
        mock_fs_1.default({
            src: {},
            conf: strify(getDefaultConf()),
        });
        expect(config_1.load('./conf')).toStrictEqual({ ...getDefaultConf() });
    });
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY29uZmlnLnRlc3QuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi9zcmMvc2VydmljZXMvY29uZmlnLnRlc3QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7O0FBQUEsOERBQTBCO0FBQzFCLHFDQUFvQztBQUVwQyxRQUFRLENBQUMsYUFBYSxFQUFFLEdBQUcsRUFBRTtJQUMzQixNQUFNLE1BQU0sR0FBRyxJQUFJLENBQUMsU0FBUyxDQUFBO0lBRTdCLFNBQVMsY0FBYztRQUNyQixPQUFPO1lBQ0wsWUFBWSxFQUFFLEtBQUs7WUFDbkIsWUFBWSxFQUFFLEtBQUs7WUFDbkIsc0JBQXNCLEVBQUUsR0FBRztZQUMzQixpQkFBaUIsRUFBRSxJQUFJO1lBQ3ZCLGdCQUFnQixFQUFFO2dCQUNoQixRQUFRLEVBQUU7b0JBQ1IsTUFBTSxFQUFFLFFBQVE7b0JBQ2hCLE1BQU0sRUFBRSxPQUFPO29CQUNmLE1BQU0sRUFBRSxPQUFPO2lCQUNoQjthQUNGO1lBQ0QsY0FBYyxFQUFFLENBQUMsUUFBUSxFQUFFLE9BQU8sQ0FBQztTQUNwQyxDQUFBO0lBQ0gsQ0FBQztJQUVELFNBQVMsQ0FBQyxHQUFHLEVBQUU7UUFDYixpQkFBSSxDQUFDLE9BQU8sRUFBRSxDQUFBO0lBQ2hCLENBQUMsQ0FBQyxDQUFBO0lBRUYsRUFBRSxDQUFDLGtDQUFrQyxFQUFFLEdBQUcsRUFBRTtRQUMxQyxpQkFBSSxDQUFDO1lBQ0gsR0FBRyxFQUFFLEVBQUU7U0FDUixDQUFDLENBQUE7UUFFRixNQUFNLENBQUMsR0FBRyxFQUFFLENBQUMsYUFBSSxDQUFDLFdBQVcsQ0FBQyxDQUFDLENBQUMsWUFBWSxDQUFDLHVCQUF1QixDQUFDLENBQUE7SUFDdkUsQ0FBQyxDQUFDLENBQUE7SUFFRixFQUFFLENBQUMsK0NBQStDLEVBQUUsR0FBRyxFQUFFO1FBQ3ZELGlCQUFJLENBQUM7WUFDSCxJQUFJLEVBQUUsTUFBTSxDQUFDLGNBQWMsRUFBRSxDQUFDO1NBQy9CLENBQUMsQ0FB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cG9ydCB7IEFwcCwgbG9hZCB9IGZyb20gJy4vY29uZmlnJ1xuXG5kZXNjcmliZSgnY29uZmlnLmxvYWQnLCAoKSA9PiB7XG4gIGNvbnN0IHN0cmlmeSA9IEpTT04uc3RyaW5naWZ5XG5cbiAgZnVuY3Rpb24gZ2V0RGVmYXVsdENvbmYoKTogQXBwIHtcbiAgICByZXR1cm4ge1xuICAgICAgY29udHJhY3RzRGlyOiAnc3JjJyxcbiAgICAgIGFydGlmYWN0c0RpcjogJ2FiaScsXG4gICAgICBjb250cmFjdEFic3RyYWN0aW9uRGlyOiAnLicsXG4gICAgICB1c2VEb2NrZXJpc2VkU29sYzogdHJ1ZSxcbiAgICAgIGNvbXBpbGVyU2V0dGluZ3M6IHtcbiAgICAgICAgdmVyc2lvbnM6IHtcbiAgICAgICAgICAndjAuNCc6ICcwLjQuMjQnLFxuICAgICAgICAgICd2MC41JzogJzAuNS4wJyxcbiAgICAgICAgICAndjAuNic6ICcwLjYuMicsXG4gICAgICAgIH0sXG4gICAgICB9LFxuICAgICAgcHVibGljVmVyc2lvbnM6IFsnMC40LjI0JywgJzAuNS4wJ10sXG4gICAgfVxuICB9XG5cbiAgYWZ0ZXJFYWNoKCgpID0+IHtcbiAgICBtb2NrLnJlc3RvcmUoKVxuICB9KVxuXG4gIGl0KCdzaG91bGQgdGhyb3cgb24gYSBtaXNzaW5nIGNvbmZpZycsICgpID0+IHtcbiAgICBtb2NrKHtcbiAgICAgIHNyYzogJycsXG4gICAgfSlcblxuICAgIGV4cGVjdCgoKSA9PiBsb2FkKCcuL2RvZXNudGUnKSkudG9UaHJvd0Vycm9yKCdDb3VsZCBub3QgbG9hZCBjb25maWcnKVxuICB9KVxuXG4gIGl0KCdzaG91bGQgdGhyb3cgb24gYSBtaXNzaW5nIGNvbnRyYWN0cyBkaXJlY3RvcnknLCAoKSA9PiB7XG4gICAgbW9jayh7XG4gICAgICBjb25mOiBzdHJpZnkoZ2V0RGVmYXVsdENvbmYoKSksXG4gICAgfSlcblxuICAgIGV4cGVjdCgoKSA9PiBsb2FkKCcuL2NvbmYnKSkudG9UaHJvd0Vycm9yKFxuICAgICAgJ0V4cGVjdGVkIHZhbHVlIG9mIGNvbmZpZy5jb250cmFjdHNEaXIgdG8gYmUgYSBkaXJlY3RvcnknLFxuICAgIClcbiAgfSlcblxuICBpdCgnc2hvdWxkIHRocm93IG9uIGEgbm9uLXN0cmluZyBhcnRpZmFjdHMgZGlyZWN0b3J5IHZhbHVlJywgKCkgPT4ge1xuICAgIG1vY2soe1xuICAgICAgc3JjOiB7fSxcbiAgICAgIGNvbmY6IHN0cmlmeSh7IC4uLmdldERlZmF1bHRDb25mKCksIGFydGlmYWN0c0RpcjogNSB9KSxcbiAgICB9KVxuXG4gICAgZXhwZWN0KCgpID0+IGxvYWQoJy4vY29uZicpKS50b1Rocm93RXJyb3IoXG4gICAgICAnRXhwZWN0ZWQgdmFsdWUgb2YgY29uZmlnLmFydGlmYWN0c0RpciB0byBiZSBhIHN0cmluZycsXG4gICAgKVxuICB9KVxuXG4gIGl0KCdzaG91bGQgdGhyb3cgb24gYSBub24tYm9vbGVhbiB1c2VEb2NrZXJpc2VkU29sYyB2YWx1ZScsICgpID0+IHtcbiAgICBtb2NrKHtcbiAgICAgIHNyYzoge30sXG4gICAgICBjb25mOiBzdHJpZnkoeyAuLi5nZXREZWZhdWx0Q29uZigpLCB1c2VEb2NrZXJpc2VkU29sYzogJycgfSksXG4gICAgfSlcblxuICAgIGV4cGVjdCgoKSA9PiBsb2FkKCcuL2NvbmYnKSkudG9UaHJvd0Vycm9yKFxuICAgICAgJ0V4cGVjdGVkIHZhbHVlIG9mIGNvbmZpZy51c2VEb2NrZXJpc2VkU29sYyB0byBiZSBhIGJvb2xlYW4nLFxuICAgIClcbiAgfSlcblxuICBpdCgnc2hvdWxkIHRocm93IG9uIGFuIGludmFsaWQgY29udHJhY3RBYnN0cmFjdGlvbkRpciB2YWx1ZScsICgpID0+IHtcbiAgICBtb2NrKHtcbiAgICAgIHNyYzoge30sXG4gICAgICBjb25mOiBzdHJpZnkoeyAuLi5nZXREZWZhdWx0Q29uZigpLCBjb250cmFjdEFic3RyYWN0aW9uRGlyOiA1IH0pLFxuICAgIH0pXG5cbiAgICBleHBlY3QoKCkgPT4gbG9hZCgnLi9jb25mJykpLnRvVGhyb3dFcnJvcihcbiAgICAgICdFeHBlY3RlZCB2YWx1ZSBvZiBjb25maWcuY29udHJhY3RBYnN0cmFjdGlvbkRpciB0byBiZSBhIHN0cmluZycsXG4gICAgKVxuICB9KVxuXG4gIGl0KCdzaG91bGQgdGhyb3cgb24gYSBub24tdmFsaWQgY29tcGlsZXJTZXR0aW5ncyB2YWx1ZScsICgpID0+IHtcbiAgICBtb2NrKHtcbiAgICAgIHNyYzoge30sXG4gICAgICBjb25mOiBzdHJpZnkoeyAuLi5nZXREZWZhdWx0Q29uZigpLCBjb21waWxlclNldHRpbmdzOiB1bmRlZmluZWQgfSksXG4gICAgfSlcblxuICAgIGV4cGVjdCgoKSA9PiBsb2FkKCcuL2NvbmYnKSkudG9UaHJvd0Vycm9yKFxuICAgICAgJ0V4cGVjdGVkIHZhbHVlIG9mIGNvbmZpZy5jb21waWxlclNldHRpbmdzIHRvIGJlIGFuIG9iamVjdCcsXG4gICAgKVxuICB9KVxuXG4gIGl0KCdzaG91bGQgdGhyb3cgb24gYSBub24tdmFsaWQgY29tcGlsZXJTZXR0aW5ncyB2YWx1ZScsICgpID0+IHtcbiAgICBtb2NrKHtcbiAgICAgIHNyYzoge30sXG4gICAgICBjb25mOiBzdHJpZnkoeyAuLi5nZXREZWZhdWx0Q29uZigpLCBjb21waWxlclNldHRpbmdzOiB7fSB9KSxcbiAgICB9KVxuXG4gICAgZXhwZWN0KCgpID0+IGxvYWQoJy4vY29uZicpKS50b1Rocm93RXJyb3IoXG4gICAgICAnRXhwZWN0ZWQgdmFsdWUgb2YgY29uZmlnLmNvbXBpbGVyU2V0dGluZ3MudmVyc2lvbnMgdG8gYmUgYSBkaWN0aW9uYXJ5JyxcbiAgICApXG4gIH0pXG5cbiAgaXQoJ3Nob3VsZCBsb2FkIGEgY29uZmlnIGNvcnJlY3RseScsICgpID0+IHtcbiAgICBtb2NrKHtcbiAgICAgIHNyYzoge30sXG4gICAgICBjb25mOiBzdHJpZnkoZ2V0RGVmYXVsdENvbmYoKSksXG4gICAgfSlcblxuICAgIGV4cGVjdChsb2FkKCcuL2NvbmYnKSkudG9TdHJpY3RFcXVhbCh7IC4uLmdldERlZmF1bHRDb25mKCkgfSlcbiAgfSlcbn0pXG4iXX0=