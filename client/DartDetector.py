import cv2
import numpy as np
from enum import Enum, auto
from collections import deque


class State(Enum):
    EMPTY = auto()
    CHANGING = auto()
    STABLE = auto()
    REMOVING_DARTS = auto()


class Dart:
    def __init__(self, x=0):
        self.x = x


class DartDetector:
    def __init__(self, video_source=0, dartboard_line=370, frame_rate=30, debug=False):
        """
        Initialize the DartDetector object.

        :param video_source: The video source (default is the webcam).
        """
        self.cap = cv2.VideoCapture(video_source)
        if not self.cap.isOpened():
            print("Error: Unable to open video source.")
            exit()

        # Initialize variables
        self.video_source = video_source
        self.width = 0
        self.state = State.EMPTY

        self.frame_rate = frame_rate
        self.debug = debug

        self.dartboard_line = dartboard_line
        self.detection_line = self.dartboard_line - 20
        self.trigger_line = self.detection_line - 105

        self.prev_contours = []

        self.prev_frame = None
        self.saved_frame = None
        self.should_update = False
        self.timer = 0  # Frame counter to track time in "Changing" state
        self.timer_history = deque(maxlen=5)

        self.should_send_dart_data = False

        self.dart_count = 0
        self.darts = []

    def process_frame(self):
        """
        Process the current frame and return a boolean indicating if there is a change.
        """

        frame = self.get_frame()
        _, self.width, _ = frame.shape
        mask = self.get_mask(frame)
        mask = self.dilate(mask)

        contours = self.get_contours(mask)

        self.on_change(mask, contours)

        mask = cv2.cvtColor(mask, cv2.COLOR_GRAY2BGR)

        prev = self.prev_frame
        if prev is None:
            prev = np.zeros(mask.shape[:2], dtype=np.uint8)
        prev = cv2.cvtColor(prev, cv2.COLOR_GRAY2BGR)

        inverted = cv2.bitwise_not(prev)
        difference = cv2.bitwise_and(mask, inverted)

        difference = cv2.cvtColor(difference, cv2.COLOR_BGR2GRAY)
        diff_contours = self.get_contours(difference)
        dart = self.get_current_dart(diff_contours)
        self.send_dart_data(dart)

        difference = cv2.cvtColor(difference, cv2.COLOR_GRAY2BGR)

        self.update_dart_count(contours)

        # Show the frame
        visual_frame = self.get_visual_frame(frame, contours=diff_contours)

        self.show(visual_frame, self.crop(prev),
                  self.crop(difference), self.crop(mask))

        return True

    def get_frame(self):
        ret, frame = self.cap.read()
        if not ret:
            self.cap.set(cv2.CAP_PROP_POS_FRAMES, 0)
            ret, frame = self.cap.read()
        return frame

    def get_visual_frame(self, frame, contours=None):
        visual_frame = frame.copy()
        _, width, _ = frame.shape

        if contours is not None:
            for contour in contours:
                cv2.drawContours(visual_frame, [contour], 0, (0, 0, 255), 2)

        text = f"State: {self.state.name} - Dart count: {self.dart_count}"
        cv2.putText(visual_frame, text, (10, 20),
                    cv2.FONT_HERSHEY_SIMPLEX, 0.5, (0, 0, 255), 1)

        # Draw the horizontal lines
        cv2.line(visual_frame, (0, self.trigger_line), (width, self.trigger_line),
                 (0, 0, 255), 1)  # Red line (dartboard)
        cv2.line(visual_frame, (0, self.detection_line), (width, self.detection_line),
                 (0, 0, 255), 1)  # Red line (dartboard)
        cv2.line(visual_frame, (0, self.dartboard_line), (width, self.dartboard_line),
                 (0, 0, 255), 1)  # Red line (dartboard)

        # Draw the center line
        cv2.line(visual_frame, (width // 2, self.dartboard_line), (width // 2, self.trigger_line),
                 (0, 0, 255), 1)

        for dart in self.darts:
            if frame is not None:
                cv2.circle(frame, (dart.x, self.dartboard_line),
                           5, (0, 0, 255), -1)

        return visual_frame

    def get_mask(self, frame):
        # Create a region mask between the lines
        mask = np.zeros(frame.shape[:2], dtype=np.uint8)
        mask[self.trigger_line:self.detection_line, :] = 255

        # Convert frame to grayscale
        gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)

        # Apply Otsu's thresholding with binary inversion
        _, thresh = cv2.threshold(
            gray, 0, 255, cv2.THRESH_BINARY_INV + cv2.THRESH_OTSU)

        # Apply the region mask to keep only the area between lines
        final_mask = cv2.bitwise_and(thresh, mask)

        # Apply morphological operations to clean up the mask
        kernel = np.ones((1, 1), np.uint8)
        final_mask = cv2.morphologyEx(final_mask, cv2.MORPH_OPEN, kernel)
        final_mask = cv2.morphologyEx(final_mask, cv2.MORPH_CLOSE, kernel)

        return final_mask

    def get_contours(self, frame):
        contours, _ = cv2.findContours(
            frame, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
        return contours

    def get_x_at_y(self, vx, vy, x0, y0, y):
        """Finds the x-coordinate at a specific y-value along the fitted line."""
        if vy == 0:
            return int(x0.item())  # Ensure x0 is a scalar

        x = x0.item() + ((y - y0.item()) * (vx.item() / vy.item()))  # Extract scalars
        return int(x)

    def get_current_dart(self, contours):
        if len(contours) == 1 and self.should_send_dart_data:
            [vx, vy, x0, y0] = cv2.fitLine(
                contours[0], cv2.DIST_L2, 0, 0.01, 0.01)
            x = self.get_x_at_y(vx, vy, x0, y0, self.dartboard_line)
            return Dart(x)

        return None

    def send_dart_data(self, dart):
        self.dart = dart
        if dart is not None:
            self.darts.append(dart)

        return

    def on_change(self, mask, contours):
        if len(contours) == 0:
            self.prev_frame = None
            self.saved_frame = None
            self.state = State.EMPTY
            self.dart = None
            self.darts = []
            self.timer_history.clear()
            return

        if len(contours) != len(self.prev_contours):
            self.timer = 10
            self.timer_history.append(self.timer)

            if self.state == State.STABLE or self.state == State.EMPTY:
                self.state = State.CHANGING
                if self.saved_frame is not None:
                    self.prev_frame = self.dilate(
                        self.saved_frame.copy(), (8, 8))
                self.prev_contours = contours
                return

        if self.state == State.CHANGING or self.state == State.REMOVING_DARTS:
            self.timer -= 1
            self.timer_history.append(self.timer)

            # if len(self.timer_history) == self.timer_history.maxlen and all(4 <= t <= 5 for t in self.timer_history):
            #     self.state = State.REMOVING_DARTS

            if self.timer == 0:
                self.state = State.STABLE
                self.should_send_dart_data = True
            return

        if self.state == State.STABLE:
            self.saved_frame = mask.copy()
            self.should_send_dart_data = False

        self.prev_contours = contours

    def update_dart_count(self, contours):
        if self.state == State.STABLE:
            self.dart_count = len(contours)
        if self.state == State.EMPTY:
            self.dart_count = 0

    def dilate(self, frame, kernel_size=(5, 5)):
        kernel = np.ones(kernel_size, np.uint8)
        return cv2.dilate(frame, kernel, iterations=1)

    def show(self, *args):
        if not self.debug:
            return
        combined_view = np.vstack(args)
        cv2.imshow(self.video_source, combined_view)

    def crop(self, frame):
        return frame[self.trigger_line:self.detection_line, :]

    def release(self):
        """Release the video capture and close all OpenCV windows."""
        self.cap.release()
        cv2.destroyAllWindows()


def main():
    # Create an instance of FrameDiffer
    detector = DartDetector(
        video_source='./videos/CAM_02_01.mp4', dartboard_line=330)

    # Process frames in a loop
    while True:
        detector.process_frame()
        if detector.dart is not None:
            print(f"X: {detector.dart.x}, Angle: {detector.dart.angle}")

        # Break the loop on 'q' key press
        if cv2.waitKey(detector.frame_rate) & 0xFF == ord('q'):
            break

    # Release the resources
    detector.release()


if __name__ == "__main__":
    main()
