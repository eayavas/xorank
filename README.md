# XORANK

> **XORank** is a minimalist pairwise comparison engine utilizing the Elo rating system.

## Overview

**XORank** is a project designed to rank items based on user preferences through one-on-one matchups. It functions as a lightweight clone of Facemash, powered by the Elo rating algorithm to ensure mathematical accuracy in ranking.

Originally developed for the [**Essociete**](https://essociete.org) hackerspace, this tool was built to settle a competition finding the "Ultimate Web Browser" among members.

**Etymology:** The name is a portmanteau of **XOR** (Exclusive OR logic gate - *choice between two*) and **Rank**.

## Features

* **Elo Rating System:** Uses a standard logistic distribution (K=32) for dynamic scoring.
* **Passcode Authentication:** Minimalist, username-free security using unique access codes.
* **Cybercore Aesthetic:** A dark, terminal-inspired UI with CSS-only styling (No heavy JS frameworks).
* **Concurrency Safe:** Handles multiple users voting simultaneously via Go routines and SQLite transactions.
* **Session Management:** Cookie-based sessions to track user votes and prevent duplicate entries for the same matchup.
* **Responsive Design:** Optimized for both desktop terminals and mobile viewports.

## Tech Stack

* **Language:** Go (Golang) 1.23+
* **Database:** SQLite3
* **Frontend:** HTML5 / CSS3 (Grid & Flexbox)
* **Architecture:** Standard Go Project Layout (`cmd`, `internal`, `templates`)

## Installation

### Prerequisites
* **Go**: Version 1.23 or higher.
* **C Compiler**: GCC or Clang (Required for `go-sqlite3` CGO).

### Clone the Repository
```bash
git clone [https://github.com/yourusername/xorank.git](https://github.com/yourusername/xorank.git)
cd xorank
```

## Usage

### First Run Setup

If you are running this for the first time, initialize the module and start the server. The database (`elo.db`) will be created automatically.

```bash
# 1. Download dependencies
go mod tidy

# 2. Run the application
go run cmd/xorank/main.go

```

### Accessing the System

Open your browser and navigate to:
`http://localhost:8080`

**Default Access Codes (Passcodes):**

* **Admin:** `admin`
* **Test:** `1234`
* **Guest:** `0000`
* **Lucky:** `7777`

### Resetting Data

To completely reset the rankings, votes, and users, simply delete the database file before restarting the server:

```bash
rm elo.db
go run cmd/xorank/main.go

```

## Configuration

To modify the default **Items** (Browsers) or **Access Codes**, edit the `seedData` function located in:
`internal/store/sqlite.go`

## Contributing

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/AmazingFeature`).
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4. Push to the branch (`git push origin feature/AmazingFeature`).
5. Open a Pull Request.

## License

Distributed under the MIT License.


*Developed by [eayavas*](https://github.com/eayavas) for Essociete

```

```