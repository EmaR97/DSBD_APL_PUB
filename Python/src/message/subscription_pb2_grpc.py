# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import subscription_pb2 as subscription__pb2


class SubscriptionService_Stub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.GetChatIds = channel.unary_unary(
                '/message.SubscriptionService_/GetChatIds',
                request_serializer=subscription__pb2.CamIdRequest.SerializeToString,
                response_deserializer=subscription__pb2.ChatIdsResponse.FromString,
                )


class SubscriptionService_Servicer(object):
    """Missing associated documentation comment in .proto file."""

    def GetChatIds(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_SubscriptionService_Servicer_to_server(servicer, server):
    rpc_method_handlers = {
            'GetChatIds': grpc.unary_unary_rpc_method_handler(
                    servicer.GetChatIds,
                    request_deserializer=subscription__pb2.CamIdRequest.FromString,
                    response_serializer=subscription__pb2.ChatIdsResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'message.SubscriptionService_', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class SubscriptionService_(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def GetChatIds(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/message.SubscriptionService_/GetChatIds',
            subscription__pb2.CamIdRequest.SerializeToString,
            subscription__pb2.ChatIdsResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
