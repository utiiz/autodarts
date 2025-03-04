import tkinter as tk
import math


class Dartboard:
    def __init__(self):
        self.root = tk.Tk()
        self.root.title("Dartboard")

        # Canvas dimensions
        self.canvas_width = 600
        self.canvas_height = 600
        self.num_divisions = 20
        self.radius = 200
        self.canvas = tk.Canvas(self.root, width=self.canvas_width,
                                height=self.canvas_height, bg="black")

        self.center_x = self.canvas_width // 2
        self.center_y = self.canvas_height // 2
        self.canvas.configure(scrollregion=(-self.center_x, -self.center_y,
                                            self.center_x, self.center_y))
        self.canvas.pack()

        self.mapping = [
            13, 4, 18, 1, 20, 5,
            12, 9, 14, 11, 8,
            16, 7, 19, 3, 17,
            2, 15, 10, 6
        ]

        self.rings = [
            17, 16.2, 10.7, 9.9, 1.6, 0.635
        ]

        self.draw_radial_grid()

    def draw_radial_grid(self):
        # Draw concentric circles
        for i in range(1, len(self.rings) - 1):
            radius = self.rings[i - 1] * self.radius / 17
            self.canvas.create_oval(
                - radius, - radius,
                radius, radius,
                outline="cyan"
            )

        for i in range(self.num_divisions):
            # Calculate the angles for the segment borders
            start_angle = math.radians(
                9) + 2 * math.pi * (i / self.num_divisions)
            end_angle = math.radians(9) + 2 * math.pi * \
                ((i + 1) / self.num_divisions)

            # Find the midpoint angle
            midpoint_angle = (start_angle + end_angle) / 2

            # Compute the radial line endpoints (for visual clarity)
            x_end = self.radius * math.cos(start_angle)
            # Subtract y (canvas grows downward)
            y_end = self.radius * math.sin(start_angle)

            self.canvas.create_line(0, 0,
                                    x_end, y_end, fill="cyan")

            # Calculate label position at midpoint, slightly outside the circle
            label_radius = self.radius + 20
            label_x = label_radius * math.cos(midpoint_angle)
            label_y = - label_radius * math.sin(midpoint_angle)
            self.canvas.create_text(label_x, label_y, text=str(
                self.mapping[i]), fill="yellow", font=("Arial", 10))

        for i in range(len(self.rings) - 1, len(self.rings) + 1):
            radius = self.rings[i - 1] * self.radius / 17
            self.canvas.create_oval(
                - radius, - radius,
                radius, radius,
                outline="cyan",
                fill="black"
            )

    def update_dart_position(self, x, y):
        x = (x * self.radius / 17)
        y = (-y * self.radius / 17)

        # Draw a line from
        x_cam = 45 * self.radius / 17 * math.cos(math.radians(42))
        y_cam = - 45 * self.radius / 17 * math.sin(math.radians(42))

        self.canvas.create_line(x_cam, y_cam, x, y, fill="green")

        x_cam = 45 * self.radius / 17 * math.cos(math.radians(138))
        y_cam = - 45 * self.radius / 17 * math.sin(math.radians(138))

        self.canvas.create_line(x_cam, y_cam, x, y, fill="green")

        # Draw a small red dot
        self.canvas.create_oval(
            x - 5, y - 5,
            x + 5, y + 5,
            fill="red"
        )

    def start(self):
        self.root.mainloop()
