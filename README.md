<img src="https://github.com/motixo/toman/blob/main/asset/toman.jpg" width="350" alt="toman">

# Toman CLI

Toman CLI is a simple command-line tool to fetch current prices of currencies and gold-coin from TGJU.

## Usage

### 1. Fetch All Prices
Run the application without arguments to get a summary of all configured assets.

```bash
$toman
```

**Output:**

```text
USD         : 50,150 Toman
EUR         : 54,200 Toman
GOLD/COIN   : 29,500,000 Toman
TETHER      : 50,300 Toman
```

### 2. Filter Specific Currencies

Use flags to fetch only the data you need.

```bash
$toman -usd -tether
```

**Output:**

```text
USD         : 50,150 Toman
TETHER      : 50,300 Toman
```

### 3. Help Menu

View all available flags and options.

```bash
$toman -help
```


## Disclaimer

This tool is intended for **educational purposes only**. The data is scraped from public web sources (`tgju.org`).

* Please respect the website's `robots.txt` and Terms of Service.
* Do not use this for high-frequency trading or excessive polling, as it may lead to IP bans.
