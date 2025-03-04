from triangulate import Triangulator, Camera
from draw import Dartboard
from threading import Thread
import time

# Parameters
r = 45  # radius
angle_a = 42  # camera A angle
angle_b = 138  # camera B angle

dartboard = Dartboard()

camera_a = Camera(angle_a)
camera_b = Camera(angle_b)
triangulator = Triangulator(r)

triangulator.add_camera(camera_a)
triangulator.add_camera(camera_b)


def update_dart():
    point_a = triangulator.get_position(9.97, 16.28)
    score_a = triangulator.get_score(point_a[0], point_a[1])
    print(score_a)

    time.sleep(3)
    dartboard.update_dart_position(point_a[0], point_a[1])

    point_b = triangulator.get_position(349.76, 353.35)
    score_b = triangulator.get_score(point_b[0], point_b[1])
    print(score_b)

    time.sleep(3)
    dartboard.update_dart_position(point_b[0], point_b[1])

    point_c = triangulator.get_position(1.72, 359.69)
    score_c = triangulator.get_score(point_c[0], point_c[1])
    print(score_c)

    time.sleep(3)
    dartboard.update_dart_position(point_c[0], point_c[1])


Thread(target=update_dart, daemon=True).start()
dartboard.start()
