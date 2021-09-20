"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getJsonFile = exports.debug = exports.getArtifactDirs = exports.getContractDirs = void 0;
const tslib_1 = require("tslib");
const debug_1 = tslib_1.__importDefault(require("debug"));
const fs_1 = require("fs");
const shelljs_1 = require("shelljs");
/**
 * Get contract versions and their directories
 */
function getContractDirs(conf) {
    const contractsMap = [];
    for (const dir in conf.compilerSettings.versions) {
        contractsMap.push({
            dir,
            version: conf.compilerSettings.versions[dir],
        });
    }
    return contractsMap;
}
exports.getContractDirs = getContractDirs;
/**
 * Get artifact verions and their directories
 */
function getArtifactDirs(conf) {
    const artifactDirs = shelljs_1.ls(conf.artifactsDir);
    return artifactDirs.map((d) => ({
        dir: d,
        version: conf.compilerSettings.versions[d],
    }));
}
exports.getArtifactDirs = getArtifactDirs;
/**
 * Create a logger specifically for debugging. The root level namespace is based on the package name.
 *
 * @see https://www.npmjs.com/package/debug
 * @param fileName The filename that the debug logger is being used in for namespacing purposes.
 */
function debug(fileName) {
    return debug_1.default('belt').extend(fileName);
}
exports.debug = debug;
/**
 * Load a json file at the specified path.
 *
 * @param path The file path relative to the cwd to read in the json file from.
 */
function getJsonFile(path) {
    return JSON.parse(fs_1.readFileSync(path, 'utf8'));
}
exports.getJsonFile = getJsonFile;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidXRpbHMuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi9zcmMvc2VydmljZXMvdXRpbHMudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7OztBQUFBLDBEQUFxQjtBQUNyQiwyQkFBaUM7QUFDakMscUNBQTRCO0FBRzVCOztHQUVHO0FBQ0gsU0FBZ0IsZUFBZSxDQUFDLElBQWdCO0lBQzlDLE1BQU0sWUFBWSxHQUFHLEVBQUUsQ0FBQTtJQUN2QixLQUFLLE1BQU0sR0FBRyxJQUFJLElBQUksQ0FBQyxnQkFBZ0IsQ0FBQyxRQUFRLEVBQUU7UUFDaEQsWUFBWSxDQUFDLElBQUksQ0FBQztZQUNoQixHQUFHO1lBQ0gsT0FBTyxFQUFFLElBQUksQ0FBQyxnQkFBZ0IsQ0FBQyxRQUFRLENBQUMsR0FBRyxDQUFDO1NBQzdDLENBQUMsQ0FBQTtLQUNIO0lBQ0QsT0FBTyxZQUFZLENBQUE7QUFDckIsQ0FBQztBQVRELDBDQVNDO0FBRUQ7O0dBRUc7QUFDSCxTQUFnQixlQUFlLENBQUMsSUFBZ0I7SUFDOUMsTUFBTSxZQUFZLEdBQUcsWUFBRSxDQUFDLElBQUksQ0FBQyxZQUFZLENBQUMsQ0FBQTtJQUUxQyxPQUFPLFlBQVksQ0FBQyxHQUFHLENBQUMsQ0FBQyxDQUFDLEVBQUUsRUFBRSxDQUFDLENBQUM7UUFDOUIsR0FBRyxFQUFFLENBQUM7UUFDTixPQUFPLEVBQUUsSUFBSSxDQUFDLGdCQUFnQixDQUFDLFFBQVEsQ0FBQyxDQUFDLENBQUM7S0FDM0MsQ0FBQyxDQUFDLENBQUE7QUFDTCxDQUFDO0FBUEQsMENBT0M7QUFFRDs7Ozs7R0FLRztBQUNILFNBQWdCLEtBQUssQ0FBQyxRQUFnQjtJQUNwQyxPQUFPLGVBQUMsQ0FBQyxNQUFNLENBQUMsQ0FBQyxNQUFNLENBQUMsUUFBUSxDQUFDLENBQUE7QUFDbkMsQ0FBQztBQUZELHNCQUVDO0FBRUQ7Ozs7R0FJRztBQUNILFNBQWdCLFdBQVcsQ0FBQyxJQUFZO0lBQ3RDLE9BQU8sSUFBSSxDQUFDLEtBQUssQ0FBQyxpQkFBWSxDQUFDLElBQUksRUFBRSxNQUFNLENBQUMsQ0FBQyxDQUFBO0FBQy9DLENBQUM7QUFGRCxrQ0FFQyIsInNvdXJjZXNDb250ZW50IjpbImltcG9ydCBkIGZyb20gJ2RlYnVnJ1xuaW1wb3J0IHsgcmVhZEZpbGVTeW5jIH0gZnJvbSAnZnMnXG5pbXBvcnQgeyBscyB9IGZyb20gJ3NoZWxsanMnXG5pbXBvcnQgKiBhcyBjb25maWcgZnJvbSAnLi9jb25maWcnXG5cbi8qKlxuICogR2V0IGNvbnRyYWN0IHZlcnNpb25zIGFuZCB0aGVpciBkaXJlY3Rvcmllc1xuICovXG5leHBvcnQgZnVuY3Rpb24gZ2V0Q29udHJhY3REaXJzKGNvbmY6IGNvbmZpZy5BcHApIHtcbiAgY29uc3QgY29udHJhY3RzTWFwID0gW11cbiAgZm9yIChjb25zdCBkaXIgaW4gY29uZi5jb21waWxlclNldHRpbmdzLnZlcnNpb25zKSB7XG4gICAgY29udHJhY3RzTWFwLnB1c2goe1xuICAgICAgZGlyLFxuICAgICAgdmVyc2lvbjogY29uZi5jb21waWxlclNldHRpbmdzLnZlcnNpb25zW2Rpcl0sXG4gICAgfSlcbiAgfVxuICByZXR1cm4gY29udHJhY3RzTWFwXG59XG5cbi8qKlxuICogR2V0IGFydGlmYWN0IHZlcmlvbnMgYW5kIHRoZWlyIGRpcmVjdG9yaWVzXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBnZXRBcnRpZmFjdERpcnMoY29uZjogY29uZmlnLkFwcCkge1xuICBjb25zdCBhcnRpZmFjdERpcnMgPSBscyhjb25mLmFydGlmYWN0c0RpcilcblxuICByZXR1cm4gYXJ0aWZhY3REaXJzLm1hcCgoZCkgPT4gKHtcbiAgICBkaXI6IGQsXG4gICAgdmVyc2lvbjogY29uZi5jb21waWxlclNldHRpbmdzLnZlcnNpb25zW2RdLFxuICB9KSlcbn1cblxuLyoqXG4gKiBDcmVhdGUgYSBsb2dnZXIgc3BlY2lmaWNhbGx5IGZvciBkZWJ1Z2dpbmcuIFRoZSByb290IGxldmVsIG5hbWVzcGFjZSBpcyBiYXNlZCBvbiB0aGUgcGFja2FnZSBuYW1lLlxuICpcbiAqIEBzZWUgaHR0cHM6Ly93d3cubnBtanMuY29tL3BhY2thZ2UvZGVidWdcbiAqIEBwYXJhbSBmaWxlTmFtZSBUaGUgZmlsZW5hbWUgdGhhdCB0aGUgZGVidWcgbG9nZ2VyIGlzIGJlaW5nIHVzZWQgaW4gZm9yIG5hbWVzcGFjaW5nIHB1cnBvc2VzLlxuICovXG5leHBvcnQgZnVuY3Rpb24gZGVidWcoZmlsZU5hbWU6IHN0cmluZykge1xuICByZXR1cm4gZCgnYmVsdCcpLmV4dGVuZChmaWxlTmFtZSlcbn1cblxuLyoqXG4gKiBMb2FkIGEganNvbiBmaWxlIGF0IHRoZSBzcGVjaWZpZWQgcGF0aC5cbiAqXG4gKiBAcGFyYW0gcGF0aCBUaGUgZmlsZSBwYXRoIHJlbGF0aXZlIHRvIHRoZSBjd2QgdG8gcmVhZCBpbiB0aGUganNvbiBmaWxlIGZyb20uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBnZXRKc29uRmlsZShwYXRoOiBzdHJpbmcpOiB1bmtub3duIHtcbiAgcmV0dXJuIEpTT04ucGFyc2UocmVhZEZpbGVTeW5jKHBhdGgsICd1dGY4JykpXG59XG4iXX0=