import json
import requests

addresses = [
    "0x00039Fd4C460BCA501E7d43D186Aa84568A551E6",
    "0x001326c9884333dDd5f1B5EFe006aE045eE4E20E",
    "0x003751638C9aB1b4f4fE76275994a57e4F630D66",
    "0x00686b7b4a1A1f3F3897CDA9c485821dccD2919e",
    "0x0073e18d85BB9C05DC152606bAd7CA09025Db353",
    "0x0074C8330367d14A35295d5d62277260ac2E0bE0"
]

def check_balance(response):
    balance = response.get("result", {})
    return balance

def get_balance(address):
    url = "http://34.136.182.12:9200"
    
    payload = {
        "jsonrpc": "2.0",
        "method": "quai_getBalance",
        "params": [address, "latest"],
        "id": 1
    }

    headers = {
        "Content-Type": "application/json"
    }

    response = requests.post(url, headers=headers, json=payload)

    if response.status_code == 200:
        response_data = response.json()
        return response_data
    else:
        print(f"Error: {response.status_code}, {response.text}")
        return None

total_balance = 0

for addr in addresses:
    # block_number = input("Enter block number (e.g., 0x1fa): ").strip()
    response = get_balance(addr)

    balance = check_balance(response)

    print("addr", addr, "balance", int(balance, 16))

    total_balance += int(balance, 16)


print("total balance", total_balance)