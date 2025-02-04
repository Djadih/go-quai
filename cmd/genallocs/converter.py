import csv
import json

data = []
with open('cmd/genallocs/V1_Genesis_Block_FINAL.csv', newline='') as csvfile:
    reader = csv.DictReader(csvfile)
    for row in reader:
        for field in ['vested', 'award', 'unlockSchedule', 'lumpSumMonth']:
            row[field] = int(row[field]) if row[field].strip() else 0
        data.append(row)

with open('cmd/genallocs/genesis_alloc.json', 'w') as jsonfile:
    json.dump(data, jsonfile, separators=(',', ':'))

