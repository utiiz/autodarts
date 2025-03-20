import websocket
import time
import json
import threading


class InitData:
    def __init__(self, uuid):
        self.uuid = uuid

    def to_json(self):
        return {
            "uuid": self.uuid,
        }


class Message:
    def __init__(self, type, data):
        self.type = type
        self.data = data

    def to_json(self):
        return {
            "type": self.type,
            "data": self.data.to_json()
        }


class WebSocket:
    def __init__(self, UUID, on_message=None):
        self.UUID = UUID
        self.on_message = on_message
        self.ws = websocket.WebSocketApp("ws://127.0.0.1:8090/ws",
                                         on_open=self.on_open,
                                         on_message=self.on_message,
                                         on_error=self.on_error,
                                         on_close=self.on_close)

    def start(self):
        self.ws.run_forever()

    def on_error(self, ws, error):
        print(f"Error occurred: {error}")

    def on_open(self, ws):
        print("Connection successfully opened")

        init_data = InitData(self.UUID)
        init_message = Message("INIT", init_data)
        print(json.dumps(init_message.to_json()))
        ws.send(json.dumps(init_message.to_json()))

    def on_close(self, ws, close_status_code, close_msg):
        print(f"Connection closed: {close_status_code} - {close_msg}")


if __name__ == "__main__":
    print("Starting program")
    ws = WebSocket()
    ws_thread = threading.Thread(target=ws.start, daemon=True)
    ws_thread.start()

    # Keep main thread alive
    while True:
        time.sleep(5)
