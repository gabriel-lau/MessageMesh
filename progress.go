package main

// Used to track the steps and progress of the app
// Steps:
// Start the app
// Display the loading screen
// Connect to the network
// Join the PubSub
// Start the consensus
// Retrieve the blockchain

// Code here
type Progress struct {
	AppStarted       chan bool
	LoadingDisplayed chan bool
	NetworkConnected chan bool
	PubSubJoined     chan bool
	ConsensusStarted chan bool
	BlockchainLoaded chan bool
}

// NewProgress creates a new Progress tracker with buffered channels
func NewProgress() *Progress {
	return &Progress{
		AppStarted:       make(chan bool, 1),
		LoadingDisplayed: make(chan bool, 1),
		NetworkConnected: make(chan bool, 1),
		PubSubJoined:     make(chan bool, 1),
		ConsensusStarted: make(chan bool, 1),
		BlockchainLoaded: make(chan bool, 1),
	}
}

// GetCurrentStep returns the current step in the startup process
// This should be called in a separate goroutine
func (p *Progress) GetCurrentStep() string {
	select {
	case <-p.AppStarted:
		select {
		case <-p.LoadingDisplayed:
			select {
			case <-p.NetworkConnected:
				select {
				case <-p.PubSubJoined:
					select {
					case <-p.ConsensusStarted:
						select {
						case <-p.BlockchainLoaded:
							return "Application ready"
						default:
							return "Loading blockchain data..."
						}
					default:
						return "Starting consensus protocol..."
					}
				default:
					return "Joining peer-to-peer network..."
				}
			default:
				return "Connecting to network..."
			}
		default:
			return "Displaying loading screen..."
		}
	default:
		return "Starting application..."
	}
}

// CompleteStep marks a step as complete by sending true to its channel
func (p *Progress) CompleteStep(step string) {
	switch step {
	case "app":
		p.AppStarted <- true
	case "loading":
		p.LoadingDisplayed <- true
	case "network":
		p.NetworkConnected <- true
	case "pubsub":
		p.PubSubJoined <- true
	case "consensus":
		p.ConsensusStarted <- true
	case "blockchain":
		p.BlockchainLoaded <- true
	}
}

// WaitForCompletion blocks until all steps are completed
func (p *Progress) WaitForCompletion() {
	<-p.AppStarted
	<-p.LoadingDisplayed
	<-p.NetworkConnected
	<-p.PubSubJoined
	<-p.ConsensusStarted
	<-p.BlockchainLoaded
}
