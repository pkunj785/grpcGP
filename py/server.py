import logging
from concurrent.futures import ThreadPoolExecutor

import grpc
import numpy as np

from anomaly_pb2 import AnomResponse
from anomaly_pb2_grpc import AnomsServicer, add_AnomsServicer_to_server

def find_outliers(data: np.ndarray):
    """Return indices where values more than 2 standard deviation from mean"""
    out = np.where(np.abs(data - data.mean()) > 2 * data.std())
    #np.where returns tuple for each dimension, we want the 1st element
    return out[0]

class OutliersServer(AnomsServicer):
    def Expose(self, request, context):
        logging.info('detect request size: %d', len(request.metrics))
        # Convert metrics to numpy array of values only
        data = np.fromiter((m.value for m in request.metrics), dtype='float64')
        indices = find_outliers(data)
        logging.info('found %d outliers', len(indices))
        resp = AnomResponse(indices=indices)
        return resp
    
if __name__ == '__main__':
    logging.basicConfig(
        level=logging.INFO,
        format='%(asctime)s - %(levelname)s - %(message)s'
    )
    server = grpc.server(ThreadPoolExecutor())
    add_AnomsServicer_to_server(OutliersServer(), server)
    port = 9091
    server.add_insecure_port(f'[::]:{port}')
    server.start()
    logging.info('server ready on port %r', port)
    server.wait_for_termination()