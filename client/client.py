import grpc
import server_pb2
import server_pb2_grpc
import google.protobuf.empty_pb2; 

def run():
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = server_pb2_grpc.ServerStub(channel)
        response = stub.AddNumber(server_pb2.AddNumberReq(a=1, b=2))
        print("Client received: " + str(response.a))
        response = stub.ReadFromDB(google.protobuf.empty_pb2.Empty())
        print("Client received: " + str(response.a))

if __name__ == '__main__':
    run()