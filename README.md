# EWallet
Test task, based on task for intern in infotecs

# Install and run

1. Clone repository to your device:
```
git clone https://github.com/Eevangelion/ewallet.git
```
2. Run an application using next command:
```
docker-compose up --build
```

# Endpoints

- `POST /api/v1/wallet` - create new wallet with default balance
    - Output contains JSON object with wallet state:
        - id - new wallet id
        - balance - new wallet balance
- `POST api/v1/wallet/:walletId/send` - send money from one wallet to another
    - Input contains JSON object:
        - id - recipient's wallet id
        - amount - transfer amount
- `GET /api/v1/wallet/:walletId/history` - show wallet history 
    - Output contains array of JSON object with incoming and outgoing transactions of the wallet. Each object consists of:
        - time - date and time of transfer in RFC 3339 format
        - from - outgoing wallet ID
        - to - incoming wallet ID
        - amount - transfer amount
- `GET /api/v1/wallet/:walletId` - show wallet state
    - Output contains JSON object:
        - id - chosen wallet id
        - balance - chosen wallet balance