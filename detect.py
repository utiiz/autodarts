import cv2
import numpy as np
import math

# Path to the video file
video_path = './videos/autodarts.mp4'

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

# Initial positions for the horizontal lines
dartboard_line_y = 380  # Middle of the frame
trigger_line_y = int(height * 0.3)    # Above the dartboard line

# Create background subtractor
bg_subtractor = cv2.createBackgroundSubtractorMOG2(
    history=20, varThreshold=25, detectShadows=False)

# Previous frame for motion detection
prev_frame = None

# Variables to store the last detected dart trajectory
last_dart_line = None
last_dart_angle = None

# Function to apply thresholding and detect dart


def process_frame(frame, prev_frame, line1_y, line2_y, last_line, last_angle):
    # Create a region mask between the lines
    region_mask = np.zeros(frame.shape[:2], dtype=np.uint8)
    region_mask[min(line1_y, line2_y):max(line1_y, line2_y), :] = 255

    # Convert frame to grayscale
    gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)

    # Apply background subtraction to detect moving objects
    fg_mask = bg_subtractor.apply(frame)

    # Additional motion detection using frame differencing if we have a previous frame
    if prev_frame is not None:
        prev_gray = cv2.cvtColor(prev_frame, cv2.COLOR_BGR2GRAY)
        frame_diff = cv2.absdiff(gray, prev_gray)
        _, motion_mask = cv2.threshold(frame_diff, 30, 255, cv2.THRESH_BINARY)

        # Combine background subtraction with frame differencing
        combined_mask = cv2.bitwise_and(fg_mask, motion_mask)
    else:
        combined_mask = fg_mask

    # Apply the region mask to keep only motion between the lines
    final_mask = cv2.bitwise_and(combined_mask, region_mask)

    # Apply morphological operations to clean up the mask
    kernel = np.ones((5, 5), np.uint8)
    final_mask = cv2.morphologyEx(final_mask, cv2.MORPH_OPEN, kernel)
    final_mask = cv2.morphologyEx(final_mask, cv2.MORPH_CLOSE, kernel)

    # Create a copy of the original frame for visualization
    visual_frame = frame.copy()

    # Draw the horizontal lines
    cv2.line(visual_frame, (0, line1_y), (width, line1_y),
             (0, 0, 255), 1)  # Red line (dartboard)
    cv2.line(visual_frame, (0, line2_y), (width, line2_y),
             (0, 255, 0), 1)  # Green line (trigger)

    # Find contours in the final mask
    contours, _ = cv2.findContours(
        final_mask, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)

    # Filter contours by size and shape to find potential darts
    dart_contours = []
    for contour in contours:
        area = cv2.contourArea(contour)
        if area > 200 and area < 5000:  # Filter by area - adjust these thresholds
            # Check if the shape is elongated (like a dart)
            x, y, w, h = cv2.boundingRect(contour)
            aspect_ratio = max(w, h) / min(w, h)
            if aspect_ratio > 2:  # Darts are typically elongated
                dart_contours.append(contour)

    # Variables to store current dart data
    current_line = last_line
    current_angle = last_angle

    # Process dart contours
    if dart_contours:
        # Use the largest valid dart contour
        dart_contour = max(dart_contours, key=cv2.contourArea)

        # Draw the contour for visualization
        cv2.drawContours(visual_frame, [dart_contour], 0, (0, 255, 255), 2)

        # Fit a line to the dart contour
        [vx, vy, x, y] = cv2.fitLine(dart_contour, cv2.DIST_L2, 0, 0.01, 0.01)

        # Extract scalar values from the arrays to fix deprecation warning
        vx, vy = float(vx[0]), float(vy[0])
        x, y = float(x[0]), float(y[0])

        # Calculate line endpoints for visualization
        if abs(vx) > 1e-6:  # Avoid division by zero
            lefty = int(y - x * (vy / vx))
            righty = int(y + (width - x) * (vy / vx))

            # Update the current line and angle
            current_line = ((0, lefty), (width, righty))

            # Calculate angle
            angle_rad = math.atan2(vy, vx)
            angle_deg = math.degrees(angle_rad)
            if angle_deg < 0:
                angle_deg += 180

            current_angle = angle_deg

    # Draw the current or last dart trajectory line if available
    intersection_x = None
    if current_line:
        # Blue line
        cv2.line(visual_frame, current_line[0],
                 current_line[1], (255, 0, 0), 1)

        # Calculate intersection point with the dartboard line (red line)
        p1, p2 = current_line
        x1, y1 = p1
        x2, y2 = p2

        # Check if the line is not horizontal (to avoid division by zero)
        if y1 != y2:
            # Calculate the x-coordinate of intersection
            # Equation: x = x1 + (line1_y - y1) * (x2 - x1) / (y2 - y1)
            intersection_x = x1 + (line1_y - y1) * (x2 - x1) / (y2 - y1)
            intersection_x = int(intersection_x)

            # Check if the intersection point is within the frame width
            if 0 <= intersection_x < width:
                # Draw a red circle at the intersection point
                cv2.circle(visual_frame, (intersection_x, line1_y),
                           3, (0, 0, 255), -1)  # Red filled circle

                # Display the X coordinate of the intersection point
                cv2.putText(visual_frame, f"Hit position: {intersection_x}",
                            (10, 60), cv2.FONT_HERSHEY_SIMPLEX,
                            0.7, (0, 0, 255), 2)

    # Display the current or last angle if available
    if current_angle is not None:
        cv2.putText(visual_frame, f"Angle: {current_angle:.1f} degrees",
                    (10, 30), cv2.FONT_HERSHEY_SIMPLEX,
                    0.7, (255, 0, 0), 2)

    return visual_frame, final_mask, current_line, current_angle


# Main loop to process the video
while True:
    ret, frame = cap.read()
    if not ret:
        break

    # Process the frame
    visual_frame, motion_mask, last_dart_line, last_dart_angle = process_frame(
        frame, prev_frame, dartboard_line_y, trigger_line_y, last_dart_line, last_dart_angle
    )

    # Update previous frame
    prev_frame = frame.copy()

    # Display the results
    cv2.imshow('Dart Tracking', visual_frame)
    cv2.imshow('Motion Mask', motion_mask)

    # Break the loop if 'q' is pressed
    if cv2.waitKey(25) & 0xFF == ord('q'):
        break

# Release resources
cap.release()
cv2.destroyAllWindows()
