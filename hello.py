import socket
import sys
import time
import random

workerId = int(sys.argv[1])
numReplicas = int(sys.argv[2])
baseInterceptorPort = int(sys.argv[3])
workerBindPort = 6000 + workerId

print("Worker ID: " + str(workerId) +
      "\tNumber of replicas: " + str(numReplicas) +
      "\tBase interceptor port: " + str(baseInterceptorPort))

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
