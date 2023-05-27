import os
import logging
from concurrent import futures

import grpc
import psycopg2

import server_pb2
import server_pb2_grpc


class Server(server_pb2_grpc.ServerServicer):
    def __init__(self, conn: psycopg2.extensions.connection):
        self.db = conn
        self.logger = logging.getLogger('server')
        self.logger.info('Server initialized')

    def AddNumber(self, request, context):
        self.logger.info('AddNumber')
        return server_pb2.AddNumberRes(a=request.a + request.b)

    def ReadFromDB(self, request, context):
        self.logger.info('ReadFromDB')
        with self.db.cursor() as cursor:
            cursor.execute("SELECT 1")
            result = cursor.fetchone()[0]
        return server_pb2.ReadFromDBRes(a=result)


def serve():
    logging.basicConfig(level=logging.INFO)

    # Database connection
    conn = psycopg2.connect(os.environ["DATABASE_URL"])

    # Server setup
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    server_pb2_grpc.add_ServerServicer_to_server(Server(conn), server)

    # Server port
    port = '50052'
    server.add_insecure_port('[::]:' + port)
    server.start()

    logging.info("Server started, listening on port %s", port)

    # Keep the server running
    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        server.stop(0)
        logging.info("Server stopped")


if __name__ == '__main__':
    serve()