import logging

import grpc

from message import subscription_pb2_grpc, subscription_pb2, cam_pb2_grpc, cam_pb2


def get_chat_ids(cam_id, target):
    with grpc.insecure_channel(target) as channel:
        stub = subscription_pb2_grpc.SubscriptionService_Stub(channel)
        logging.debug(f'get_chat_ids request: {cam_id}')
        request = subscription_pb2.CamIdRequest(cam_id=cam_id)
        response = stub.GetChatIds(request)
        logging.debug(f'get_chat_ids response: {response}')
        return response.chat_ids


def get_cam_ids(user_id, target):
    with grpc.insecure_channel(target) as channel:
        stub = cam_pb2_grpc.SubscriptionServiceStub(channel)
        request = cam_pb2.UserIdRequest(user_id=user_id)
        response = stub.GetCamIds(request)
        return response.cam_ids
