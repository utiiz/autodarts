import cv2
import numpy as np
import math

# Path to the video file
video_path = './videos/autodarts.mp4'

# Global variables for trackbar values
min_area = 200
max_area = 5000
aspect_ratio_threshold = 2
morphology_kernel_size = 5
dartboard_line_y = 370
detection_line_y = 350
trigger_line_y = 245

# Create a window for trackbars
cv2.namedWindow('Trackbars', cv2.WINDOW_NORMAL)


def on_min_area_change(val):
    global min_area
    min_area = val


def on_max_area_change(val):
    global max_area
    max_area = val


def on_aspect_ratio_change(val):
    global aspect_ratio_threshold
    aspect_ratio_threshold = val / 10.0


def on_morphology_kernel_change(val):
    global morphology_kernel_size
    morphology_kernel_size = val if val % 2 == 1 else val + 1


def on_dartboard_line_change(val):
    global dartboard_line_y
    dartboard_line_y = val


def on_detection_line_change(val):
    global detection_line_y
    detection_line_y = val


def on_trigger_line_change(val):
    global trigger_line_y
    trigger_line_y = val


# Create trackbars
cv2.createTrackbar('Min Area', 'Trackbars', min_area, 1000, on_min_area_change)
cv2.createTrackbar('Max Area', 'Trackbars', max_area,
                   10000, on_max_area_change)
cv2.createTrackbar('Aspect Ratio (x10)', 'Trackbars', int(
    aspect_ratio_threshold * 10), 50, on_aspect_ratio_change)
cv2.createTrackbar('Morphology Kernel', 'Trackbars',
                   morphology_kernel_size, 15, on_morphology_kernel_change)
cv2.createTrackbar('Dartboard Line Y', 'Trackbars',
                   dartboard_line_y, 700, on_dartboard_line_change)
cv2.createTrackbar('Detection Line Y', 'Trackbars',
                   detection_line_y, 700, on_detection_line_change)
cv2.createTrackbar('Trigger Line Y', 'Trackbars',
                   trigger_line_y, 700, on_trigger_line_change)

# Open the video file
cap = cv2.VideoCapture(video_path)

if not cap.isOpened():
    print("Error: Could not open video file.")
    exit()

# Get the first frame to set up the lines
ret, frame = cap.read()
if not ret:
    print("Error: Could not read from video file.")
    exit()

# Get frame dimensions
height, width = frame.shape[:2]

# Tracking for multiple dart trajectories
dart_trajectories = []
previous_contour_count = 0
detection_changed = False
change_timer = 0
CHANGE_DURATION = 30  # Number of frames to keep "Changing" state

pre_change_mask = None


def process_frame(frame, line1_y, line2_y, line3_y):
    global dart_trajectories, previous_contour_count, detection_changed, change_timer, pre_change_mask

    # Create a region mask between the lines
    region_mask = np.zeros(frame.shape[:2], dtype=np.uint8)
    region_mask[min(line2_y, line3_y):max(line2_y, line3_y), :] = 255

    # Convert frame to grayscale
    gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)

    # Apply Otsu's thresholding with binary inversion
    _, thresh = cv2.threshold(
        gray, 0, 255, cv2.THRESH_BINARY_INV + cv2.THRESH_OTSU)

    # Apply the region mask to keep only the area between lines
    final_mask = cv2.bitwise_and(thresh, region_mask)

    # Apply morphological operations to clean up the mask
    kernel = np.ones(
        (morphology_kernel_size, morphology_kernel_size), np.uint8)
    final_mask = cv2.morphologyEx(final_mask, cv2.MORPH_OPEN, kernel)
    final_mask = cv2.morphologyEx(final_mask, cv2.MORPH_CLOSE, kernel)

    # Create a copy of the original frame for visualization
    visual_frame = frame.copy()

    # Draw the horizontal lines
    cv2.line(visual_frame, (0, line1_y), (width, line1_y),
             (0, 0, 255), 1)  # Red line (dartboard)
    cv2.line(visual_frame, (0, line2_y), (width, line2_y),
             (0, 255, 0), 1)  # Green line (detection)
    cv2.line(visual_frame, (0, line3_y), (width, line3_y),
             (0, 255, 255), 1)  # Yellow line (trigger)

    # Find contours in the final mask
    contours, _ = cv2.findContours(
        final_mask, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)

    # Filter contours by size and shape to find potential darts
    dart_contours = []
    for contour in contours:
        area = cv2.contourArea(contour)
        if min_area < area < max_area:  # Adjustable area filter
            # Check if the shape is elongated (like a dart)
            x, y, w, h = cv2.boundingRect(contour)
            aspect_ratio = max(w, h) / min(w, h)
            if aspect_ratio > aspect_ratio_threshold:  # Adjustable aspect ratio
                dart_contours.append(contour)

    # Check if detection has changed
    if len(dart_contours) != previous_contour_count:
        detection_changed = True
        change_timer = CHANGE_DURATION

    # Decrement change timer
    if change_timer > 0:
        change_timer -= 1
    else:
        detection_changed = False

    previous_contour_count = len(dart_contours)

    # Clear trajectories if no darts are detected
    if not dart_contours and not detection_changed:
        dart_trajectories.clear()

    # Process each dart contour
    for i, dart_contour in enumerate(dart_contours):
        # Draw the selected contour in cyan
        cv2.drawContours(visual_frame, [dart_contour], 0, (255, 255, 0), 2)

        # Fit a line to the dart contour
        try:
            [vx, vy, x, y] = cv2.fitLine(
                dart_contour, cv2.DIST_L2, 0, 0.01, 0.01)

            # Extract scalar values from the arrays to fix deprecation warning
            vx, vy = float(vx[0]), float(vy[0])
            x, y = float(x[0]), float(y[0])

            # Calculate line endpoints for visualization
            if abs(vx) > 1e-6:  # Avoid division by zero
                lefty = int(y - x * (vy / vx))
                righty = int(y + (width - x) * (vy / vx))

                # Calculate angle
                angle_rad = math.atan2(vy, vx)
                angle_deg = math.degrees(angle_rad)
                if angle_deg < 0:
                    angle_deg += 180

                # Draw the trajectory line in blue
                current_line = ((0, lefty), (width, righty))
                cv2.line(visual_frame,
                         current_line[0], current_line[1], (255, 0, 0), 1)

                # If we haven't tracked this trajectory before, add it
                if i >= len(dart_trajectories) and not detection_changed:
                    # Calculate intersection point with dartboard line
                    intersection_x = None
                    if abs(vy) > 1e-6:  # Avoid division by zero
                        dart_x = x + (line1_y - y) * (vx / vy)
                        if 0 <= dart_x <= width:
                            intersection_x = dart_x

                    dart_trajectories.append({
                        'line': current_line,
                        'angle': angle_deg,
                        'intersection_x': intersection_x
                    })

        except cv2.error:
            # Skip if line fitting fails
            continue

    # Display all tracked angles
    for i, trajectory in enumerate(dart_trajectories):
        cv2.putText(visual_frame,
                    f"Dart {i+1} Angle: {trajectory['angle']:.1f}",
                    (10, 30 + i * 30),
                    cv2.FONT_HERSHEY_SIMPLEX,
                    0.7, (255, 0, 0), 2)
        cv2.circle(visual_frame,
                   (int(trajectory['intersection_x']), line1_y),
                   5, (0, 0, 255), -1)

    # Display detection status
    status_color = (0, 0, 255) if detection_changed else (0, 255, 0)
    status_text = "Changing" if detection_changed else "Stable"
    cv2.putText(visual_frame,
                status_text,
                (width - 150, 50),
                cv2.FONT_HERSHEY_SIMPLEX,
                1, status_color, 2)

    cv2.putText(visual_frame,
                f"Dart Count : {len(dart_trajectories)}",
                (width - 300, 100),
                cv2.FONT_HERSHEY_SIMPLEX,
                1, (255, 255, 0), 2)

    # Convert mask to 3-channel for side-by-side display
    final_mask_color = cv2.cvtColor(final_mask, cv2.COLOR_GRAY2BGR)

    return visual_frame, final_mask_color


# Main loop to process the video
while True:
    ret, frame = cap.read()
    if not ret:
        # Reset video to the beginning if it reaches the end
        cap.set(cv2.CAP_PROP_POS_FRAMES, 0)
        continue

    # Process the frame
    visual_frame, motion_mask = process_frame(
        frame, dartboard_line_y, detection_line_y, trigger_line_y
    )
    cropped_frame = motion_mask[trigger_line_y:detection_line_y, :]

    # Resize frames to fit side by side

    # Concatenate frames horizontally
    combined_view = np.vstack((visual_frame, cropped_frame))

    # Display the results
    cv2.imshow('Dart Detection', combined_view)

    # Break the loop if 'q' is pressed
    if cv2.waitKey(25) & 0xFF == ord('q'):
        break

# Release resources
cap.release()
cv2.destroyAllWindows()
