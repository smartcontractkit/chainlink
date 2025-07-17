# CRE Local Developer Environment - Component Diagram

```mermaid
graph TB
    subgraph "Local Developer Environment"
        subgraph "Control Plane"
            Setup[Env Command]
            Config[TOML Config]
        end
        
        subgraph "Blockchain Layer"
            BC1[Blockchain Node 1<br/>Registry Chain]
        end
        
        subgraph "Job Distribution"
            JD[Job Distributor<br/>Port: Various]
        end
        
        subgraph "Node Topology - Simplified"
            subgraph "Unified DON"
                N1[Node 1<br/>Bootstrap + Gateway <br/> <i>Cron</i>]
                N2[Node 2<br/> <i>Cron</i>]
                N3[Node 3<br/> <i>Cron</i>]
                N4[Node 4<br/> <i>Cron</i>]
            end
        end
        

        
        subgraph "Capabilities"
            Cron[Cron Capability]
        end
        
        subgraph "Workflows"
            WF[User Workflow]
        end

        subgraph "Core Node"
            CoreNode[Core Node]
        end

        subgraph "Observability Stack - Beholder"
            subgraph "Chip Ingress"
                ChipIngress[Chip Ingress<br/>chip-ingress:local-cre]
            end
            
            subgraph "Red Panda"
                RedPanda[Red Panda<br/>Message Streaming]
                Topics[Kafka Topics]
            end
            
            subgraph "Monitoring"
                Metrics[Metrics Collection]
                Logs[Log Aggregation]
            end
        end
        
    end
    
    subgraph "Developer Interface"
        Dev[Developer]
    end
    
    %% Control Flow
    Dev --> Setup
    Config --> Setup
    Setup --> JD
    Setup --> BC1
    Setup --> |"(+ N2, N3, N4)"| N1
    Config --> Cron
    Config --> CoreNode
    
    %% Job Distribution
    JD -->|"(+ N2, N3, N4)"| N1

   
    
    %% Blockchain Connections
    N1 --> BC1
    N2 --> BC1
    N3 --> BC1
    N4 --> BC1

    
    %% Dev Ex
    Dev --> Cron
    Dev --> WF
    Dev --> CoreNode 
    
    
    %% Observability
    N1 --> ChipIngress
    N2 --> ChipIngress
    N3 --> ChipIngress
    N4 --> ChipIngress

    N1 --> Metrics
    N2 --> Metrics
    N3 --> Metrics
    N4 --> Metrics
    
    N1 --> Logs
    N2 --> Logs
    N3 --> Logs
    N4 --> Logs
    
    ChipIngress --> RedPanda
    RedPanda --> Topics
    
    
    
    
    %% Styling
    classDef blockchain fill:#e1f5fe
    classDef node fill:#f3e5f5
    classDef capability fill:#e8f5e8
    classDef observability fill:#fff3e0
    classDef storage fill:#fce4ec
    classDef external fill:#f1f8e9
    
    class BC1 blockchain
    class N1,N2,N3,N4 node
    class Cron capability
    class ChipIngress,RedPanda,Topics,Metrics,Logs observability

```

## Architecture Overview

### **Control Plane**
- **Setup Command**: Handles prerequisite checks and environment initialization
- **TOML Config**: Configuration files defining node sets, capabilities, and infrastructure

### **Node Topologies**

#### **Simplified Topology**
- Single unified DON with all capabilities
- Node 1 acts as both Bootstrap and Gateway
- All nodes share the same capability set



### **Capabilities**
- **Cron**: Scheduled task execution

### **Observability Stack (Beholder)**
- **Chip Ingress**: Data collection and routing service
- **Red Panda**: Kafka-compatible message streaming
- **Monitoring**: Metrics collection and log aggregation


### **Developer Interface**
- **CLI Commands**: Direct environment control
- **Browser UI**: Gateway web interface
- **Terminal**: Development workflow integration

## Key Features

1. **Capability Distribution**: Flexible assignment of capabilities across different node types
2. **Observability Integration**: Built-in monitoring and logging via Beholder stack
3. **Docker-based**: Containerized deployment for consistency and isolation
4. **Developer-friendly**: Simple CLI commands for environment management

### Advanced Features
#### **Full Topology**
- **Workflow DON**: Handles workflow execution (OCR3, Compute, WebAPI, Cron, LogEvent)
- **Capabilities DON**: Manages blockchain interactions (WriteEVM, ReadContract, WebAPI)
- **Gateway DON**: Provides external interface and routing
#### **Capabilities** In Progress...
- **OCR3**: Off-chain reporting consensus mechanism
- **Custom Compute**: Computational workflows
- **WebAPI**: HTTP triggers and targets
- **WriteEVM**: Blockchain write operations
- **Cron**: Scheduled task execution
- **LogEvent**: Blockchain event monitoring
- **ReadContract**: Smart contract reading