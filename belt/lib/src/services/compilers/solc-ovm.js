"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.compileAll = void 0;
const ovm_compiler_1 = require("@krebernisak/ovm-compiler");
const utils_1 = require("../utils");
const solc_1 = require("./solc");
/**
 * Generate solidity artifacts for all of the solidity versions under a specified contract
 * directory.
 *
 * @param conf The application configuration, e.g. where to read solidity files, where to output, etc..
 */
async function compileAll(conf) {
    return Promise.all(utils_1.getContractDirs(conf).map(async ({ dir, version }) => {
        const opts = solc_1.getCompilerOptions(conf, dir, version);
        const c = new ovm_compiler_1.Compiler({
            ...opts,
            // Update version string to be detected by forked 0x/sol-compiler
            solcVersion: opts.solcVersion + '_ovm',
            isOfflineMode: true,
        });
        // Compiler#getCompilerOutputsAsync throws on compilation errors
        // this method prints any errors and warnings for us
        await c.compileAsync();
    }));
}
exports.compileAll = compileAll;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoic29sYy1vdm0uanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi9zcmMvc2VydmljZXMvY29tcGlsZXJzL3NvbGMtb3ZtLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7OztBQUFBLDREQUFvRDtBQUVwRCxvQ0FBMEM7QUFDMUMsaUNBQTJDO0FBRTNDOzs7OztHQUtHO0FBQ0ksS0FBSyxVQUFVLFVBQVUsQ0FBQyxJQUFnQjtJQUMvQyxPQUFPLE9BQU8sQ0FBQyxHQUFHLENBQ2hCLHVCQUFlLENBQUMsSUFBSSxDQUFDLENBQUMsR0FBRyxDQUFDLEtBQUssRUFBRSxFQUFFLEdBQUcsRUFBRSxPQUFPLEVBQUUsRUFBRSxFQUFFO1FBQ25ELE1BQU0sSUFBSSxHQUFHLHlCQUFrQixDQUFDLElBQUksRUFBRSxHQUFHLEVBQUUsT0FBTyxDQUFDLENBQUE7UUFFbkQsTUFBTSxDQUFDLEdBQUcsSUFBSSx1QkFBUSxDQUFDO1lBQ3JCLEdBQUcsSUFBSTtZQUNQLGlFQUFpRTtZQUNqRSxXQUFXLEVBQUUsSUFBSSxDQUFDLFdBQVcsR0FBRyxNQUFNO1lBQ3RDLGFBQWEsRUFBRSxJQUFJO1NBQ3BCLENBQUMsQ0FBQTtRQUVGLGdFQUFnRTtRQUNoRSxvREFBb0Q7UUFDcEQsTUFBTSxDQUFDLENBQUMsWUFBWSxFQUFFLENBQUE7SUFDeEIsQ0FBQyxDQUFDLENBQ0gsQ0FBQTtBQUNILENBQUM7QUFqQkQsZ0NBaUJDIiwic291cmNlc0NvbnRlbnQiOlsiaW1wb3J0IHsgQ29tcGlsZXIgfSBmcm9tICdAa3JlYmVybmlzYWsvb3ZtLWNvbXBpbGVyJ1xuaW1wb3J0ICogYXMgY29uZmlnIGZyb20gJy4uL2NvbmZpZydcbmltcG9ydCB7IGdldENvbnRyYWN0RGlycyB9IGZyb20gJy4uL3V0aWxzJ1xuaW1wb3J0IHsgZ2V0Q29tcGlsZXJPcHRpb25zIH0gZnJvbSAnLi9zb2xjJ1xuXG4vKipcbiAqIEdlbmVyYXRlIHNvbGlkaXR5IGFydGlmYWN0cyBmb3IgYWxsIG9mIHRoZSBzb2xpZGl0eSB2ZXJzaW9ucyB1bmRlciBhIHNwZWNpZmllZCBjb250cmFjdFxuICogZGlyZWN0b3J5LlxuICpcbiAqIEBwYXJhbSBjb25mIFRoZSBhcHBsaWNhdGlvbiBjb25maWd1cmF0aW9uLCBlLmcuIHdoZXJlIHRvIHJlYWQgc29saWRpdHkgZmlsZXMsIHdoZXJlIHRvIG91dHB1dCwgZXRjLi5cbiAqL1xuZXhwb3J0IGFzeW5jIGZ1bmN0aW9uIGNvbXBpbGVBbGwoY29uZjogY29uZmlnLkFwcCkge1xuICByZXR1cm4gUHJvbWlzZS5hbGwoXG4gICAgZ2V0Q29udHJhY3REaXJzKGNvbmYpLm1hcChhc3luYyAoeyBkaXIsIHZlcnNpb24gfSkgPT4ge1xuICAgICAgY29uc3Qgb3B0cyA9IGdldENvbXBpbGVyT3B0aW9ucyhjb25mLCBkaXIsIHZlcnNpb24pXG5cbiAgICAgIGNvbnN0IGMgPSBuZXcgQ29tcGlsZXIoe1xuICAgICAgICAuLi5vcHRzLFxuICAgICAgICAvLyBVcGRhdGUgdmVyc2lvbiBzdHJpbmcgdG8gYmUgZGV0ZWN0ZWQgYnkgZm9ya2VkIDB4L3NvbC1jb21waWxlclxuICAgICAgICBzb2xjVmVyc2lvbjogb3B0cy5zb2xjVmVyc2lvbiArICdfb3ZtJyxcbiAgICAgICAgaXNPZmZsaW5lTW9kZTogdHJ1ZSxcbiAgICAgIH0pXG5cbiAgICAgIC8vIENvbXBpbGVyI2dldENvbXBpbGVyT3V0cHV0c0FzeW5jIHRocm93cyBvbiBjb21waWxhdGlvbiBlcnJvcnNcbiAgICAgIC8vIHRoaXMgbWV0aG9kIHByaW50cyBhbnkgZXJyb3JzIGFuZCB3YXJuaW5ncyBmb3IgdXNcbiAgICAgIGF3YWl0IGMuY29tcGlsZUFzeW5jKClcbiAgICB9KSxcbiAgKVxufVxuIl19