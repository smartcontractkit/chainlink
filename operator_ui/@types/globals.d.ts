declare namespace NodeJS {
  interface Document {
    cookie: string
  }

  interface Global {
    document: Document
    fetch: any
  }
}
