```mermaid
sequenceDiagram
    title MessageMesh Startup Sequence

    %% Participants
    participant App
    participant Network
    participant SystemMonitor
    participant P2PService
    participant PubSubService
    participant ConsensusService
    participant UIDataLoop

    %% Startup Function
    App->>App: startup(ctx)
    App->>App: Save context
    App->>Network: ConnectToNetwork()
    
    %% ConnectToNetwork Function
    Network->>SystemMonitor: NewSystemMonitor()
    SystemMonitor-->>Network: Return monitor
    Network->>Network: Start runMonitoring goroutine
    
    %% P2P Setup
    Network->>P2PService: NewP2PService()
    P2PService->>P2PService: setupHost()
    P2PService->>P2PService: bootstrapDHT()
    P2PService-->>Network: Return P2PService
    
    %% Connect to peers
    Network->>P2PService: AdvertiseConnect()
    P2PService->>P2PService: Advertise service
    P2PService->>P2PService: Find peers
    P2PService->>P2PService: Start peer connection handler
    P2PService-->>Network: Connected to peers
    
    %% Join PubSub
    Network->>PubSubService: JoinPubSub(P2PService)
    PubSubService->>PubSubService: Join topic
    PubSubService->>PubSubService: Subscribe to topic
    PubSubService->>PubSubService: Start SubLoop goroutine
    PubSubService->>PubSubService: Start PubLoop goroutine
    PubSubService->>PubSubService: Start PeerJoinedLoop goroutine
    PubSubService-->>Network: Return PubSubService
    
    %% Wait for network setup
    Network->>Network: Sleep for 5 seconds
    
    %% Start consensus
    Network->>ConsensusService: StartConsensus(network)
    ConsensusService->>ConsensusService: Initialize blockchain with genesis block
    ConsensusService->>ConsensusService: Setup Raft consensus
    ConsensusService->>ConsensusService: Start networkLoop goroutine
    ConsensusService->>ConsensusService: Start blockchainLoop goroutine
    ConsensusService-->>Network: Return ConsensusService
    
    %% Start UI Data Loop
    App->>UIDataLoop: Start UIDataLoop goroutine
    UIDataLoop->>UIDataLoop: Emit initial events (getUserPeerID, getPeerList)
    UIDataLoop->>UIDataLoop: Start event loop
    
    %% Completion
    Note over App,UIDataLoop: Startup Complete 