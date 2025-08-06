package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// AppConfig holds all configuration values
type AppConfig struct {
	ServerURL string
	APIPort   int
	RTMPPort  int
	RTSPPort  int
	HLSPort   int
}

// Global config variable
var config AppConfig

// loadConfig reads the configuration from config.ini file
func loadConfig() error {
	// Set default values
	config = AppConfig{
		ServerURL: "https://cloud.samsonix.com",
		APIPort:   443,
		RTMPPort:  1935,
		RTSPPort:  6604,
		HLSPort:   16604,
	}

	file, err := os.Open("config.ini")
	if err != nil {
		// If config file doesn't exist, use defaults and create one
		fmt.Println("Config file not found, using defaults and creating config.ini")
		return createDefaultConfig()
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key = value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes from value if present
		if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}

		switch key {
		case "server_url":
			config.ServerURL = value
		case "api_port":
			if port, err := strconv.Atoi(value); err == nil {
				config.APIPort = port
			}
		case "rtmp_port":
			if port, err := strconv.Atoi(value); err == nil {
				config.RTMPPort = port
			}
		case "rtsp_port":
			if port, err := strconv.Atoi(value); err == nil {
				config.RTSPPort = port
			}
		case "hls_port":
			if port, err := strconv.Atoi(value); err == nil {
				config.HLSPort = port
			}
		}
	}

	return scanner.Err()
}

// createDefaultConfig creates a default config.ini file
func createDefaultConfig() error {
	content := `# Application Configuration File

# Server URL
server_url = "https://cloud.samsonix.com"

# API Port
api_port = 443

# RTMP Port
rtmp_port = 1935

# RTSP Port
rtsp_port = 6604

# HLS Port
hls_port = 16604

# Other configuration options can be added below
# Example:
# map_port = 8080
`
	return os.WriteFile("config.ini", []byte(content), 0644)
}

type Config struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Result   int    `json:"result"`
	JSession string `json:"jsession"`
}

type Device struct {
	VID string `json:"vid"`
	DID string `json:"did"`
}

type StatusResponse struct {
	Result  int      `json:"result"`
	Onlines []Device `json:"onlines"`
}

type VehicleResponse struct {
	Result   int `json:"result"`
	Companys []struct {
		ID   int    `json:"id"`
		Name string `json:"nm"`
		PID  int    `json:"pId"`
	} `json:"companys"`
	Vehicles []struct {
		ID         int    `json:"id"`
		Name       string `json:"nm"`
		PID        int    `json:"pid"`
		PName      string `json:"pnm"`
		DeviceList []struct {
			ID          string `json:"id"`
			Channels    int    `json:"cc"`
			ChanName    string `json:"cn"`
			SIM         string `json:"sim"`
			InstallTime string `json:"ist"`
		} `json:"dl"`
		VehicleType  string `json:"vehiType"`
		VehicleColor string `json:"vehiColor"`
		VehicleBand  string `json:"vehiBand"`
		OwnerName    string `json:"ownerName"`
		EngineNum    string `json:"engineNum"`
		FrameNum     string `json:"frameNum"`
	} `json:"vehicles"`
}

// EquipmentStatus represents the bit-by-bit status flags for equipment
type EquipmentStatus struct {
	// S1 flags (32 bits)
	GPSValid           bool // s1:0 - GPS positioning status (0=invalid, 1=valid)
	ACCStatus          bool // s1:1 - ACC status (0=off, 1=on)
	LeftTurn           bool // s1:2 - Left turn status
	RightTurn          bool // s1:3 - Right turn status
	FatigueWarning     bool // s1:4 - Fatigue driving warning
	ForwardRotation    bool // s1:5 - Positive rotation state
	ReverseState       bool // s1:6 - Reverse state
	GPSAntennaPresent  bool // s1:7 - GPS antenna present
	HardDriveStatus    int  // s1:8-9 - Hard drive status (0=not present, 1=present, 2=power down)
	ThreeGModuleStatus int  // s1:10-12 - 3G module status (0-5)
	QuiescentState     bool // s1:13 - Quiescent state
	OverspeedState     bool // s1:14 - Overspeed state
	GPSSupplement      bool // s1:15 - GPS supplement
	BatteryStatus      bool // s1:16 - Battery status
	NightState         bool // s1:17 - Night state
	OvercrowdingStatus bool // s1:18 - Overcrowding status
	ParkingACCStatus   bool // s1:19 - Parking ACC status
	IO1Status          bool // s1:20 - IO1 status
	IO2Status          bool // s1:21 - IO2 status
	IO3Status          bool // s1:22 - IO3 status
	IO4Status          bool // s1:23 - IO4 status
	IO5Status          bool // s1:24 - IO5 status
	IO6Status          bool // s1:25 - IO6 status
	IO7Status          bool // s1:26 - IO7 status
	IO8Status          bool // s1:27 - IO8 status
	Drive2Status       bool // s1:28 - Drive 2 status
	HardDisk2Status    int  // s1:29-30 - Hard disk 2 status
	HardDiskInvalid    bool // s1:31 - Hard disk status

	// S2 flags (32 bits)
	OutOfAreaAlarm            bool // s2:0 - Out of area alarm
	LineAlarm                 bool // s2:1 - Line alarm
	HighSpeedInAreaAlarm      bool // s2:2 - High speed in area
	LowSpeedInAreaAlarm       bool // s2:3 - Low speed in area
	HighSpeedOutsideAreaAlarm bool // s2:4 - High speed outside area
	LowSpeedOutsideAreaAlarm  bool // s2:5 - Low speed outside area
	ParkingInAreaAlarm        bool // s2:6 - Parking in area alarm
	OutOfAreaParkingAlarm     bool // s2:7 - Out of area parking alarm
	DailyFlowWarning          bool // s2:8 - Daily flow warning
	DailyFlowExceeded         bool // s2:9 - Daily flow exceeded
	MonthlyTrafficWarning     bool // s2:10 - Monthly traffic warning
	MonthlyFlowExceeded       bool // s2:11 - Monthly flow exceeded
	BackupBatteryPowered      bool // s2:12 - Host powered by backup battery
	DoorOpen                  bool // s2:13 - Door open
	VehicleFortification      bool // s2:14 - Vehicle fortification
	BatteryVoltageLow         bool // s2:15 - Battery voltage too low
	EngineStatus              bool // s2:17 - Engine status
	LastValidGPSInfo          bool // s2:18 - Last valid GPS information
	OnBoardStatus             bool // s2:19 - On board status (0=no load, 1=heavy load)
	OperationStatus           bool // s2:20 - Operation status (1=shutdown)
	LatLngNotEncrypted        bool // s2:21 - Latitude and longitude not encrypted
	NormalOilCircuit          bool // s2:22 - Normal oil circuit (1=disconnected)
	CircuitOK                 bool // s2:23 - Circuit OK (1=disconnected)
	DoorUnlock                bool // s2:24 - Door unlock (1=locked)
	AreaOverspeedPlatform     bool // s2:25 - Area overspeed alarm (platform)
	AreaOverspeedPlatform2    bool // s2:26 - Area overspeed alarm (platform)
	IntoAreaAlarm             bool // s2:27 - Into area alarm (platform)
	LineOffset                bool // s2:28 - Line offset (platform)
	TimePeriodOverspeed       bool // s2:29 - Time period overspeed (platform)
	TimePeriodLowSpeed        bool // s2:30 - Time period low speed (platform)
	FatigueDriving            bool // s2:31 - Fatigue driving (platform)

	// S3 flags (32 bits)
	VideoLostChannels    uint8 // s3:0-7 - Channel video lost
	VideoChannels        uint8 // s3:8-15 - Channel video
	IOInputs916          uint8 // s3:16-23 - IO inputs 9-16
	IOOutput14           uint8 // s3:24-27 - IO output 1-4
	PositioningType      uint8 // s3:28-29 - Positioning (0=GPS, 1=base station, 2=WiFi)
	AbnormalDrivingState bool  // s3:30 - Abnormal driving state (passenger cars forbidden)
	MountainForbidden    bool  // s3:31 - Mountain forbidden line

	// S4 flags (32 bits)
	PositioningCoordType      uint8 // s4:0-2 - Positioning type (0=WGS84, 1=GCJ-02, 2=BD09)
	EmergencyAlarm            bool  // s4:3 - Emergency alarm
	AreaOverspeedAlarm        bool  // s4:4 - Area overspeed alarm
	FatigueDrivingReport      bool  // s4:5 - Fatigue driving report
	DangerousDrivingAlarm     bool  // s4:6 - Dangerous driving behavior alarm
	GNSSModuleFault           bool  // s4:7 - GNSS module fault alarm
	GNSSAntennaDisconnected   bool  // s4:8 - GNSS antenna not connected/cut off
	GNSSAntennaShortCircuit   bool  // s4:9 - GNSS antenna short circuit
	TerminalLCDFault          bool  // s4:10 - Terminal LCD/display failure
	TTSModuleFault            bool  // s4:11 - TTS module fault
	CameraFailure             bool  // s4:12 - Camera failure
	CumulativeDrivingOvertime bool  // s4:13 - Cumulative driving overtime
	OvertimeParking           bool  // s4:14 - Overtime parking
	IntoArea                  bool  // s4:15 - Into area
	RouteAlarm                bool  // s4:16 - Route alarm
	TravelTimeAbnormal        bool  // s4:17 - Insufficient/excessive travel time
	RouteDeviationAlarm       bool  // s4:18 - Route deviation alarm
	VSSFailure                bool  // s4:19 - Vehicle VSS failure
	FuelQuantityAbnormal      bool  // s4:20 - Abnormal fuel quantity
	VehicleTheftAlarm         bool  // s4:21 - Vehicle theft alarm
	IllegalIgnitionAlarm      bool  // s4:22 - Illegal ignition alarm
	IllegalDisplacementAlarm  bool  // s4:23 - Illegal displacement alarm
	CollisionRolloverAlarm    bool  // s4:24 - Collision rollover alarm
	OvertimeStop              bool  // s4:25 - Overtime stop (platform)
	KeyPointNotReachedAlarm   bool  // s4:26 - Key point not reached (platform)
	LineOverspeedAlarm        bool  // s4:27 - Line overspeed alarm (platform)
	LineLowSpeedAlarm         bool  // s4:28 - Line low speed alarm (platform)
	RoadOverspeedAlarm        bool  // s4:29 - Road overspeed alarm (platform)
	OutOfAreaAlarmPlatform    bool  // s4:30 - Out of area alarm (platform)
	KeyPointNotLeaveAlarm     bool  // s4:31 - Key points not leave alarm (platform)
}

// ParseEquipmentStatus parses the s1, s2, s3, s4 integers into a structured EquipmentStatus
func ParseEquipmentStatus(s1, s2, s3, s4 int) EquipmentStatus {
	status := EquipmentStatus{}

	// Parse S1 flags
	status.GPSValid = (s1 & 0x01) != 0
	status.ACCStatus = (s1 & 0x02) != 0
	status.LeftTurn = (s1 & 0x04) != 0
	status.RightTurn = (s1 & 0x08) != 0
	status.FatigueWarning = (s1 & 0x10) != 0
	status.ForwardRotation = (s1 & 0x20) != 0
	status.ReverseState = (s1 & 0x40) != 0
	status.GPSAntennaPresent = (s1 & 0x80) != 0
	status.HardDriveStatus = (s1 >> 8) & 0x03
	status.ThreeGModuleStatus = (s1 >> 10) & 0x07
	status.QuiescentState = (s1 & 0x2000) != 0
	status.OverspeedState = (s1 & 0x4000) != 0
	status.GPSSupplement = (s1 & 0x8000) != 0
	status.BatteryStatus = (s1 & 0x10000) != 0
	status.NightState = (s1 & 0x20000) != 0
	status.OvercrowdingStatus = (s1 & 0x40000) != 0
	status.ParkingACCStatus = (s1 & 0x80000) != 0
	status.IO1Status = (s1 & 0x100000) != 0
	status.IO2Status = (s1 & 0x200000) != 0
	status.IO3Status = (s1 & 0x400000) != 0
	status.IO4Status = (s1 & 0x800000) != 0
	status.IO5Status = (s1 & 0x1000000) != 0
	status.IO6Status = (s1 & 0x2000000) != 0
	status.IO7Status = (s1 & 0x4000000) != 0
	status.IO8Status = (s1 & 0x8000000) != 0
	status.Drive2Status = (s1 & 0x10000000) != 0
	status.HardDisk2Status = (s1 >> 29) & 0x03
	status.HardDiskInvalid = (s1 & 0x80000000) != 0

	// Parse S2 flags
	status.OutOfAreaAlarm = (s2 & 0x01) != 0
	status.LineAlarm = (s2 & 0x02) != 0
	status.HighSpeedInAreaAlarm = (s2 & 0x04) != 0
	status.LowSpeedInAreaAlarm = (s2 & 0x08) != 0
	status.HighSpeedOutsideAreaAlarm = (s2 & 0x10) != 0
	status.LowSpeedOutsideAreaAlarm = (s2 & 0x20) != 0
	status.ParkingInAreaAlarm = (s2 & 0x40) != 0
	status.OutOfAreaParkingAlarm = (s2 & 0x80) != 0
	status.DailyFlowWarning = (s2 & 0x100) != 0
	status.DailyFlowExceeded = (s2 & 0x200) != 0
	status.MonthlyTrafficWarning = (s2 & 0x400) != 0
	status.MonthlyFlowExceeded = (s2 & 0x800) != 0
	status.BackupBatteryPowered = (s2 & 0x1000) != 0
	status.DoorOpen = (s2 & 0x2000) != 0
	status.VehicleFortification = (s2 & 0x4000) != 0
	status.BatteryVoltageLow = (s2 & 0x8000) != 0
	status.EngineStatus = (s2 & 0x20000) != 0
	status.LastValidGPSInfo = (s2 & 0x40000) != 0
	status.OnBoardStatus = (s2 & 0x80000) != 0
	status.OperationStatus = (s2 & 0x100000) != 0
	status.LatLngNotEncrypted = (s2 & 0x200000) != 0
	status.NormalOilCircuit = (s2 & 0x400000) != 0
	status.CircuitOK = (s2 & 0x800000) != 0
	status.DoorUnlock = (s2 & 0x1000000) != 0
	status.AreaOverspeedPlatform = (s2 & 0x2000000) != 0
	status.AreaOverspeedPlatform2 = (s2 & 0x4000000) != 0
	status.IntoAreaAlarm = (s2 & 0x8000000) != 0
	status.LineOffset = (s2 & 0x10000000) != 0
	status.TimePeriodOverspeed = (s2 & 0x20000000) != 0
	status.TimePeriodLowSpeed = (s2 & 0x40000000) != 0
	status.FatigueDriving = (s2 & 0x80000000) != 0

	// Parse S3 flags
	status.VideoLostChannels = uint8(s3 & 0xFF)
	status.VideoChannels = uint8((s3 >> 8) & 0xFF)
	status.IOInputs916 = uint8((s3 >> 16) & 0xFF)
	status.IOOutput14 = uint8((s3 >> 24) & 0x0F)
	status.PositioningType = uint8((s3 >> 28) & 0x03)
	status.AbnormalDrivingState = (s3 & 0x40000000) != 0
	status.MountainForbidden = (s3 & 0x80000000) != 0

	// Parse S4 flags
	status.PositioningCoordType = uint8(s4 & 0x07)
	status.EmergencyAlarm = (s4 & 0x08) != 0
	status.AreaOverspeedAlarm = (s4 & 0x10) != 0
	status.FatigueDrivingReport = (s4 & 0x20) != 0
	status.DangerousDrivingAlarm = (s4 & 0x40) != 0
	status.GNSSModuleFault = (s4 & 0x80) != 0
	status.GNSSAntennaDisconnected = (s4 & 0x100) != 0
	status.GNSSAntennaShortCircuit = (s4 & 0x200) != 0
	status.TerminalLCDFault = (s4 & 0x400) != 0
	status.TTSModuleFault = (s4 & 0x800) != 0
	status.CameraFailure = (s4 & 0x1000) != 0
	status.CumulativeDrivingOvertime = (s4 & 0x2000) != 0
	status.OvertimeParking = (s4 & 0x4000) != 0
	status.IntoArea = (s4 & 0x8000) != 0
	status.RouteAlarm = (s4 & 0x10000) != 0
	status.TravelTimeAbnormal = (s4 & 0x20000) != 0
	status.RouteDeviationAlarm = (s4 & 0x40000) != 0
	status.VSSFailure = (s4 & 0x80000) != 0
	status.FuelQuantityAbnormal = (s4 & 0x100000) != 0
	status.VehicleTheftAlarm = (s4 & 0x200000) != 0
	status.IllegalIgnitionAlarm = (s4 & 0x400000) != 0
	status.IllegalDisplacementAlarm = (s4 & 0x800000) != 0
	status.CollisionRolloverAlarm = (s4 & 0x1000000) != 0
	status.OvertimeStop = (s4 & 0x2000000) != 0
	status.KeyPointNotReachedAlarm = (s4 & 0x4000000) != 0
	status.LineOverspeedAlarm = (s4 & 0x8000000) != 0
	status.LineLowSpeedAlarm = (s4 & 0x10000000) != 0
	status.RoadOverspeedAlarm = (s4 & 0x20000000) != 0
	status.OutOfAreaAlarmPlatform = (s4 & 0x40000000) != 0
	status.KeyPointNotLeaveAlarm = (s4 & 0x80000000) != 0

	return status
}

type AlarmResponse struct {
	Result    int `json:"result"`
	AlarmList []struct {
		DevIDNO string `json:"DevIDNO"`
		Desc    string `json:"desc"`
		GUID    string `json:"guid"`
		HD      int    `json:"hd"`
		Img     string `json:"img"`
		Info    int    `json:"info"`
		P1      int    `json:"p1"`
		P2      int    `json:"p2"`
		P3      int    `json:"p3"`
		P4      int    `json:"p4"`
		SrcTm   string `json:"srcTm"`
		StType  int    `json:"stType"`
		Time    string `json:"time"`
		Type    int    `json:"type"`
		Gps     struct {
			DCT  int    `json:"dct"`
			GD   int    `json:"gd"`
			GT   string `json:"gt"`
			HX   int    `json:"hx"`
			Lat  int    `json:"lat"`
			LC   int    `json:"lc"`
			LID  int    `json:"lid"`
			Lng  int    `json:"lng"`
			MLat string `json:"mlat"`
			MLng string `json:"mlng"`
			SP   int    `json:"sp"`
		} `json:"Gps"`
	} `json:"alarmlist"`
	Pagination struct {
		TotalPages   int `json:"totalPages"`
		CurrentPage  int `json:"currentPage"`
		PageRecords  int `json:"pageRecords"`
		TotalRecords int `json:"totalRecords"`
	} `json:"pagination"`
}

type AlarmResponseAlarm = struct {
	DevIDNO string `json:"DevIDNO"`
	Desc    string `json:"desc"`
	GUID    string `json:"guid"`
	HD      int    `json:"hd"`
	Img     string `json:"img"`
	Info    int    `json:"info"`
	P1      int    `json:"p1"`
	P2      int    `json:"p2"`
	P3      int    `json:"p3"`
	P4      int    `json:"p4"`
	SrcTm   string `json:"srcTm"`
	StType  int    `json:"stType"`
	Time    string `json:"time"`
	Type    int    `json:"type"`
	Gps     struct {
		DCT  int    `json:"dct"`
		GD   int    `json:"gd"`
		GT   string `json:"gt"`
		HX   int    `json:"hx"`
		Lat  int    `json:"lat"`
		LC   int    `json:"lc"`
		LID  int    `json:"lid"`
		Lng  int    `json:"lng"`
		MLat string `json:"mlat"`
		MLng string `json:"mlng"`
		SP   int    `json:"sp"`
	} `json:"Gps"`
}

// RTSPLinkOptions contains the parameters needed to build an RTSP URL
type RTSPLinkOptions struct {
	ServerHost string // RTSP server hostname
	ServerPort int    // RTSP server port (default 6604)
	JSession   string // Session token from login
	DevIDNO    string // Device ID number
	Channel    int    // Channel number (starts from 0)
	Stream     int    // Stream type (0=main stream, 1=sub stream)
	AVType     int    // 1=live video, 2=listening
}

// GenerateRTSPLink creates a properly formatted RTSP URL for video streaming
func GenerateRTSPLink(opts RTSPLinkOptions) string {
	// Set default port from config if not specified
	if opts.ServerPort == 0 {
		opts.ServerPort = config.RTSPPort
	}

	// Default to live video if not specified
	if opts.AVType == 0 {
		opts.AVType = 1
	}

	// Format the RTSP URL according to the API documentation
	return fmt.Sprintf("rtsp://%s:%d/3/3?AVType=%d&jsession=%s&DevIDNO=%s&Channel=%d&Stream=%d",
		opts.ServerHost,
		opts.ServerPort,
		opts.AVType,
		opts.JSession,
		opts.DevIDNO,
		opts.Channel,
		opts.Stream)
}

// HLSLinkOptions contains the parameters needed to build an HLS URL

type HLSLinkOptions struct {
	ServerHost  string // HLS server hostname
	ServerPort  int    // HLS server port (default 16604)
	JSession    string // Session token from login
	DevIDNO     string // Device ID number
	Channel     int    // Channel number (starts from 0)
	Stream      int    // Stream type (0=main stream, 1=sub stream)
	RequestType int    // 1 for real-time video
}

// GenerateHLSLink creates a properly formatted HLS URL for video streaming
// HLS(HTTP Live streaming) is a streaming media transmission protocol based on HTTP, which is proposed by Apple as a protocol interaction method for transmitting audio and video.
// Provides the real- time video request address based on the HLS protocol. Currently supports h264, does not support h265.
func GenerateHLSLink(opts HLSLinkOptions) string {
	// Set default port from config if not specified
	if opts.ServerPort == 0 {
		opts.ServerPort = config.HLSPort
	}

	// Default to real-time video if not specified
	if opts.RequestType == 0 {
		opts.RequestType = 1
	}

	// Format the HLS URL according to the API documentation
	return fmt.Sprintf("https://%s:%d/hls/%d_%s_%d_%d.m3u8?jsession=%s",
		opts.ServerHost,
		opts.ServerPort,
		opts.RequestType,
		opts.DevIDNO,
		opts.Channel,
		opts.Stream,
		opts.JSession)
}

// RTMPLinkOptions contains the parameters needed to build an RTMP URL
type RTMPLinkOptions struct {
	ServerHost string // RTMP server hostname
	ServerPort int    // RTMP server port (default from config)
	JSession   string // Session token from login
	DevIDNO    string // Device ID number
	Channel    int    // Channel number (starts from 0)
	Stream     int    // Stream type (0=main stream, 1=sub stream)
	AVType     int    // 1=live video, 2=listening
}

// GenerateRTMPLink creates a properly formatted RTMP URL for video streaming
func GenerateRTMPLink(opts RTMPLinkOptions) string {
	// Set default port from config if not specified
	if opts.ServerPort == 0 {
		opts.ServerPort = config.RTMPPort
	}

	// Default to live video if not specified
	if opts.AVType == 0 {
		opts.AVType = 1
	}

	// Format the RTMP URL according to the API documentation
	return fmt.Sprintf("rtmp://%s:%d/3/3?AVType=%d&jsession=%s&DevIDNO=%s&Channel=%d&Stream=%d",
		opts.ServerHost,
		opts.ServerPort,
		opts.AVType,
		opts.JSession,
		opts.DevIDNO,
		opts.Channel,
		opts.Stream)
}

func httpGetJSON(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", "GoClient")

	resp, err := client.Do(req)
	if err != nil && isCertError(err) {
		// Retry with InsecureSkipVerify
		insecureClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
		req2, err2 := http.NewRequest("GET", url, nil)
		if err2 != nil {
			return nil, fmt.Errorf("failed to create retry request: %v", err2)
		}
		req2.Header.Set("User-Agent", "GoClient")
		resp, err = insecureClient.Do(req2)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func isCertError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "x509:") || strings.Contains(msg, "certificate signed by unknown authority")
}

func login(account, password string) (string, error) {
	url := fmt.Sprintf("%s?account=%s&password=%s", getLoginURL(), account, password)
	data, err := httpGetJSON(url)
	if err != nil {
		return "", err
	}
	var res LoginResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return "", err
	}
	if res.Result != 0 {
		return "", fmt.Errorf("login failed (result code %d)", res.Result)
	}
	return res.JSession, nil
}

func getDevices(jsession string) ([]Device, error) {
	url := fmt.Sprintf("%s?jsession=%s", getStatusURL(), jsession)

	data, err := httpGetJSON(url)
	if err != nil {
		return nil, err
	}
	var res StatusResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Onlines, nil
}

func getVehicleInfo(jsession string) (*VehicleResponse, error) {
	url := fmt.Sprintf("%s?jsession=%s", getVehicleInfoURL(), jsession)

	fmt.Printf("Requesting vehicle info from: %s\n", url)

	data, err := httpGetJSON(url)
	if err != nil {
		fmt.Printf("HTTP request error: %v\n", err)
		return nil, err
	}

	fmt.Printf("Raw response: %s\n", string(data))

	var res VehicleResponse
	if err := json.Unmarshal(data, &res); err != nil {
		fmt.Printf("JSON parsing error: %v\n", err)
		return nil, err
	}

	if res.Result != 0 {
		fmt.Printf("API error response with code: %d\n", res.Result)
		return nil, fmt.Errorf("vehicle info request failed (result code %d)", res.Result)
	}

	return &res, nil
}

func getDeviceAlarms(jsession, devIDNO string, toMap int) (*AlarmResponse, error) {
	url := fmt.Sprintf("%s?jsession=%s&DevIDNO=%s&toMap=%d", getAlarmURL(), jsession, devIDNO, toMap)

	fmt.Printf("Requesting alarms from: %s\n", url)

	data, err := httpGetJSON(url)
	if err != nil {
		fmt.Printf("HTTP request error: %v\n", err)
		return nil, err
	}

	fmt.Printf("Raw alarm response: %s\n", string(data))

	var res AlarmResponse
	if err := json.Unmarshal(data, &res); err != nil {
		fmt.Printf("JSON parsing error: %v\n", err)
		return nil, err
	}

	if res.Result != 0 {
		return nil, fmt.Errorf("alarm request failed (result code %d)", res.Result)
	}

	return &res, nil
}

func generateLinks(jsession, did, vid, account, password string) map[string]string {
	return map[string]string{
		"Web Player ID": fmt.Sprintf("%s/808gps/open/player/video.html?lang=en&devIdno=%s&account=%s&password=%s", getWebPlayerURL(), did, account, password),
		"Web Player VI": fmt.Sprintf("%s/808gps/open/player/video.html?lang=en&vehiIdno=%s&account=%s&password=%s", getWebPlayerURL(), vid, account, password),
		"Live API":      fmt.Sprintf("%s?jsession=%s&DevIDNO=%s&Chn=1&Sec=300&Label=test", getLiveAPIBaseURL(), jsession, did),
	}
}

func saveToFile(account string, allLinks map[string]map[string]string) error {
	filename := fmt.Sprintf("%s-%ddev-%s.txt", account, len(allLinks), time.Now().Format("2006-01-02"))
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	for name, links := range allLinks {
		fmt.Fprintf(f, "Device: %s\n", name)
		for k, v := range links {
			fmt.Fprintf(f, "  %s: %s\n", k, v)
		}
		fmt.Fprintln(f, strings.Repeat("-", 60))
	}
	return nil
}

func buildCompanyHierarchy(companies []struct {
	ID   int    `json:"id"`
	Name string `json:"nm"`
	PID  int    `json:"pId"`
}) map[int][]struct {
	ID   int    `json:"id"`
	Name string `json:"nm"`
	PID  int    `json:"pId"`
} {
	hierarchy := make(map[int][]struct {
		ID   int    `json:"id"`
		Name string `json:"nm"`
		PID  int    `json:"pId"`
	})

	for _, company := range companies {
		hierarchy[company.PID] = append(hierarchy[company.PID], company)
	}

	return hierarchy
}

func printCompanyTree(builder *strings.Builder, hierarchy map[int][]struct {
	ID   int    `json:"id"`
	Name string `json:"nm"`
	PID  int    `json:"pId"`
}, parentID int, prefix string) {
	children, exists := hierarchy[parentID]
	if !exists {
		return
	}

	for i, company := range children {
		isLast := i == len(children)-1

		if isLast {
			builder.WriteString(fmt.Sprintf("%s└── %s\n", prefix, company.Name))
			printCompanyTree(builder, hierarchy, company.ID, prefix+"    ")
		} else {
			builder.WriteString(fmt.Sprintf("%s├── %s\n", prefix, company.Name))
			printCompanyTree(builder, hierarchy, company.ID, prefix+"│   ")
		}
	}
}

func getStatusDescription(status EquipmentStatus) string {
	var descriptions []string

	// Build status descriptions for relevant flags
	if status.GPSValid {
		descriptions = append(descriptions, "GPS Valid")
	}

	if status.ACCStatus {
		descriptions = append(descriptions, "ACC On")
	}

	if status.LeftTurn {
		descriptions = append(descriptions, "Left Turn")
	}

	if status.RightTurn {
		descriptions = append(descriptions, "Right Turn")
	}

	if status.QuiescentState {
		descriptions = append(descriptions, "Quiescent")
	}

	if status.OverspeedState {
		descriptions = append(descriptions, "Overspeeding")
	}

	if status.BatteryStatus {
		descriptions = append(descriptions, "Battery Low")
	}

	if status.NightState {
		descriptions = append(descriptions, "Night Mode")
	}

	if status.DoorOpen {
		descriptions = append(descriptions, "Door Open")
	}

	// Add important alarms
	var alarms []string
	if status.EmergencyAlarm {
		alarms = append(alarms, "Emergency")
	}
	if status.AreaOverspeedAlarm {
		alarms = append(alarms, "Area Overspeed")
	}
	if status.FatigueDrivingReport {
		alarms = append(alarms, "Fatigue Driving")
	}
	if status.DangerousDrivingAlarm {
		alarms = append(alarms, "Dangerous Driving")
	}
	if status.VehicleTheftAlarm {
		alarms = append(alarms, "Vehicle Theft")
	}
	if status.IllegalIgnitionAlarm {
		alarms = append(alarms, "Illegal Ignition")
	}
	if status.CollisionRolloverAlarm {
		alarms = append(alarms, "Collision/Rollover")
	}

	if len(alarms) > 0 {
		descriptions = append(descriptions, fmt.Sprintf("ALARMS: %s", strings.Join(alarms, ", ")))
	}

	return strings.Join(descriptions, ", ")
}

func appIcon() fyne.Resource {
	// Return the default Fyne icon instead of trying to load a custom one
	// This avoids the PNG decoding error
	return theme.FyneLogo()
}

func logAlarmsToFile(alarms []AlarmResponseAlarm) {
	if len(alarms) == 0 {
		return
	}

	f, err := os.OpenFile("alarms.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to write log: %v\n", err)
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(f, "=== Alarm log at %s ===\n", timestamp)

	for _, alarm := range alarms {
		fmt.Fprintf(f, "Device: %s\n", alarm.DevIDNO)
		fmt.Fprintf(f, "Time: %s\n", alarm.Time)
		fmt.Fprintf(f, "Type: %d\n", alarm.Type)
		fmt.Fprintf(f, "Description: %s\n", alarm.Desc)

		if alarm.Gps.Lat != 0 && alarm.Gps.Lng != 0 {
			fmt.Fprintf(f, "Location: %.6f, %.6f\n",
				float64(alarm.Gps.Lat)/1000000.0, float64(alarm.Gps.Lng)/1000000.0)
			fmt.Fprintf(f, "Mapped Location: %s, %s\n", alarm.Gps.MLat, alarm.Gps.MLng)
			fmt.Fprintf(f, "Speed: %.1f km/h\n", float64(alarm.Gps.SP)/10.0)
		}

		status := "Unprocessed"
		if alarm.HD == 1 {
			status = "Processed"
		}
		fmt.Fprintf(f, "Status: %s\n", status)
		fmt.Fprintln(f, strings.Repeat("-", 60))
	}
}

// URL generation functions that use config
func getLoginURL() string {
	return fmt.Sprintf("%s/StandardApiAction_login.action", config.ServerURL)
}

func getStatusURL() string {
	return fmt.Sprintf("%s/StandardApiAction_getDeviceOlStatus.action", config.ServerURL)
}

func getLiveAPIBaseURL() string {
	return fmt.Sprintf("%s/StandardApiAction_realTimeVedio.action", config.ServerURL)
}

func getVehicleInfoURL() string {
	return fmt.Sprintf("%s/StandardApiAction_queryUserVehicle.action", config.ServerURL)
}

func getAlarmURL() string {
	return fmt.Sprintf("%s/StandardApiAction_vehicleAlarm.action", config.ServerURL)
}

func getWebPlayerURL() string {
	return strings.Replace(config.ServerURL, "https://", "http://", 1)
}

// getServerHostname extracts hostname from server URL for streaming services
func getServerHostname() string {
	// Remove protocol prefix (https:// or http://)
	hostname := strings.TrimPrefix(config.ServerURL, "https://")
	hostname = strings.TrimPrefix(hostname, "http://")

	// Remove any path or port if present
	parts := strings.Split(hostname, "/")
	hostname = parts[0]

	// Remove port if present (for cases like hostname:port)
	parts = strings.Split(hostname, ":")
	hostname = parts[0]

	return hostname
}

func main() {
	// Load configuration
	err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	myApp := app.New()
	myWindow := myApp.NewWindow("CMSV Video Generator Link")
	myWindow.SetIcon(appIcon())

	accountEntry := widget.NewEntry()
	accountEntry.SetPlaceHolder("Enter Account")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Enter Password")

	output := widget.NewMultiLineEntry()
	output.SetPlaceHolder("Results will appear here...")
	output.SetMinRowsVisible(15)

	// Create a dropdown for device selection
	deviceSelector := widget.NewSelect([]string{"Login first to see devices"}, func(selected string) {
		// This will be handled when a device is selected
	})
	deviceSelector.PlaceHolder = "Select Device IDNO"
	deviceSelector.Disable() // Disable until logged in

	// Create a dropdown for coordinate system selection
	coordSystems := []string{
		"0 - WGS84 (Default)",
		"1 - Google (GJ02)",
		"2 - Baidu (BD09)",
	}
	coordSystemSelector := widget.NewSelect(coordSystems, nil)
	coordSystemSelector.SetSelected(coordSystems[0])

	var allLinks map[string]map[string]string
	var deviceMap map[string]Device // Map to store device names to their IDs
	var jsessionCache string        // Store the session for reuse

	loginBtn := widget.NewButton("Login and Fetch Devices", func() {
		account := strings.TrimSpace(accountEntry.Text)
		password := strings.TrimSpace(passwordEntry.Text)

		if account == "" || password == "" {
			dialog.ShowError(fmt.Errorf("please enter both account and password"), myWindow)
			return
		}

		jsession, err := login(account, password)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Login failed: %v", err), myWindow)
			return
		}

		jsessionCache = jsession // Cache the jsession for later use

		devices, err := getDevices(jsession)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Device fetch failed: %v", err), myWindow)
			return
		}

		// Update the device selector with actual devices
		deviceOptions := []string{"All Devices"}
		deviceMap = make(map[string]Device)

		for _, d := range devices {
			deviceName := fmt.Sprintf("%s (%s)", d.VID, d.DID)
			deviceOptions = append(deviceOptions, deviceName)
			deviceMap[deviceName] = d
		}

		deviceSelector.Options = deviceOptions
		deviceSelector.Enable()
		deviceSelector.SetSelected("All Devices")

		allLinks = make(map[string]map[string]string)
		builder := strings.Builder{}

		builder.WriteString(fmt.Sprintf("Found %d devices\n\n", len(devices)))
		for _, d := range devices {
			key := fmt.Sprintf("%s (%s)", d.VID, d.DID)
			links := generateLinks(jsession, d.DID, d.VID, account, password)
			allLinks[key] = links

			builder.WriteString(fmt.Sprintf("Device: %s\n", key))
			for name, link := range links {
				builder.WriteString(fmt.Sprintf("  %s: %s\n", name, link))
			}
			builder.WriteString(strings.Repeat("-", 60) + "\n")
		}

		output.SetText(builder.String())
	})

	saveBtn := widget.NewButton("Save to File", func() {
		if allLinks == nil {
			dialog.ShowInformation("Info", "No data to save yet", myWindow)
			return
		}
		err := saveToFile(accountEntry.Text, allLinks)
		if err != nil {
			dialog.ShowError(err, myWindow)
		} else {
			dialog.ShowInformation("Success", "File saved successfully", myWindow)
		}
	})

	vehicleInfoBtn := widget.NewButton("VEHICLE INFORMATION", func() {
		account := strings.TrimSpace(accountEntry.Text)
		password := strings.TrimSpace(passwordEntry.Text)

		if account == "" || password == "" {
			dialog.ShowError(fmt.Errorf("please enter both account and password"), myWindow)
			return
		}

		jsession, err := login(account, password)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Login failed: %v", err), myWindow)
			return
		}

		jsessionCache = jsession // Cache the jsession for reuse
		fmt.Printf("Using jsession: %s\n", jsession)

		vehicleInfo, err := getVehicleInfo(jsession)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Vehicle info fetch failed: %v", err), myWindow)
			return
		}

		builder := strings.Builder{}

		// Build and display company hierarchy as a clean tree without IDs
		builder.WriteString("=== COMPANY HIERARCHY ===\n")
		hierarchy := buildCompanyHierarchy(vehicleInfo.Companys)
		printCompanyTree(&builder, hierarchy, 2, "")
		builder.WriteString("\n")

		// Display vehicle information
		builder.WriteString("=== VEHICLE INFORMATION ===\n")
		for _, vehicle := range vehicleInfo.Vehicles {
			builder.WriteString(fmt.Sprintf("Vehicle: %s (ID: %d)\n", vehicle.Name, vehicle.ID))
			builder.WriteString(fmt.Sprintf("  Company: %s\n", vehicle.PName))
			builder.WriteString(fmt.Sprintf("  Type: %s, Band: %s, Color: %s\n",
				vehicle.VehicleType, vehicle.VehicleBand, vehicle.VehicleColor))
			builder.WriteString(fmt.Sprintf("  Owner: %s\n", vehicle.OwnerName))
			builder.WriteString(fmt.Sprintf("  Engine #: %s, Frame #: %s\n",
				vehicle.EngineNum, vehicle.FrameNum))

			// Display device information for each vehicle
			builder.WriteString("  Devices:\n")
			for _, device := range vehicle.DeviceList {
				builder.WriteString(fmt.Sprintf("    - %s (%s)\n", device.ID, device.SIM))
				builder.WriteString(fmt.Sprintf("      Channels: %d, Channel Name: %s\n", device.Channels, device.ChanName))
				builder.WriteString(fmt.Sprintf("      Installed: %s\n", device.InstallTime))
			}
			builder.WriteString(strings.Repeat("-", 60) + "\n")
		}

		output.SetText(builder.String())
	})

	alarmBtn := widget.NewButton("GET DEVICE ALARMS", func() {
		if jsessionCache == "" {
			dialog.ShowError(fmt.Errorf("please login first"), myWindow)
			return
		}

		selectedDevice := deviceSelector.Selected
		if selectedDevice == "" {
			dialog.ShowError(fmt.Errorf("please select a device"), myWindow)
			return
		}

		// Get the device ID based on selection
		var deviceID string
		if selectedDevice == "All Devices" {
			deviceID = "" // Empty string means all devices
		} else if device, ok := deviceMap[selectedDevice]; ok {
			deviceID = device.DID
		} else {
			dialog.ShowError(fmt.Errorf("invalid device selection"), myWindow)
			return
		}

		// Get coordinate system selection
		toMap := 0 // Default to WGS84
		selectedCoordSystem := coordSystemSelector.Selected
		if strings.HasPrefix(selectedCoordSystem, "1 -") {
			toMap = 1 // Google
		} else if strings.HasPrefix(selectedCoordSystem, "2 -") {
			toMap = 2 // Baidu
		}

		alarmData, err := getDeviceAlarms(jsessionCache, deviceID, toMap)
		if err != nil {
			dialog.ShowError(fmt.Errorf("alarm fetch failed: %v", err), myWindow)
			return
		}

		builder := strings.Builder{}
		builder.WriteString("=== DEVICE ALARMS ===\n")

		if len(alarmData.AlarmList) == 0 {
			builder.WriteString("No alarms found for this device\n")
		} else {
			builder.WriteString(fmt.Sprintf("Found %d alarms\n\n", len(alarmData.AlarmList)))

			for _, alarm := range alarmData.AlarmList {
				// Display alarm details
				builder.WriteString(fmt.Sprintf("Device: %s\n", alarm.DevIDNO))
				builder.WriteString(fmt.Sprintf("Time: %s\n", alarm.Time))
				builder.WriteString(fmt.Sprintf("Type: %d\n", alarm.Type))
				builder.WriteString(fmt.Sprintf("Description: %s\n", alarm.Desc))

				if alarm.Gps.Lat != 0 && alarm.Gps.Lng != 0 {
					builder.WriteString(fmt.Sprintf("Location: %.6f, %.6f\n",
						float64(alarm.Gps.Lat)/1000000.0, float64(alarm.Gps.Lng)/1000000.0))
					builder.WriteString(fmt.Sprintf("Mapped Location: %s, %s\n", alarm.Gps.MLat, alarm.Gps.MLng))
					builder.WriteString(fmt.Sprintf("Speed: %.1f km/h\n", float64(alarm.Gps.SP)/10.0))
				}

				status := "Unprocessed"
				if alarm.HD == 1 {
					status = "Processed"
				}
				builder.WriteString(fmt.Sprintf("Status: %s\n", status))
				builder.WriteString(strings.Repeat("-", 60) + "\n")
			}

			// Log alarms to file for future reference
			var alarmSlice []AlarmResponseAlarm
			for _, a := range alarmData.AlarmList {
				alarmSlice = append(alarmSlice, a)
			}
			logAlarmsToFile(alarmSlice)
		}

		output.SetText(builder.String())
	})

	// Add refresh button to continuously fetch alarms
	var refreshTicker *time.Ticker
	var stopRefresh chan bool
	var refreshBtn *widget.Button
	refreshBtn = widget.NewButton("AUTO REFRESH ALARMS", func() {
		if refreshTicker != nil {
			// Stop auto-refresh
			stopRefresh <- true
			refreshTicker = nil
			refreshBtn.SetText("AUTO REFRESH ALARMS")
			dialog.ShowInformation("Auto-refresh", "Auto-refresh stopped", myWindow)
			return
		}

		// Start auto-refresh
		const timeoutSec = 5
		message := fmt.Sprintf("Start auto-refreshing alarms every %d seconds?", timeoutSec)

		dialog.ShowConfirm("Auto-refresh", message, func(start bool) {
			if !start {
				return
			}

			// Setup channels
			refreshTicker = time.NewTicker(time.Duration(timeoutSec) * time.Second)
			stopRefresh = make(chan bool)
			refreshBtn.SetText("STOP AUTO REFRESH")

			// Start refresh goroutine
			go func() {
				for {
					select {
					case <-refreshTicker.C:
						// Get device ID based on selection
						var deviceID string
						if deviceSelector.Selected == "All Devices" {
							deviceID = "" // Empty string means all devices
						} else if device, ok := deviceMap[deviceSelector.Selected]; ok {
							deviceID = device.DID
						} else {
							continue
						}

						// Get coordinate system
						toMap := 0 // Default to WGS84
						selectedCoordSystem := coordSystemSelector.Selected
						if strings.HasPrefix(selectedCoordSystem, "1 -") {
							toMap = 1 // Google
						} else if strings.HasPrefix(selectedCoordSystem, "2 -") {
							toMap = 2 // Baidu
						}

						// Fetch alarms
						alarmData, err := getDeviceAlarms(jsessionCache, deviceID, toMap)
						if err != nil {
							continue // Skip this iteration on error
						}

						// Build output
						builder := strings.Builder{}
						builder.WriteString(fmt.Sprintf("=== AUTO REFRESH ALARMS (%s) ===\n",
							time.Now().Format("15:04:05")))

						if len(alarmData.AlarmList) == 0 {
							builder.WriteString("No alarms found for this device\n")
						} else {
							builder.WriteString(fmt.Sprintf("Found %d alarms\n\n", len(alarmData.AlarmList)))

							for _, alarm := range alarmData.AlarmList {
								// Display alarm details
								builder.WriteString(fmt.Sprintf("Device: %s\n", alarm.DevIDNO))
								builder.WriteString(fmt.Sprintf("Time: %s\n", alarm.Time))
								builder.WriteString(fmt.Sprintf("Type: %d\n", alarm.Type))
								builder.WriteString(fmt.Sprintf("Description: %s\n", alarm.Desc))

								builder.WriteString(strings.Repeat("-", 60) + "\n")
							}
						}

						// Update UI thread-safely
						// Update UI on main thread
						output.SetText(builder.String())
						output.Refresh()

					case <-stopRefresh:
						if refreshTicker != nil {
							refreshTicker.Stop()
						}
						return
					}
				}
			}()
		}, myWindow)
	})

	// Add RTSP link generation button
	rtspBtn := widget.NewButton("Generate RTSP Link", func() {
		// Ensure we have a valid session and selected device
		if jsessionCache == "" {
			dialog.ShowInformation("Error", "Please login first", myWindow)
			return
		}

		selectedDevice := deviceSelector.Selected
		if selectedDevice == "" || selectedDevice == "All Devices" {
			dialog.ShowInformation("Error", "Please select a specific device", myWindow)
			return
		}

		device, ok := deviceMap[selectedDevice]
		if !ok {
			dialog.ShowError(fmt.Errorf("invalid device selection"), myWindow)
			return
		}

		// Show config dialog for RTSP parameters
		serverEntry := widget.NewEntry()
		serverEntry.SetText(getServerHostname())

		streamOptions := []string{"Main Stream (0)", "Sub Stream (1)"}
		streamSelector := widget.NewSelect(streamOptions, nil)
		streamSelector.SetSelected(streamOptions[1]) // Default to sub stream

		channelOptions := []string{"Channel 0", "Channel 1", "Channel 2", "Channel 3"}
		channelSelector := widget.NewSelect(channelOptions, nil)
		channelSelector.SetSelected(channelOptions[0]) // Default to channel 0

		configContainer := container.NewVBox(
			widget.NewLabel("Configure RTSP Stream:"),
			container.NewGridWithColumns(2,
				widget.NewLabel("Server:"),
				serverEntry,
				widget.NewLabel("Stream Type:"),
				streamSelector,
				widget.NewLabel("Channel:"),
				channelSelector,
			),
		)

		dialog.ShowCustomConfirm("RTSP Configuration", "Generate", "Cancel", configContainer, func(generate bool) {
			if !generate {
				return
			}

			// Parse channel number from selection
			channelStr := channelSelector.Selected
			channelNum := 0 // Default
			if len(channelStr) > 0 {
				channelNum, _ = strconv.Atoi(string(channelStr[len(channelStr)-1]))
			}

			// Parse stream type from selection
			streamStr := streamSelector.Selected
			streamType := 1 // Default to sub stream
			if strings.Contains(streamStr, "(0)") {
				streamType = 0 // Main stream
			}

			// Generate the RTSP link
			rtspOptions := RTSPLinkOptions{
				ServerHost: serverEntry.Text,
				ServerPort: 6604, // Default port
				JSession:   jsessionCache,
				DevIDNO:    device.DID,
				Channel:    channelNum,
				Stream:     streamType,
				AVType:     1, // Live video
			}

			rtspLink := GenerateRTSPLink(rtspOptions)

			// Show the generated link
			linkEntry := widget.NewMultiLineEntry()
			linkEntry.SetText(rtspLink)
			linkEntry.TextStyle = fyne.TextStyle{Monospace: true}

			linkContainer := container.NewVBox(
				widget.NewLabel("RTSP Link Generated:"),
				linkEntry,
				widget.NewButton("Copy to Clipboard", func() {
					myWindow.Clipboard().SetContent(rtspLink)
					dialog.ShowInformation("Copied", "RTSP link copied to clipboard", myWindow)
				}),
			)

			dialog.ShowCustom("RTSP Link", "Close", linkContainer, myWindow)
		}, myWindow)
	})

	// Add RTMP link generation button
	rtmpBtn := widget.NewButton("Generate RTMP Link", func() {
		// Ensure we have a valid session and selected device
		if jsessionCache == "" {
			dialog.ShowInformation("Error", "Please login first", myWindow)
			return
		}

		selectedDevice := deviceSelector.Selected
		if selectedDevice == "" || selectedDevice == "All Devices" {
			dialog.ShowInformation("Error", "Please select a specific device", myWindow)
			return
		}

		device, ok := deviceMap[selectedDevice]
		if !ok {
			dialog.ShowError(fmt.Errorf("invalid device selection"), myWindow)
			return
		}

		// Show config dialog for RTMP parameters
		serverEntry := widget.NewEntry()
		serverEntry.SetText(getServerHostname())

		streamOptions := []string{"Main Stream (0)", "Sub Stream (1)"}
		streamSelector := widget.NewSelect(streamOptions, nil)
		streamSelector.SetSelected(streamOptions[1]) // Default to sub stream

		channelOptions := []string{"Channel 0", "Channel 1", "Channel 2", "Channel 3"}
		channelSelector := widget.NewSelect(channelOptions, nil)
		channelSelector.SetSelected(channelOptions[0]) // Default to channel 0

		configContainer := container.NewVBox(
			widget.NewLabel("Configure RTMP Stream:"),
			container.NewGridWithColumns(2,
				widget.NewLabel("Server:"),
				serverEntry,
				widget.NewLabel("Channel:"),
				channelSelector,
				widget.NewLabel("Stream Type:"),
				streamSelector,
			),
		)

		dialog.ShowCustomConfirm("RTMP Configuration", "Generate", "Cancel", configContainer, func(generate bool) {
			if !generate {
				return
			}

			// Parse channel number from selection
			channelStr := channelSelector.Selected
			channelNum := 0 // Default
			if len(channelStr) > 0 {
				channelNum, _ = strconv.Atoi(string(channelStr[len(channelStr)-1]))
			}

			// Parse stream type from selection
			streamStr := streamSelector.Selected
			streamType := 1 // Default to sub stream
			if strings.Contains(streamStr, "(0)") {
				streamType = 0 // Main stream
			}

			// Generate the RTMP link
			rtmpOptions := RTMPLinkOptions{
				ServerHost: serverEntry.Text,
				ServerPort: 6604, // Default port for RTMP
				JSession:   jsessionCache,
				DevIDNO:    device.DID,
				Channel:    channelNum,
				Stream:     streamType,
				AVType:     1, // Live video
			}

			rtmpLink := GenerateRTMPLink(rtmpOptions)

			// Show the generated link
			linkEntry := widget.NewMultiLineEntry()
			linkEntry.SetText(rtmpLink)
			linkEntry.TextStyle = fyne.TextStyle{Monospace: true}

			linkContainer := container.NewVBox(
				widget.NewLabel("RTMP URL:"),
				linkEntry,
			)

			dialog.ShowCustom("RTMP Link", "Close", linkContainer, myWindow)
		}, myWindow)
	})

	// Add HLS link generation button
	hlsBtn := widget.NewButton("Generate HLS Link", func() {
		// Ensure we have a valid session and selected device
		if jsessionCache == "" {
			dialog.ShowInformation("Error", "Please login first", myWindow)
			return
		}

		selectedDevice := deviceSelector.Selected
		if selectedDevice == "" || selectedDevice == "All Devices" {
			dialog.ShowInformation("Error", "Please select a specific device", myWindow)
			return
		}

		device, ok := deviceMap[selectedDevice]
		if !ok {
			dialog.ShowError(fmt.Errorf("invalid device selection"), myWindow)
			return
		}

		// Show config dialog for HLS parameters
		serverEntry := widget.NewEntry()
		serverEntry.SetText(getServerHostname())

		streamOptions := []string{"Main Stream (0)", "Sub Stream (1)"}
		streamSelector := widget.NewSelect(streamOptions, nil)
		streamSelector.SetSelected(streamOptions[1]) // Default to sub stream

		channelOptions := []string{"Channel 0", "Channel 1", "Channel 2", "Channel 3"}
		channelSelector := widget.NewSelect(channelOptions, nil)
		channelSelector.SetSelected(channelOptions[0]) // Default to channel 0

		configContainer := container.NewVBox(
			widget.NewLabel("Configure HLS Stream:"),
			container.NewGridWithColumns(2,
				widget.NewLabel("Server:"),
				serverEntry,
				widget.NewLabel("Channel:"),
				channelSelector,
				widget.NewLabel("Stream Type:"),
				streamSelector,
			),
		)

		dialog.ShowCustomConfirm("HLS Configuration", "Generate", "Cancel", configContainer, func(generate bool) {
			if !generate {
				return
			}

			// Parse channel number from selection
			channelStr := channelSelector.Selected
			channelNum := 0 // Default
			if len(channelStr) > 0 {
				channelNum, _ = strconv.Atoi(string(channelStr[len(channelStr)-1]))
			}

			// Parse stream type from selection
			streamStr := streamSelector.Selected
			streamType := 1 // Default to sub stream
			if strings.Contains(streamStr, "(0)") {
				streamType = 0 // Main stream
			}

			// Generate the HLS link
			hlsOptions := HLSLinkOptions{
				ServerHost:  serverEntry.Text,
				ServerPort:  16604, // Default port for HLS
				JSession:    jsessionCache,
				DevIDNO:     device.DID,
				Channel:     channelNum,
				Stream:      streamType,
				RequestType: 1, // Real-time video
			}

			hlsLink := GenerateHLSLink(hlsOptions)

			// Show the generated link
			linkEntry := widget.NewMultiLineEntry()
			linkEntry.SetText(hlsLink)
			linkEntry.TextStyle = fyne.TextStyle{Monospace: true}

			// Add HTML video player code
			htmlCode := fmt.Sprintf(`<video controls preload="none" width="352" height="288" data-setup="{}">
    <source src="%s" type="application/x-mpegURL">
</video>`, hlsLink)

			htmlEntry := widget.NewMultiLineEntry()
			htmlEntry.SetText(htmlCode)
			htmlEntry.TextStyle = fyne.TextStyle{Monospace: true}

			linkContainer := container.NewVBox(
				widget.NewLabel("HLS URL:"),
				linkEntry,
				widget.NewLabel("HTML Video Player Code:"),
				htmlEntry,
			)

			dialog.ShowCustom("HLS Link", "Close", linkContainer, myWindow)
		}, myWindow)
	})

	// Create the final UI layout
	// Create the final UI layout
	// Create the final UI layout
	content := container.NewVBox(
		container.NewGridWithColumns(2,
			widget.NewLabel("Account:"),
			accountEntry,
			widget.NewLabel("Password:"),
			passwordEntry,
		),
		loginBtn,
		deviceSelector,
		container.NewGridWithColumns(6, // Changed from 5 to 6 columns
			vehicleInfoBtn,
			alarmBtn,
			refreshBtn,
			rtspBtn,
			hlsBtn,
			rtmpBtn, // Added RTMP button
		),
		widget.NewLabel("Coordinate System:"),
		coordSystemSelector,
		saveBtn,
		output,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}
