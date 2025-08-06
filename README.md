# CMSV API Client

A GUI application for interacting with CMSV (Commercial Vehicle Monitoring System) API. This application provides an easy-to-use interface for managing vehicle tracking devices, generating streaming links, and monitoring alarms.

## Features

- **User Authentication**: Login to CMSV system with account credentials
- **Device Management**: View and manage connected tracking devices
- **Vehicle Information**: Display detailed vehicle and company hierarchy information
- **Real-time Alarms**: Monitor device alarms with auto-refresh capability
- **Streaming Links Generation**: Generate RTSP, RTMP, and HLS streaming URLs
- **Configurable Interface**: Customize UI elements visibility through configuration
- **Multiple Coordinate Systems**: Support for WGS84, Google (GCJ-02), and Baidu (BD-09) coordinates

## Configuration

The application uses a `config.ini` file to manage settings. If the file doesn't exist, it will be created automatically with default values.

### Server Configuration
```ini
# Server URL
server_url = "https://ahd.samsonix.com"

# API Port
api_port = 443

# RTMP Port
rtmp_port = 1935

# RTSP Port
rtsp_port = 6604

# HLS Port
hls_port = 16604
```

### UI Elements Visibility (1 = show, 0 = hide)
```ini
show_login_button = 1
show_save_button = 1
show_vehicle_info_button = 1
show_device_alarms_button = 1
show_auto_refresh_button = 1
show_rtsp_button = 1
show_rtmp_button = 1
show_hls_button = 1
show_company_hierarchy = 0
```

## Installation

### Prerequisites
- Go 1.19 or later
- Windows, macOS, or Linux

### Building from Source
1. Clone the repository or download the source code
2. Navigate to the project directory
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Build the application:
   ```bash
   go build -o cmsv_api main.go
   ```

## Usage

### Running the Application
Execute the built binary:
```bash
./cmsv_api
```

### Basic Workflow
1. **Login**: Enter your CMSV account credentials and click "Login and Fetch Devices"
2. **Select Device**: Choose a device from the dropdown menu
3. **Monitor Alarms**: Click "GET DEVICE ALARMS" to view current alarms
4. **Generate Streaming Links**: Use RTSP, RTMP, or HLS buttons to generate streaming URLs
5. **Save Data**: Use "Save to File" to export device information

### Features Overview

#### Device Management
- Login with CMSV credentials
- View all authorized devices
- Real-time device status monitoring

#### Vehicle Information
- Display vehicle details including owner, engine number, frame number
- Show company hierarchy (can be disabled in config)
- Device installation information

#### Alarm Monitoring
- View device alarms with detailed information
- Support for different coordinate systems
- Auto-refresh functionality for real-time monitoring
- Alarm logging to file

#### Streaming Links
- **RTSP**: Real-Time Streaming Protocol links for video players
- **RTMP**: Real-Time Messaging Protocol for streaming servers
- **HLS**: HTTP Live Streaming for web browsers
- Configurable stream quality (main/sub stream)
- Multiple channel support

## API Documentation

See `api_description.md` for detailed API endpoint documentation including:
- User authentication endpoints
- Device management APIs
- Vehicle information retrieval
- Alarm monitoring APIs
- Error codes and troubleshooting

## File Structure

```
cmsv_api/
├── main.go              # Main application file
├── config.ini           # Configuration file
├── api_description.md   # API documentation
├── README.md           # This file
├── go.mod              # Go module file
├── go.sum              # Go dependencies
└── alarms.log          # Alarm log file (created automatically)
```

## Configuration Options

### UI Customization
You can hide/show interface elements by modifying the config.ini file:

- `show_login_button = 0` - Hide the login button
- `show_company_hierarchy = 0` - Hide company hierarchy in vehicle information
- `show_rtsp_button = 0` - Hide RTSP link generation button

### Server Configuration
- Change `server_url` to point to your CMSV server
- Modify port settings for different streaming protocols

## Troubleshooting

### Common Issues

1. **SSL Certificate Errors**: The application automatically handles SSL certificate issues by falling back to insecure connections when necessary.

2. **Login Failures**: Check your credentials and ensure the server URL is correct in config.ini.

3. **No Devices Found**: Ensure your account has proper permissions to access devices.

4. **Streaming Links Not Working**: Verify that the streaming ports are accessible and not blocked by firewalls.

### Log Files
- Application logs are displayed in the output area
- Alarm data is automatically saved to `alarms.log`

## Dependencies

- [Fyne](https://fyne.io/) - Cross-platform GUI toolkit
- Go standard library for HTTP requests and configuration parsing

## License

This project is for internal use with CMSV systems. Please ensure compliance with your organization's software usage policies.

## Support

For technical support or feature requests, please contact your system administrator or refer to the CMSV system documentation.
