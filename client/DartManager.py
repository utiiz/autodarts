import cv2
import json
import websockets
from threading import Thread
from DartDetector import DartDetector
from Triangulator import Triangulator, Camera
from Dartboard import Dartboard


class Dart:
    def __init__(self, x_1=None, x_2=None):
        self.x_1 = x_1
        self.x_2 = x_2

    def clear(self):
        self.x_1 = None
        self.x_2 = None


class DartManager:
    def __init__(self, dartboard, frame_rate=30, debug=False):
        self.dart = None
        self.dartboard = dartboard
        self.frame_rate = frame_rate
        self.debug = debug

    def update(self):
        detector1 = DartDetector(
            video_source='./videos/CAM_01_03.mp4',
            dartboard_line=347, frame_rate=self.frame_rate, debug=self.debug)
        detector2 = DartDetector(
            video_source='./videos/CAM_02_03.mp4',
            dartboard_line=330, frame_rate=self.frame_rate, debug=self.debug)

        camera_a = Camera(angle=45, fov=73.4)
        camera_b = Camera(angle=135, fov=73.4)
        triangulator = Triangulator(radius=45, width=1280)

        triangulator.add_camera(camera_a)
        triangulator.add_camera(camera_b)

        dart = Dart()
        while True:
            detector1.process_frame()
            detector2.process_frame()

            if detector1.dart is not None:
                dart.x_1 = detector1.dart.x
            if detector2.dart is not None:
                dart.x_2 = detector2.dart.x

            if dart.x_1 is not None and dart.x_2 is not None:
                print(f"X 1: {dart.x_1}, X 2: {dart.x_2}")
                point = triangulator.get_position(dart.x_1, dart.x_2)
                score = triangulator.get_score(point[0], point[1])
                self.dartboard.update_dart_position(point[0], point[1])
                # asyncio.run(self.send_data(point, score))
                print(f"Score: {score}")
                print(f"X: {point[0]}, Y: {point[1]}")
                print("--------------------")
                dart.clear()

            # Break the loop on 'q' key press
            if cv2.waitKey(30) & 0xFF == ord('q'):
                break

        # Release the resources
        detector1.release()
        detector2.release()

    async def send_data(self, point, score):
        uri = "ws://localhost:8080/ws"
        data = {
            "x": point[0],
            "y": point[1],
            "score": {
                "bed": score[0],
                "segment": score[1],
                "score": score[2]
            }
        }
        print(json.loads(json.dumps(data)))
        async with websockets.connect(uri) as websocket:
            await websocket.send(json.dumps(data))


def main():
    dartboard = Dartboard()
    dart_manager = DartManager(dartboard=dartboard, debug=True)
    Thread(target=dart_manager.update, daemon=True).start()
    dartboard.start()


if __name__ == "__main__":
    main()
