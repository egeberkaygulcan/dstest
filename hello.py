import time
import random

print(f' Replica {random.randint(1, 20)} sleeping...')
start = time.time()

while time.time() - start < 5:
    time.sleep(0.1)
    # print('Still sleeping...')

print('Woke up!')