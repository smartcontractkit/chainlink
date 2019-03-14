import "reflect-metadata"
import { createConnection } from "typeorm"
import seed from "./seed"
import server from "./server"

createConnection().then(async dbConnection => {
  seed(dbConnection)
  server(dbConnection)
}).catch(error => console.log(error))
