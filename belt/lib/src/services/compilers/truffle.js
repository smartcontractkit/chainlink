"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.compileAll = void 0;
const fs_1 = require("fs");
const path_1 = require("path");
const shelljs_1 = require("shelljs");
const utils_1 = require("../utils");
/**
 * Generate @truffle/contract abstractions for all of the solidity versions under a specified contract
 * directory.
 *
 * @param conf The application configuration, e.g. where to read artifacts, where to output, etc..
 */
async function compileAll(conf) {
    utils_1.getArtifactDirs(conf).forEach(({ dir }) => {
        getContractPathsPer(conf, dir).forEach((p) => {
            const json = utils_1.getJsonFile(p);
            const fileName = path_1.basename(p, '.json');
            const file = fillTemplate(fileName, {
                contractName: json.contractName,
                abi: json.compilerOutput.abi,
                evm: json.compilerOutput.evm,
                metadata: json.compilerOutput.metadata,
            });
            write(path_1.join(conf.contractAbstractionDir, 'truffle', dir), fileName, file);
        });
    });
}
exports.compileAll = compileAll;
/**
 * Create a truffle contract abstraction file
 *
 * @param contractName The name of the contract that will be exported
 * @param contractArgs The arguments to pass to @truffle/contract
 */
function fillTemplate(contractName, contractArgs) {
    return `'use strict'
Object.defineProperty(exports, '__esModule', { value: true })
const contract = require('@truffle/contract')
const ${contractName} = contract(${JSON.stringify(contractArgs, null, 1)})

if (process.env.NODE_ENV === 'test') {
  try {
    eval('${contractName}.setProvider(web3.currentProvider)')
  } catch (e) {}
}

exports.${contractName} = ${contractName}
`;
}
function getContractPathsPer({ artifactsDir }, version) {
    return [...shelljs_1.ls(path_1.join(artifactsDir, version, '/**/*.json'))];
}
function write(outPath, fileName, file) {
    shelljs_1.mkdir('-p', outPath);
    fs_1.writeFileSync(path_1.join(outPath, `${fileName}.js`), file);
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidHJ1ZmZsZS5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3NyYy9zZXJ2aWNlcy9jb21waWxlcnMvdHJ1ZmZsZS50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7QUFDQSwyQkFBa0M7QUFDbEMsK0JBQXFDO0FBQ3JDLHFDQUFtQztBQUVuQyxvQ0FBdUQ7QUFFdkQ7Ozs7O0dBS0c7QUFDSSxLQUFLLFVBQVUsVUFBVSxDQUFDLElBQWdCO0lBQy9DLHVCQUFlLENBQUMsSUFBSSxDQUFDLENBQUMsT0FBTyxDQUFDLENBQUMsRUFBRSxHQUFHLEVBQUUsRUFBRSxFQUFFO1FBQ3hDLG1CQUFtQixDQUFDLElBQUksRUFBRSxHQUFHLENBQUMsQ0FBQyxPQUFPLENBQUMsQ0FBQyxDQUFDLEVBQUUsRUFBRTtZQUMzQyxNQUFNLElBQUksR0FBUSxtQkFBVyxDQUFDLENBQUMsQ0FBQyxDQUFBO1lBQ2hDLE1BQU0sUUFBUSxHQUFHLGVBQVEsQ0FBQyxDQUFDLEVBQUUsT0FBTyxDQUFDLENBQUE7WUFDckMsTUFBTSxJQUFJLEdBQUcsWUFBWSxDQUFDLFFBQVEsRUFBRTtnQkFDbEMsWUFBWSxFQUFFLElBQUksQ0FBQyxZQUFZO2dCQUMvQixHQUFHLEVBQUUsSUFBSSxDQUFDLGNBQWMsQ0FBQyxHQUFHO2dCQUM1QixHQUFHLEVBQUUsSUFBSSxDQUFDLGNBQWMsQ0FBQyxHQUFHO2dCQUM1QixRQUFRLEVBQUUsSUFBSSxDQUFDLGNBQWMsQ0FBQyxRQUFRO2FBQ3ZDLENBQUMsQ0FBQTtZQUVGLEtBQUssQ0FBQyxXQUFJLENBQUMsSUFBSSxDQUFDLHNCQUFzQixFQUFFLFNBQVMsRUFBRSxHQUFHLENBQUMsRUFBRSxRQUFRLEVBQUUsSUFBSSxDQUFDLENBQUE7UUFDMUUsQ0FBQyxDQUFDLENBQUE7SUFDSixDQUFDLENBQUMsQ0FBQTtBQUNKLENBQUM7QUFmRCxnQ0FlQztBQUVEOzs7OztHQUtHO0FBQ0gsU0FBUyxZQUFZLENBQ25CLFlBQW9CLEVBQ3BCLFlBQTRCO0lBRTVCLE9BQU87OztRQUdELFlBQVksZUFBZSxJQUFJLENBQUMsU0FBUyxDQUFDLFlBQVksRUFBRSxJQUFJLEVBQUUsQ0FBQyxDQUFDOzs7O1lBSTVELFlBQVk7Ozs7VUFJZCxZQUFZLE1BQU0sWUFBWTtDQUN2QyxDQUFBO0FBQ0QsQ0FBQztBQUVELFNBQVMsbUJBQW1CLENBQUMsRUFBRSxZQUFZLEVBQWMsRUFBRSxPQUFlO0lBQ3hFLE9BQU8sQ0FBQyxHQUFHLFlBQUUsQ0FBQyxXQUFJLENBQUMsWUFBWSxFQUFFLE9BQU8sRUFBRSxZQUFZLENBQUMsQ0FBQyxDQUFDLENBQUE7QUFDM0QsQ0FBQztBQUVELFNBQVMsS0FBSyxDQUFDLE9BQWUsRUFBRSxRQUFnQixFQUFFLElBQVk7SUFDNUQsZUFBSyxDQUFDLElBQUksRUFBRSxPQUFPLENBQUMsQ0FBQTtJQUNwQixrQkFBYSxDQUFDLFdBQUksQ0FBQyxPQUFPLEVBQUUsR0FBRyxRQUFRLEtBQUssQ0FBQyxFQUFFLElBQUksQ0FBQyxDQUFBO0FBQ3RELENBQUMiLCJzb3VyY2VzQ29udGVudCI6WyJpbXBvcnQgeyBDb250cmFjdE9iamVjdCB9IGZyb20gJ0B0cnVmZmxlL2NvbnRyYWN0LXNjaGVtYSdcbmltcG9ydCB7IHdyaXRlRmlsZVN5bmMgfSBmcm9tICdmcydcbmltcG9ydCB7IGJhc2VuYW1lLCBqb2luIH0gZnJvbSAncGF0aCdcbmltcG9ydCB7IGxzLCBta2RpciB9IGZyb20gJ3NoZWxsanMnXG5pbXBvcnQgKiBhcyBjb25maWcgZnJvbSAnLi4vY29uZmlnJ1xuaW1wb3J0IHsgZ2V0QXJ0aWZhY3REaXJzLCBnZXRKc29uRmlsZSB9IGZyb20gJy4uL3V0aWxzJ1xuXG4vKipcbiAqIEdlbmVyYXRlIEB0cnVmZmxlL2NvbnRyYWN0IGFic3RyYWN0aW9ucyBmb3IgYWxsIG9mIHRoZSBzb2xpZGl0eSB2ZXJzaW9ucyB1bmRlciBhIHNwZWNpZmllZCBjb250cmFjdFxuICogZGlyZWN0b3J5LlxuICpcbiAqIEBwYXJhbSBjb25mIFRoZSBhcHBsaWNhdGlvbiBjb25maWd1cmF0aW9uLCBlLmcuIHdoZXJlIHRvIHJlYWQgYXJ0aWZhY3RzLCB3aGVyZSB0byBvdXRwdXQsIGV0Yy4uXG4gKi9cbmV4cG9ydCBhc3luYyBmdW5jdGlvbiBjb21waWxlQWxsKGNvbmY6IGNvbmZpZy5BcHApIHtcbiAgZ2V0QXJ0aWZhY3REaXJzKGNvbmYpLmZvckVhY2goKHsgZGlyIH0pID0+IHtcbiAgICBnZXRDb250cmFjdFBhdGhzUGVyKGNvbmYsIGRpcikuZm9yRWFjaCgocCkgPT4ge1xuICAgICAgY29uc3QganNvbjogYW55ID0gZ2V0SnNvbkZpbGUocClcbiAgICAgIGNvbnN0IGZpbGVOYW1lID0gYmFzZW5hbWUocCwgJy5qc29uJylcbiAgICAgIGNvbnN0IGZpbGUgPSBmaWxsVGVtcGxhdGUoZmlsZU5hbWUsIHtcbiAgICAgICAgY29udHJhY3ROYW1lOiBqc29uLmNvbnRyYWN0TmFtZSxcbiAgICAgICAgYWJpOiBqc29uLmNvbXBpbGVyT3V0cHV0LmFiaSxcbiAgICAgICAgZXZtOiBqc29uLmNvbXBpbGVyT3V0cHV0LmV2bSxcbiAgICAgICAgbWV0YWRhdGE6IGpzb24uY29tcGlsZXJPdXRwdXQubWV0YWRhdGEsXG4gICAgICB9KVxuXG4gICAgICB3cml0ZShqb2luKGNvbmYuY29udHJhY3RBYnN0cmFjdGlvbkRpciwgJ3RydWZmbGUnLCBkaXIpLCBmaWxlTmFtZSwgZmlsZSlcbiAgICB9KVxuICB9KVxufVxuXG4vKipcbiAqIENyZWF0ZSBhIHRydWZmbGUgY29udHJhY3QgYWJzdHJhY3Rpb24gZmlsZVxuICpcbiAqIEBwYXJhbSBjb250cmFjdE5hbWUgVGhlIG5hbWUgb2YgdGhlIGNvbnRyYWN0IHRoYXQgd2lsbCBiZSBleHBvcnRlZFxuICogQHBhcmFtIGNvbnRyYWN0QXJncyBUaGUgYXJndW1lbnRzIHRvIHBhc3MgdG8gQHRydWZmbGUvY29udHJhY3RcbiAqL1xuZnVuY3Rpb24gZmlsbFRlbXBsYXRlKFxuICBjb250cmFjdE5hbWU6IHN0cmluZyxcbiAgY29udHJhY3RBcmdzOiBDb250cmFjdE9iamVjdCxcbik6IHN0cmluZyB7XG4gIHJldHVybiBgJ3VzZSBzdHJpY3QnXG5PYmplY3QuZGVmaW5lUHJvcGVydHkoZXhwb3J0cywgJ19fZXNNb2R1bGUnLCB7IHZhbHVlOiB0cnVlIH0pXG5jb25zdCBjb250cmFjdCA9IHJlcXVpcmUoJ0B0cnVmZmxlL2NvbnRyYWN0JylcbmNvbnN0ICR7Y29udHJhY3ROYW1lfSA9IGNvbnRyYWN0KCR7SlNPTi5zdHJpbmdpZnkoY29udHJhY3RBcmdzLCBudWxsLCAxKX0pXG5cbmlmIChwcm9jZXNzLmVudi5OT0RFX0VOViA9PT0gJ3Rlc3QnKSB7XG4gIHRyeSB7XG4gICAgZXZhbCgnJHtjb250cmFjdE5hbWV9LnNldFByb3ZpZGVyKHdlYjMuY3VycmVudFByb3ZpZGVyKScpXG4gIH0gY2F0Y2ggKGUpIHt9XG59XG5cbmV4cG9ydHMuJHtjb250cmFjdE5hbWV9ID0gJHtjb250cmFjdE5hbWV9XG5gXG59XG5cbmZ1bmN0aW9uIGdldENvbnRyYWN0UGF0aHNQZXIoeyBhcnRpZmFjdHNEaXIgfTogY29uZmlnLkFwcCwgdmVyc2lvbjogc3RyaW5nKSB7XG4gIHJldHVybiBbLi4ubHMoam9pbihhcnRpZmFjdHNEaXIsIHZlcnNpb24sICcvKiovKi5qc29uJykpXVxufVxuXG5mdW5jdGlvbiB3cml0ZShvdXRQYXRoOiBzdHJpbmcsIGZpbGVOYW1lOiBzdHJpbmcsIGZpbGU6IHN0cmluZykge1xuICBta2RpcignLXAnLCBvdXRQYXRoKVxuICB3cml0ZUZpbGVTeW5jKGpvaW4ob3V0UGF0aCwgYCR7ZmlsZU5hbWV9LmpzYCksIGZpbGUpXG59XG4iXX0=