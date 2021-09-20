"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.compileAll = void 0;
const path_1 = require("path");
const ts_generator_1 = require("ts-generator");
const TypeChain_1 = require("typechain/dist/TypeChain");
const utils_1 = require("../utils");
/**
 * Generate ethers.js contract abstractions for all of the solidity versions under a specified contract
 * directory.
 *
 * @param conf The application configuration, e.g. where to read artifacts, where to output, etc..
 */
async function compileAll(conf) {
    const cwd = process.cwd();
    return Promise.all(utils_1.getArtifactDirs(conf).map(async ({ dir }) => {
        const c = compiler(conf, cwd, dir);
        await ts_generator_1.tsGenerator({ cwd, loggingLvl: 'verbose' }, c);
    }));
}
exports.compileAll = compileAll;
/**
 * Create a typechain compiler instance that reads in a subdirectory of artifacts e.g. (abi/v0.4, abi/v0.5.. etc)
 * and outputs ethers contract abstractions under the same version prefix, (ethers/v0.4, ethers/v0.4.. etc)
 *
 * @param config The application level config for compilation
 * @param cwd The current working directory during this programs execution
 * @param subDir The subdirectory to use as a namespace when reading artifacts and outputting
 * contract abstractions
 */
function compiler({ artifactsDir, contractAbstractionDir }, cwd, subDir) {
    return new TypeChain_1.TypeChain({
        cwd,
        rawConfig: {
            files: path_1.join(artifactsDir, subDir, '**', '*.json'),
            outDir: path_1.join(contractAbstractionDir, 'ethers', subDir),
            target: 'ethers-v4',
        },
    });
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZXRoZXJzLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL3NlcnZpY2VzL2NvbXBpbGVycy9ldGhlcnMudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7O0FBQUEsK0JBQTJCO0FBQzNCLCtDQUEwQztBQUMxQyx3REFBb0Q7QUFFcEQsb0NBQTBDO0FBRTFDOzs7OztHQUtHO0FBQ0ksS0FBSyxVQUFVLFVBQVUsQ0FBQyxJQUFnQjtJQUMvQyxNQUFNLEdBQUcsR0FBRyxPQUFPLENBQUMsR0FBRyxFQUFFLENBQUE7SUFFekIsT0FBTyxPQUFPLENBQUMsR0FBRyxDQUNoQix1QkFBZSxDQUFDLElBQUksQ0FBQyxDQUFDLEdBQUcsQ0FBQyxLQUFLLEVBQUUsRUFBRSxHQUFHLEVBQUUsRUFBRSxFQUFFO1FBQzFDLE1BQU0sQ0FBQyxHQUFHLFFBQVEsQ0FBQyxJQUFJLEVBQUUsR0FBRyxFQUFFLEdBQUcsQ0FBQyxDQUFBO1FBQ2xDLE1BQU0sMEJBQVcsQ0FBQyxFQUFFLEdBQUcsRUFBRSxVQUFVLEVBQUUsU0FBUyxFQUFFLEVBQUUsQ0FBQyxDQUFDLENBQUE7SUFDdEQsQ0FBQyxDQUFDLENBQ0gsQ0FBQTtBQUNILENBQUM7QUFURCxnQ0FTQztBQUVEOzs7Ozs7OztHQVFHO0FBQ0gsU0FBUyxRQUFRLENBQ2YsRUFBRSxZQUFZLEVBQUUsc0JBQXNCLEVBQWMsRUFDcEQsR0FBVyxFQUNYLE1BQWM7SUFFZCxPQUFPLElBQUkscUJBQVMsQ0FBQztRQUNuQixHQUFHO1FBQ0gsU0FBUyxFQUFFO1lBQ1QsS0FBSyxFQUFFLFdBQUksQ0FBQyxZQUFZLEVBQUUsTUFBTSxFQUFFLElBQUksRUFBRSxRQUFRLENBQUM7WUFDakQsTUFBTSxFQUFFLFdBQUksQ0FBQyxzQkFBc0IsRUFBRSxRQUFRLEVBQUUsTUFBTSxDQUFDO1lBQ3RELE1BQU0sRUFBRSxXQUFXO1NBQ3BCO0tBQ0YsQ0FBQyxDQUFBO0FBQ0osQ0FBQyIsInNvdXJjZXNDb250ZW50IjpbImltcG9ydCB7IGpvaW4gfSBmcm9tICdwYXRoJ1xuaW1wb3J0IHsgdHNHZW5lcmF0b3IgfSBmcm9tICd0cy1nZW5lcmF0b3InXG5pbXBvcnQgeyBUeXBlQ2hhaW4gfSBmcm9tICd0eXBlY2hhaW4vZGlzdC9UeXBlQ2hhaW4nXG5pbXBvcnQgKiBhcyBjb25maWcgZnJvbSAnLi4vY29uZmlnJ1xuaW1wb3J0IHsgZ2V0QXJ0aWZhY3REaXJzIH0gZnJvbSAnLi4vdXRpbHMnXG5cbi8qKlxuICogR2VuZXJhdGUgZXRoZXJzLmpzIGNvbnRyYWN0IGFic3RyYWN0aW9ucyBmb3IgYWxsIG9mIHRoZSBzb2xpZGl0eSB2ZXJzaW9ucyB1bmRlciBhIHNwZWNpZmllZCBjb250cmFjdFxuICogZGlyZWN0b3J5LlxuICpcbiAqIEBwYXJhbSBjb25mIFRoZSBhcHBsaWNhdGlvbiBjb25maWd1cmF0aW9uLCBlLmcuIHdoZXJlIHRvIHJlYWQgYXJ0aWZhY3RzLCB3aGVyZSB0byBvdXRwdXQsIGV0Yy4uXG4gKi9cbmV4cG9ydCBhc3luYyBmdW5jdGlvbiBjb21waWxlQWxsKGNvbmY6IGNvbmZpZy5BcHApIHtcbiAgY29uc3QgY3dkID0gcHJvY2Vzcy5jd2QoKVxuXG4gIHJldHVybiBQcm9taXNlLmFsbChcbiAgICBnZXRBcnRpZmFjdERpcnMoY29uZikubWFwKGFzeW5jICh7IGRpciB9KSA9PiB7XG4gICAgICBjb25zdCBjID0gY29tcGlsZXIoY29uZiwgY3dkLCBkaXIpXG4gICAgICBhd2FpdCB0c0dlbmVyYXRvcih7IGN3ZCwgbG9nZ2luZ0x2bDogJ3ZlcmJvc2UnIH0sIGMpXG4gICAgfSksXG4gIClcbn1cblxuLyoqXG4gKiBDcmVhdGUgYSB0eXBlY2hhaW4gY29tcGlsZXIgaW5zdGFuY2UgdGhhdCByZWFkcyBpbiBhIHN1YmRpcmVjdG9yeSBvZiBhcnRpZmFjdHMgZS5nLiAoYWJpL3YwLjQsIGFiaS92MC41Li4gZXRjKVxuICogYW5kIG91dHB1dHMgZXRoZXJzIGNvbnRyYWN0IGFic3RyYWN0aW9ucyB1bmRlciB0aGUgc2FtZSB2ZXJzaW9uIHByZWZpeCwgKGV0aGVycy92MC40LCBldGhlcnMvdjAuNC4uIGV0YylcbiAqXG4gKiBAcGFyYW0gY29uZmlnIFRoZSBhcHBsaWNhdGlvbiBsZXZlbCBjb25maWcgZm9yIGNvbXBpbGF0aW9uXG4gKiBAcGFyYW0gY3dkIFRoZSBjdXJyZW50IHdvcmtpbmcgZGlyZWN0b3J5IGR1cmluZyB0aGlzIHByb2dyYW1zIGV4ZWN1dGlvblxuICogQHBhcmFtIHN1YkRpciBUaGUgc3ViZGlyZWN0b3J5IHRvIHVzZSBhcyBhIG5hbWVzcGFjZSB3aGVuIHJlYWRpbmcgYXJ0aWZhY3RzIGFuZCBvdXRwdXR0aW5nXG4gKiBjb250cmFjdCBhYnN0cmFjdGlvbnNcbiAqL1xuZnVuY3Rpb24gY29tcGlsZXIoXG4gIHsgYXJ0aWZhY3RzRGlyLCBjb250cmFjdEFic3RyYWN0aW9uRGlyIH06IGNvbmZpZy5BcHAsXG4gIGN3ZDogc3RyaW5nLFxuICBzdWJEaXI6IHN0cmluZyxcbik6IFR5cGVDaGFpbiB7XG4gIHJldHVybiBuZXcgVHlwZUNoYWluKHtcbiAgICBjd2QsXG4gICAgcmF3Q29uZmlnOiB7XG4gICAgICBmaWxlczogam9pbihhcnRpZmFjdHNEaXIsIHN1YkRpciwgJyoqJywgJyouanNvbicpLFxuICAgICAgb3V0RGlyOiBqb2luKGNvbnRyYWN0QWJzdHJhY3Rpb25EaXIsICdldGhlcnMnLCBzdWJEaXIpLFxuICAgICAgdGFyZ2V0OiAnZXRoZXJzLXY0JyxcbiAgICB9LFxuICB9KVxufVxuIl19