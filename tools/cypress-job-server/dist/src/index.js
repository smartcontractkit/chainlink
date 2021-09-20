"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const tslib_1 = require("tslib");
/* eslint-disable @typescript-eslint/no-var-requires */
const net_1 = tslib_1.__importDefault(require("net"));
(async function () {
    const [customResponse] = process.argv.slice(2);
    const defaultResponse = '{"last": "3843.95"}';
    const response = customResponse || defaultResponse;
    const port = process.env.JOB_SERVER_PORT || 6692;
    const server = new net_1.default.Server((socket) => {
        socket.on('data', () => {
            socket.write(`HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: ${response.length}

${response}`);
            socket.end();
        });
    });
    server.on('close', () => {
        server.unref();
    });
    const endpoint = server.listen(port);
    const address = endpoint.address();
    if (address && typeof address != 'string') {
        console.log(`Job Server listening on port ${address.port}`);
    }
    else {
        console.error('Invalid server setup. Address should be of type net.AddressInfo');
        process.exit(1);
    }
})();
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiaW5kZXguanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi9zcmMvaW5kZXgudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7O0FBQUEsdURBQXVEO0FBQ3ZELHNEQUNDO0FBQUEsQ0FBQyxLQUFLO0lBQ0wsTUFBTSxDQUFDLGNBQWMsQ0FBQyxHQUFHLE9BQU8sQ0FBQyxJQUFJLENBQUMsS0FBSyxDQUFDLENBQUMsQ0FBQyxDQUFBO0lBQzlDLE1BQU0sZUFBZSxHQUFHLHFCQUFxQixDQUFBO0lBQzdDLE1BQU0sUUFBUSxHQUFHLGNBQWMsSUFBSSxlQUFlLENBQUE7SUFDbEQsTUFBTSxJQUFJLEdBQUcsT0FBTyxDQUFDLEdBQUcsQ0FBQyxlQUFlLElBQUksSUFBSSxDQUFBO0lBRWhELE1BQU0sTUFBTSxHQUFHLElBQUksYUFBRyxDQUFDLE1BQU0sQ0FBQyxDQUFDLE1BQWtCLEVBQUUsRUFBRTtRQUNuRCxNQUFNLENBQUMsRUFBRSxDQUFDLE1BQU0sRUFBRSxHQUFHLEVBQUU7WUFDckIsTUFBTSxDQUFDLEtBQUssQ0FBQzs7a0JBRUQsUUFBUSxDQUFDLE1BQU07O0VBRS9CLFFBQVEsRUFBRSxDQUFDLENBQUE7WUFDUCxNQUFNLENBQUMsR0FBRyxFQUFFLENBQUE7UUFDZCxDQUFDLENBQUMsQ0FBQTtJQUNKLENBQUMsQ0FBQyxDQUFBO0lBQ0YsTUFBTSxDQUFDLEVBQUUsQ0FBQyxPQUFPLEVBQUUsR0FBRyxFQUFFO1FBQ3RCLE1BQU0sQ0FBQyxLQUFLLEVBQUUsQ0FBQTtJQUNoQixDQUFDLENBQUMsQ0FBQTtJQUVGLE1BQU0sUUFBUSxHQUFHLE1BQU0sQ0FBQyxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUE7SUFDcEMsTUFBTSxPQUFPLEdBQUcsUUFBUSxDQUFDLE9BQU8sRUFBRSxDQUFBO0lBQ2xDLElBQUksT0FBTyxJQUFJLE9BQU8sT0FBTyxJQUFJLFFBQVEsRUFBRTtRQUN6QyxPQUFPLENBQUMsR0FBRyxDQUFDLGdDQUFnQyxPQUFPLENBQUMsSUFBSSxFQUFFLENBQUMsQ0FBQTtLQUM1RDtTQUFNO1FBQ0wsT0FBTyxDQUFDLEtBQUssQ0FDWCxpRUFBaUUsQ0FDbEUsQ0FBQTtRQUNELE9BQU8sQ0FBQyxJQUFJLENBQUMsQ0FBQyxDQUFDLENBQUE7S0FDaEI7QUFDSCxDQUFDLENBQUMsRUFBRSxDQUFBIiwic291cmNlc0NvbnRlbnQiOlsiLyogZXNsaW50LWRpc2FibGUgQHR5cGVzY3JpcHQtZXNsaW50L25vLXZhci1yZXF1aXJlcyAqL1xuaW1wb3J0IG5ldCBmcm9tICduZXQnXG47KGFzeW5jIGZ1bmN0aW9uICgpIHtcbiAgY29uc3QgW2N1c3RvbVJlc3BvbnNlXSA9IHByb2Nlc3MuYXJndi5zbGljZSgyKVxuICBjb25zdCBkZWZhdWx0UmVzcG9uc2UgPSAne1wibGFzdFwiOiBcIjM4NDMuOTVcIn0nXG4gIGNvbnN0IHJlc3BvbnNlID0gY3VzdG9tUmVzcG9uc2UgfHwgZGVmYXVsdFJlc3BvbnNlXG4gIGNvbnN0IHBvcnQgPSBwcm9jZXNzLmVudi5KT0JfU0VSVkVSX1BPUlQgfHwgNjY5MlxuXG4gIGNvbnN0IHNlcnZlciA9IG5ldyBuZXQuU2VydmVyKChzb2NrZXQ6IG5ldC5Tb2NrZXQpID0+IHtcbiAgICBzb2NrZXQub24oJ2RhdGEnLCAoKSA9PiB7XG4gICAgICBzb2NrZXQud3JpdGUoYEhUVFAvMS4xIDIwMCBPS1xuQ29udGVudC1UeXBlOiBhcHBsaWNhdGlvbi9qc29uXG5Db250ZW50LUxlbmd0aDogJHtyZXNwb25zZS5sZW5ndGh9XG5cbiR7cmVzcG9uc2V9YClcbiAgICAgIHNvY2tldC5lbmQoKVxuICAgIH0pXG4gIH0pXG4gIHNlcnZlci5vbignY2xvc2UnLCAoKSA9PiB7XG4gICAgc2VydmVyLnVucmVmKClcbiAgfSlcblxuICBjb25zdCBlbmRwb2ludCA9IHNlcnZlci5saXN0ZW4ocG9ydClcbiAgY29uc3QgYWRkcmVzcyA9IGVuZHBvaW50LmFkZHJlc3MoKVxuICBpZiAoYWRkcmVzcyAmJiB0eXBlb2YgYWRkcmVzcyAhPSAnc3RyaW5nJykge1xuICAgIGNvbnNvbGUubG9nKGBKb2IgU2VydmVyIGxpc3RlbmluZyBvbiBwb3J0ICR7YWRkcmVzcy5wb3J0fWApXG4gIH0gZWxzZSB7XG4gICAgY29uc29sZS5lcnJvcihcbiAgICAgICdJbnZhbGlkIHNlcnZlciBzZXR1cC4gQWRkcmVzcyBzaG91bGQgYmUgb2YgdHlwZSBuZXQuQWRkcmVzc0luZm8nLFxuICAgIClcbiAgICBwcm9jZXNzLmV4aXQoMSlcbiAgfVxufSkoKVxuIl19