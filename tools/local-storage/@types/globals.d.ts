declare namespace NodeJS {
  interface Global {
    localStorage: {
      clear: Function
      getItem: Function
      setItem: Function
    }
  }
}
