# ⚡ Minimal BitTorrent Client (Go)

A minimal BitTorrent client written in Go 🚀 that demonstrates core peer-to-peer file sharing concepts like handshake, piece exchange, and integrity verification.

---

## 📦 Features

- 🤝 Peer handshake & connection
- 📡 BitTorrent wire protocol support
- 🧩 Piece-based downloading (16KB blocks)
- 🔐 SHA-1 integrity verification
- ⚙️ Concurrent workers (multi-peer support)
- 💾 Random-access file writing
- 📊 Download progress tracking
- 🧪 CLI flags for flexible usage

---

## 👨‍💻 Usage

```bash
go run cmd/main.go tmp/sample.torrent tmp/out.md
```

---

## 🏗️ Architecture

```mermaid
flowchart TD

    A[Torrent File] --> B[Metadata Parser]
    B --> C[Task Generator]

    C -->|piece tasks| D[Task Queue]

    subgraph Peer Workers
        E1[Peer 1]
        E2[Peer 2]
        E3[Peer N]
    end

    D --> E1
    D --> E2
    D --> E3

    subgraph Download Flow
        F1[Handshake]
        F2[Bitfield Exchange]
        F3[Interested]
        F4[Request Blocks]
        F5[Receive Blocks]
        F6[Assemble Piece]
        F7[Verify SHA-1]
    end

    E1 --> F1 --> F2 --> F3 --> F4 --> F5 --> F6 --> F7
    E2 --> F1
    E3 --> F1

    F7 -->|valid| G[Piece Results]
    F7 -->|invalid| D

    G --> H[File Writer]
    H --> I[Output File]
```
