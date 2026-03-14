# CLAUDE.md — MessageMesh

## Overview

MessageMesh is a decentralized peer-to-peer (P2P) desktop messaging application. Peers discover each other via a Kademlia DHT, communicate over libp2p pub/sub, and maintain a consistent shared blockchain via the Raft consensus algorithm. All direct messages are AES-256 encrypted, with key exchange protected by RSA key pairs. The desktop UI is built with Svelte and exposed through the Wails framework.

**Key features:**
- Fully decentralized — no central server
- End-to-end encrypted direct messages
- Blockchain-backed message history (Raft consensus)
- Desktop GUI (Wails) with optional headless/server mode
- Real-time peer discovery and presence tracking

---

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                   Wails Desktop App                 │
│  ┌─────────────────────┐  ┌────────────────────┐   │
│  │  Svelte/TS Frontend │◄─►│   Go Backend       │   │
│  │  (frontend/src/)    │  │  (app.go, backend/) │   │
│  └─────────────────────┘  └────────┬───────────┘   │
└───────────────────────────────────-│────────────────┘
                                     │
            ┌────────────────────────┼───────────────┐
            ▼                        ▼               ▼
     ┌─────────────┐        ┌──────────────┐  ┌──────────────┐
     │  P2P Layer  │        │  Consensus   │  │  Key Store   │
     │  (libp2p)   │        │  (Raft)      │  │  (BoltDB)    │
     │  p2p.go     │        │  consensus.go│  │  keys.go     │
     │  pubSub.go  │        │  blockchain  │  │  keys.db     │
     └─────────────┘        └──────────────┘  └──────────────┘
```

**Layers:**
1. **P2P Network** — libp2p host with TCP, TLS, Yamux, Kad-DHT for discovery, pub/sub on topic `"messagemesh"`
2. **Consensus** — Raft FSM maintains a blockchain of message/account/key-exchange blocks shared across peers
3. **Encryption** — RSA key pairs identify peers; AES-256 symmetric keys exchanged via `FirstMessage` blocks
4. **Desktop UI** — Wails bridges Go exported functions to the Svelte frontend via auto-generated JS bindings
5. **Monitoring** — `monitoring/monitor.go` collects CPU, memory, network, and goroutine metrics

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language (backend) | Go 1.23.2 |
| Desktop framework | Wails v2.9.2 |
| P2P networking | go-libp2p v0.15.1, go-libp2p-pubsub v0.5.6, go-libp2p-kad-dht v0.15.0 |
| Consensus | go-libp2p-raft v0.1.8, hashicorp/raft v1.1.1 |
| Key-value store | BoltDB (go.etcd.io/bbolt v1.3.11) |
| System metrics | gopsutil v3.24.5 |
| Env config | godotenv v1.5.1 |
| UI framework | Svelte v5.14.4 + TypeScript v5.7.2 |
| Bundler | Vite v6.0.4 |
| CSS | Tailwind CSS v3.4.9 |
| UI components | Flowbite-Svelte v0.47.4 |
| Package manager | pnpm v9.15.0 |
| Formatter | Prettier v4+ |
| Containerization | Docker, Docker Compose |
| CI/CD | GitHub Actions |

---

## Directory Structure

```
MessageMesh/
├── app.go                  # Wails App struct; all functions exported to frontend
├── main.go                 # Entry point; initializes Wails in GUI or headless mode
├── wails.json              # Wails config (app name, frontend commands, author)
├── go.mod / go.sum         # Go module and dependency checksums
├── Dockerfile              # Container build definition
├── docker-compose.yaml     # Multi-container orchestration
├── .env                    # Local environment vars (gitignored)
├── .prettierrc             # Prettier formatting config
│
├── backend/
│   ├── interface.go        # Service interface types: Network, P2PService, PubSubService, ConsensusService
│   ├── p2p.go              # libp2p host creation (TCP, TLS, Yamux, Kad-DHT)
│   ├── pubSub.go           # Pub/sub topic subscription and message routing
│   ├── server.go           # High-level network orchestration (ConnectToNetwork, Send*)
│   ├── consensus.go        # Raft FSM: ADD_MESSAGE_BLOCK, ADD_ACCOUNT_BLOCK, ADD_FIRST_MESSAGE_BLOCK
│   ├── keys.go             # RSA key generation, BoltDB persistence, AES-256 encrypt/decrypt
│   ├── ui.go               # UIDataLoop: emits Wails events to the frontend
│   └── models/
│       ├── blockchain.go   # Blockchain struct with SHA256 block hashing and validation
│       ├── message.go      # Message model (sender, receiver, message, timestamp)
│       ├── account.go      # Account model (username, public key)
│       └── firstMessage.go # Key-exchange model (symmetric key, sender/receiver)
│
├── debug/
│   ├── log.go              # Colored, timestamped console logging by module
│   └── env.go              # Loads .env file via godotenv
│
├── monitoring/
│   └── monitor.go          # Collects CPU/memory/network/goroutine stats
│
├── frontend/
│   ├── src/
│   │   ├── main.ts         # Svelte app entry point
│   │   ├── App.svelte      # Root component; manages chat state, peer selection, event handling
│   │   ├── app.css         # Global styles
│   │   └── components/
│   │       ├── ChatComponent.svelte          # Message display, input, latency tracking
│   │       ├── ChatListComponent.svelte      # Peer/user list
│   │       └── NavigationRailComponent.svelte # Side navigation
│   ├── wailsjs/            # Auto-generated Wails JS bindings (do not edit manually)
│   │   └── go/main/        # Mirrors exported Go functions in app.go
│   ├── package.json        # Frontend scripts and dependencies
│   ├── vite.config.ts      # Vite bundler config
│   ├── svelte.config.js    # Svelte preprocessor config
│   ├── tailwind.config.ts  # Tailwind theme (custom green palette) + Flowbite plugin
│   └── tsconfig.json       # TypeScript config (ESNext, checkJs: true)
│
├── build/
│   ├── bin/                # Compiled binaries output (gitignored)
│   ├── windows/            # Windows build resources
│   ├── darwin/             # macOS build resources
│   └── appicon.png         # Application icon
│
├── metrics/
│   ├── response_time.csv   # Measured message round-trip latency data
│   └── response_time.ipynb # Jupyter notebook for latency analysis
│
├── diagrams/               # Architecture diagrams
│
└── .github/workflows/
    └── build.yml           # CI: builds Windows + macOS binaries on push/PR to main
```

---

## Development Workflow

### Prerequisites

- Go 1.23+
- Node.js 20+
- pnpm 9.15+ (`npm install -g pnpm`)
- Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Running in Development

```bash
# Start dev mode (Go backend + frontend with HMR)
wails dev
```

This runs `pnpm install` and `pnpm run dev` automatically for the frontend.

### Building for Production

```bash
# Build for current platform
wails build

# Build for specific platform
wails build -platform=windows/amd64
wails build -platform=darwin/amd64

# Output: build/bin/MessageMesh (or .exe on Windows)
```

### Frontend-only Commands

```bash
cd frontend
pnpm install        # Install dependencies
pnpm run dev        # Start Vite dev server
pnpm run build      # Build frontend assets to dist/
pnpm run check      # Run svelte-check (type checking)
```

### Docker

```bash
# Build image
docker build -t gabriellh/messagemesh .

# Run container (headless mode)
docker run -d -p 8080:8080 gabriellh/messagemesh

# Multi-node via compose
docker-compose up
```

---

## Environment Configuration

Create a `.env` file in the repository root (gitignored):

```env
HEADLESS=false     # Set to 'true' to run as a server without GUI
USERNAME=yourname  # Display name in the P2P network
```

Loaded at startup by `debug/env.go` via `godotenv`.

---

## Key Source Files

### `app.go`
The central Go/frontend bridge. Defines `App` struct with all methods exported to Wails:
- `GetPeerList()` — returns known peers
- `GetUserPeerID()` — returns own peer ID
- `SendEncryptedMessage(receiver, message)` — sends AES-encrypted message
- `GetDecryptedMessage(sender, encryptedMsg)` — decrypts incoming message
- `GetMessagesFromPeer(peerID)` — retrieves message history from blockchain

### `backend/interface.go`
Defines the four core service interfaces that the entire backend is built around:
- `Network` — orchestrates all services
- `P2PService` — manages libp2p host, DHT, pub/sub
- `PubSubService` — manages topic and message/peer channels
- `ConsensusService` — manages Raft state machine

### `backend/server.go`
High-level network logic: connects to the DHT, sends messages, handles encryption handshakes.

### `backend/consensus.go`
Raft FSM implementation. Three operation types written to the shared blockchain:
- `ADD_MESSAGE_BLOCK` — stores an encrypted direct message
- `ADD_ACCOUNT_BLOCK` — registers a peer's username and public key
- `ADD_FIRST_MESSAGE_BLOCK` — exchanges AES symmetric keys for E2E encryption

### `backend/ui.go`
`UIDataLoop` runs in a goroutine and periodically emits Wails events (`peer:list`, `blockchain:update`, `connected`) to the frontend to drive reactive UI updates.

---

## Backend Conventions (Go)

### Package Structure
- Domain logic lives in `backend/` with models in `backend/models/`
- Utilities in `debug/` (logging, env) and `monitoring/`
- `app.go` is the thin exported API layer; heavy logic delegates to `backend/`

### Interface-Driven Design
All major services are defined as interfaces in `backend/interface.go`. Implementations are in separate files (`p2p.go`, `pubSub.go`, etc.). Always program against interfaces.

### Message Envelope Pattern
Messages published over pub/sub are JSON with a `type` field for polymorphic dispatch:
```go
type Envelope struct {
    Type    string          `json:"type"`    // "Message", "FirstMessage", "Account"
    Payload json.RawMessage `json:"payload"`
}
```

### Channel-Based Concurrency
`PubSubService` uses Go channels for async message routing:
- `Inbound` — messages received from the network
- `Outbound` — messages to publish to the network
- `PeerJoin` / `PeerLeave` — peer lifecycle events
- `PeerIDs` — current peer list

### Error Handling
Errors are typically logged via `debug/log.go` colored logging. For non-critical paths, errors may be logged and execution continues. Be explicit about error propagation in new code.

### Storage
BoltDB stores RSA key pairs at `db/keys.db` (gitignored). Do not introduce additional database dependencies; use BoltDB for any new persistent key-value needs.

---

## Frontend Conventions (Svelte/TypeScript)

### Svelte 5 Runes
Use Svelte 5 rune syntax throughout — not the legacy `$:` reactive declarations:
```svelte
let count = $state(0)
let { value = $bindable() } = $props()
$effect(() => { /* side effects */ })
```

### Wails Backend Calls
Import from the auto-generated bindings in `wailsjs/go/main/`. Never call Go functions directly.
```typescript
import { SendEncryptedMessage, GetPeerList } from '../wailsjs/go/main/App'

const peers = await GetPeerList()
```

**Never manually edit files in `frontend/wailsjs/`** — they are regenerated by Wails on each build.

### Real-Time Events
Use Wails event system for push updates from the backend:
```typescript
import { EventsOn } from '../wailsjs/runtime/runtime'

EventsOn('peer:list', (peers) => { /* handle update */ })
EventsOn('blockchain:update', (blocks) => { /* handle update */ })
```

### Styling
- Use Tailwind utility classes for all styling
- Use Flowbite-Svelte components for buttons, modals, inputs, etc.
- The custom Tailwind theme uses a green color palette (defined in `tailwind.config.ts`)
- Dark mode is supported via Tailwind `class` strategy

### Prettier Formatting
Configured in `.prettierrc`:
- Line width: **160 characters**
- Indent: **4 spaces** (tabs: false)
- Quotes: **single**
- Trailing commas: **disabled**
- Run `pnpm run check` before committing to catch type errors

---

## Testing

There are no unit test files in this repository. Quality assurance relies on:

- **`pnpm run check`** — `svelte-check` for TypeScript/Svelte type checking (run before every commit)
- **Manual testing** — spin up multiple instances (or Docker containers) to test P2P behavior
- **Performance metrics** — message round-trip latency logged to `metrics/response_time.csv`, analyzed in `metrics/response_time.ipynb`

When adding new features, manually verify:
1. Peers discover each other and appear in the peer list
2. Messages are sent, received, decrypted, and appear in the blockchain
3. App works in both GUI mode (`HEADLESS=false`) and headless mode (`HEADLESS=true`)

---

## CI/CD

GitHub Actions workflow at `.github/workflows/build.yml`:
- **Trigger:** push or PR to `main`
- **Matrix:** Windows (`windows-latest`), macOS (`macos-latest`)
- **Steps:** checkout → setup Go 1.21 → setup Node 20 → install pnpm 9.15 → install Wails → `wails build`
- **Artifacts:** compiled binaries uploaded per platform

---

## Git Workflow

- **Main branch:** `main` (protected; CI runs on push/PR)
- **Feature branches:** use descriptive names, e.g. `feature/encryption-improvements`
- **Commit style** (from project history): short imperative sentences, e.g.:
  - `Added signing`
  - `Prevent empty messages`
  - `Online and offline status and name hover`
- Keep commits focused; avoid large unrelated changes in a single commit
