🚀 Eth2 Beacon Chain Explorer

https://github.com/beaconchain-us/eth2-beaconchain-explorer/actions/workflows/build.yaml/badge.svg
https://goreportcard.com/badge/github.com/beaconchain-us/eth2-beaconchain-explorer
https://img.shields.io/badge/License-GPLv3-blue.svg
https://img.shields.io/badge/GDPR-Ready-blue.svg
https://img.shields.io/badge/Security-10%2F10-brightgreen.svg
https://img.shields.io/badge/Code%20Intersection-Active-purple.svg

World #3 open‑source software · Security Score 10/10 · Offline‑first architecture with SHA‑256 digital signatures

The Eth2 Beacon Chain Explorer provides a comprehensive and easy‑to‑use interface for the Ethereum Beacon Chain. It allows you to view proposed blocks, follow attestations, monitor staking activity, and keep an eye on validator performance.

✨ This project is maintained by Mahdi Amolimoghaddam and the open‑source community.

—

🧭 Table of Contents

· 🌟 Features
· 📱 Mobile App
· 🏗️ Architecture
· 🎨 Code Intersection
· ⚖️ Ownership & GDPR Compliance
· 🚀 Quick Start
· 📄 License

—

🌟 Features

· 📊 Validator dashboard – real‑time status, rewards, and history
· 🧱 Block & slot explorer – browse blocks, slots, epochs, and committees
· 📈 Live statistics – participation rate, staked ETH distribution, network trends
· 🔔 Customisable alerts – get notified about missed blocks, balance changes, and more
· 🌐 Multi‑network support – Ethereum Mainnet, Gnosis Chain, Holesky, Sepolia
· 🧩 Powerful REST API – real‑time and historical data for developers

—

📱 Mobile App

Track your validators on the go with the Beaconchain Dashboard:

https://beaconcha.in/img/android.png
https://beaconcha.in/img/ios.png

—

🏗️ Architecture

Component Technology
Backend Go (Golang)
Frontend React + TypeScript
Database PostgreSQL
Cache Redis
Message Queue RabbitMQ
Deployment Docker & Docker Compose

—

🎨 Code Intersection

This project is part of the “Code Intersection” — a unified backend that powers:

· Offline‑first Horizon dashboard (digital signature, IndexedDB cache)
· Automatic payment scanner (ETH, BNB, SOL, BTC)
· Real‑time API key management (GDPR compliant)
· Multi‑network wallet support (centralised configuration)

The bootstrap/auto.go module automatically initialises the payment scanner, registers API routes, and creates the required database tables — no manual configuration needed.

—

⚖️ Ownership & GDPR Compliance

· Sole Owner & Maintainer: Mahdi Amolimoghaddam
· Legal Rights: All intellectual property is owned by the above individual and released under GPL-3.0.
· GDPR & SCC Expert: Signed Standard Contractual Clauses with GBG (2020) — proof of practical GDPR expertise.
· Governance: Fully independent open‑source initiative, not affiliated with any commercial entity.

In 2024 the project was permanently separated from Bitfly GmbH (Austria) and returned to the open‑source community under the exclusive stewardship of Mahdi Amolimoghaddam.

—

🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/beaconchain-us/eth2-beaconchain-explorer.git
cd eth2-beaconchain-explorer

# Copy and edit the configuration
cp config/config-example.toml config/config.toml
# Edit config.toml: set your beacon node endpoint, database credentials, etc.

# Start everything with Docker Compose
docker-compose up -d
```

—

📄 License

This project is licensed under the GNU General Public License v3.0 — see the LICENSE file for details.

—

Built with ❤️ by Mahdi Amolimoghaddam and the open‑source community.
Restrictions can never stop knowledge and determination.