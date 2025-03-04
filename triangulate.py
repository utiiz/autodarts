import math
from scipy.optimize import fsolve


class Camera:
    def __init__(self, angle):
        self.fov = 100
        self.angle = angle


class Triangulator:
    def __init__(self, r):
        self.r = r
        self.cameras = []

    def add_camera(self, camera):
        self.cameras.append(camera)

    def get_position(self, alpha, beta):
        angle_a_rad = math.radians(self.cameras[0].angle)
        angle_b_rad = math.radians(self.cameras[1].angle)

        def line_a(x):
            return math.tan(math.radians(alpha + self.cameras[0].angle)) * (x - self.r * math.cos(angle_a_rad)) + self.r * math.sin(angle_a_rad)

        def line_b(x):
            return math.tan(math.radians(beta + self.cameras[1].angle)) * (x - self.r * math.cos(angle_b_rad)) + self.r * math.sin(angle_b_rad)

        def intersection_equation(x):
            return line_a(x) - line_b(x)

        x_p = fsolve(intersection_equation, 0)[0]
        y_p = line_a(x_p)
        return x_p, y_p

    def determine_segment(self, x, y, num_segments=20, offset=9):
        angle_rad = math.atan2(y, x)
        angle_deg = math.degrees(angle_rad)

        if angle_deg < 0:
            angle_deg += 360

        angle_with_offset = (angle_deg - offset) % 360
        segment_size = 360 / num_segments
        segment_index = int(angle_with_offset // segment_size) + 1
        return segment_index

    def transform_segment(self, index):
        mapping = [13, 4, 18, 1, 20, 5, 12, 9, 14,
                   11, 8, 16, 7, 19, 3, 17, 2, 15, 10, 6]

        if index < 1 or index > len(mapping):
            raise ValueError("Index out of range")
        return mapping[index - 1]

    def get_score(self, x, y):
        zones = [
            (17, "D"),
            (16.2, "S"),
            (10.7, "T"),
            (9.9, "S"),
            (1.6, "S-BULL"),
            (0.635, "D-BULL"),
        ]

        distance = math.sqrt(x**2 + y**2)
        applicable_zone = None
        for radius, score in zones:
            if distance <= radius:
                applicable_zone = (radius, score)

        if applicable_zone is None:
            return None

        radius, score = applicable_zone
        segment = self.determine_segment(x, y, offset=9)
        transformed_segment = self.transform_segment(segment)

        if score == "D":
            return [score, transformed_segment, transformed_segment * 2]
        elif score == "T":
            return [score, transformed_segment, transformed_segment * 3]
        elif score == "S":
            return [score, transformed_segment, transformed_segment]
        elif score == "S-BULL":
            return [score, 25, 25]
        elif score == "D-BULL":
            return [score, 50, 50]
