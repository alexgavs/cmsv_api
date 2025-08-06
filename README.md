# CMSV API Project

This project is a Go application for interacting with the CMSV8 API, providing a user interface for device management, vehicle information, and real-time monitoring.

## Features
- User login/logout
- Retrieve user vehicle information
- Get device online status
- Get real-time device status
- RTSP stream configuration for devices

## API Endpoints
See `api_description.md` for detailed API documentation, including request/response examples and error codes.

## Requirements
- Go 1.18+
- Fyne UI library

## Getting Started
1. Clone the repository:
   ```sh
   git clone <repo-url>
   cd cmsv_api
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```
3. Run the application:
   ```sh
   go run main.go
   ```

## Usage
- Launch the application and log in with your CMSV8 account.
- Select a device to view or configure.
- Use the UI to access vehicle and device information, and configure RTSP streams.

## Documentation
- API details: [api_description.md](api_description.md)
- Main application logic: [main.go](main.go)

## License
This project is licensed under the MIT License.

