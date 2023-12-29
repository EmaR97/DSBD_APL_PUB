import logging
import threading
from concurrent import futures

import grpc

from message import subscription_pb2, subscription_pb2_grpc
from mongo import Subscription


class SubscriptionServiceServicer(subscription_pb2_grpc.SubscriptionService_Servicer):
    def GetChatIds(self, request, context):
        cam_id = request.cam_id
        logging.debug(f'get_chat_ids request: {cam_id}')
        chat_ids = [subscription.id.chat_id for subscription in Subscription.get_by_cam_id(cam_id) if
                    subscription.to_notify()]
        logging.debug(f'get_chat_ids response: {chat_ids}')
        return subscription_pb2.ChatIdsResponse(chat_ids=chat_ids)


def run_grpc_server(address, stop_event: threading.Event):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    subscription_pb2_grpc.add_SubscriptionService_Servicer_to_server(SubscriptionServiceServicer(), server)
    server.add_insecure_port(address)
    server.start()
    stop_event.wait()  # Wait until the stop event is set
    server.stop(grace=None)
    server.wait_for_termination()
