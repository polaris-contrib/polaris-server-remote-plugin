import sys
import time
from concurrent import futures

import grpc

import plugin_pb2
import plugin_pb2_grpc


class PluginServicer(plugin_pb2_grpc.PluginServicer):
    def Call(self, request, context):
        return plugin_pb2.Response()

    def Ping(self, request, context):
        return plugin_pb2.PongResponse()


def serve():
    # Start the server
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    plugin_pb2_grpc.add_PluginServicer_to_server(PluginServicer(), server)
    server.add_insecure_port('127.0.0.1:8982')
    server.start()

    # Output information
    print("1|1|tcp|127.0.0.1:8982|grpc")
    sys.stdout.flush()

    try:
        while True:
            time.sleep(60 * 60 * 24)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == '__main__':
    serve()
