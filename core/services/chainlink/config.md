# Configuration Transition

- subgraph names are packages
- thick lines indicate control flow
- dotted lines indicate implicit interface implementation
- regular w/ dot indicate implementation types

```mermaid
flowchart LR

    subgraph cmd
    
        subgraph cmd/app
            NewApp([func NewApp])     
        end
    
        cli>$ chainlink node start]
        
        RunNode([func Client.RunNode])
        
        NewApplication([func NewApplication])
            
        cli == 1. Before ==> NewApp
        cli == 2. Action ==> RunNode
        RunNode ==> NewApplication
    
    end
    
    toml{{TOML?}}    
    
    subgraph services/chainlink
        
        Config[[Config]]
        
        NewTOMLGeneralConfig([func NewTOMLGeneralConfig])
        
        generalConfig --o Config
        
        NewTOMLGeneralConfig --> generalConfig
       
    end
    
    subgraph config
    
        BasicConfig(BasicConfig)
        
        NewGeneralConfig([func NewGeneralConfig])
        
        generalConfig2[generalConfig]
        
        NewGeneralConfig --> generalConfig2
    
        subgraph config/v2
        
            Core[[Core]]
        
        end
    
    end
    
    Config --o Core
    
    NewApp ==> toml
    toml == yes ==> NewTOMLGeneralConfig
    toml == no ==> NewGeneralConfig
    generalConfig -.-> BasicConfig
    generalConfig2 -.-> BasicConfig
    
    
    subgraph chains/evm
    
        LoadChainSet([func LoadChainSet])
        tomlChain{{TOML?}}
        LoadChainSet ==> tomlChain
    
        subgraph chains/evm/config  
        
            NewChainScopedConfig([func NewChainScopedConfig])
        
            ChainScopedOnlyConfig(ChainScopedOnlyConfig)
            
            chainScopedConfig
            
            NewChainScopedConfig --> chainScopedConfig
            
            chainScopedConfig -.-> ChainScopedOnlyConfig
            
            subgraph chains/evm/config/v2 
            
                NewTOMLChainScopedConfig([func NewTOMLChainScopedConfig])
            
                ChainScoped
                
                NewTOMLChainScopedConfig --> ChainScoped
                
                ChainScoped -.-> ChainScopedOnlyConfig
                
                EVMConfig[[EVMConfig]]
                
            end
        
        end
        
        tomlChain == no ==>NewChainScopedConfig
        tomlChain == yes ==>NewTOMLChainScopedConfig
        Config --o EVMConfig
    end
    
    chainScopedConfig --o generalConfig
    ChainScoped --o generalConfig2 
    
    NewApplication ==> LoadChainSet
    
```
