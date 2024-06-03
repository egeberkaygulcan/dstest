import requests
import socket
import sys
import time
import random

from requests.adapters import HTTPAdapter
from urllib3 import PoolManager

workerId = int(sys.argv[1])
numReplicas = int(sys.argv[2])
baseInterceptorPort = int(sys.argv[3])
workerBindPort = 6000 + workerId

time.sleep(1)

print("Worker ID: " + str(workerId) +
      "\tNumber of replicas: " + str(numReplicas) +
      "\tBase interceptor port: " + str(baseInterceptorPort))

# blatantly stolen from https://stackoverflow.com/a/47203137 with no regrets
class SourcePortAdapter(HTTPAdapter):
    """"Transport adapter" that allows us to set the source port."""
    def __init__(self, port, *args, **kwargs):
        self._source_port = port
        super(SourcePortAdapter, self).__init__(*args, **kwargs)

    def init_poolmanager(self, connections, maxsize, block=False):
        self.poolmanager = PoolManager(
            num_pools=connections, maxsize=maxsize,
            block=block, source_address=('', self._source_port))

# TCP client
# send 2 tcp requests to random interceptors and then exit

time.sleep(workerId)
for i in range(2):
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        s.bind(("127.0.0.1", workerBindPort))
        randomPort = baseInterceptorPort + random.randrange(0, numReplicas)
        print("Worker " + str(workerId) + " sending request to port: " + str(randomPort))
        s.connect(('127.0.0.1', randomPort))
        s.sendall(b'Hello, world')
        #data = s.recvmsg(1024)
        #print('Worker ' + str(workerId) + ' received', repr(data))
        s.close()
        time.sleep(numReplicas)
'''

# HTTP client
# send 2 http requests to random interceptors and then exit
time.sleep(workerId)
s = requests.Session()
s.mount('http://', SourcePortAdapter(workerBindPort))
for i in range(2):
    randomPort = baseInterceptorPort + random.randrange(0, numReplicas)
    path = "http://localhost:" + str(randomPort) + "/hello/from/worker/" + str(workerId)
    print("Worker " + str(workerId) + " sending request to path: " + path)
    s.get(path)
s.close()
'''
