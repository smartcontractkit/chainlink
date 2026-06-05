package main

import (
    "fmt"
    "net/url"
    "github.com/smartcontractkit/chainlink-common/pkg/config"
)

func main() {
    u, _ := url.Parse("postgres://postgres:postgres@localhost:5432/db")
    cUrl := (*config.URL)(u)
    sUrl := config.NewSecretURL(cUrl)
    fmt.Printf("sUrl.String(): %s\n", sUrl.String())
    
    // Now convert it back to url.URL
    back := *sUrl.URL()
    fmt.Printf("back.String(): %s\n", back.String())
    
    // What about User.Password()?
    pass, ok := back.User.Password()
    fmt.Printf("Password: %s, ok: %v\n", pass, ok)
}
