# 🔍 Eth2 Beacon Chain Explorer

[![Build](https://github.com/beaconchain-com/eth2-beaconchain-explorer/actions/workflows/build.yaml/badge.svg)](https://github.com/beaconchain-com/eth2-beaconchain-explorer/actions/workflows/build.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/beaconchain-com/eth2-beaconchain-explorer)](https://goreportcard.com/report/github.com/beaconchain-com/eth2-beaconchain-explorer)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](LICENSE)

The **Eth2 Beacon Chain Explorer** provides a comprehensive and easy‑to‑use interface for the Ethereum Beacon Chain. It allows you to view proposed blocks, follow attestations, monitor staking activity, and keep an eye on validator performance.

This project is maintained by **Mahdi Amolimoghaddam** and the open‑source community.

## ✨ Features

- 📊 **Validator dashboard** – real‑time status, rewards, and history  
- 🧱 **Block & slot explorer** – browse blocks, slots, epochs, and committees  
- 📈 **Live statistics** – participation rate, staked ETH distribution, network trends  
- 🔔 **Customisable alerts** – get notified about missed blocks, balance changes, and more  
- 🌐 **Multi‑network support** – Ethereum Mainnet, Gnosis Chain, Holesky, Sepolia  
- 🧩 **Powerful REST API** – real‑time and historical data for developers  

## 📱 Mobile App

Track your validators on the go with the **Beaconchain Dashboard**:

[![Google Play](https://beaconcha.in/img/android.png)](https://play.google.com/store/apps/details?id=in.beaconcha.mobile)
[![App Store](https://beaconcha.in/img/ios.png)](https://apps.apple.com/at/app/beaconchain-dashboard/id1541822121)

## 🏗️ Architecture

| Component | Technology |
| :--- | :--- |
| **Backend** | Go (Golang) |
| **Frontend** | React + TypeScript |
| **Database** | PostgreSQL |
| **Cache** | Redis |
| **Message Queue** | RabbitMQ |
| **Deployment** | Docker & Docker Compose |

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/beaconchain-com/eth2-beaconchain-explorer.git
cd eth2-beaconchain-explorer

# Copy and edit the configuration
cp config/config-example.toml config/config.toml
# Edit config.toml: set your beacon node endpoint, database credentials, etc.

# Start everything with Docker Compose
docker-compose up -d
