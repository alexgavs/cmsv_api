# API Documentation

## User Login

### Interface Description
User Login

The credential obtained by this interface is used for subsequent access to the business interface.
The validity period of the voucher is **4 hours**. The voucher can be reused during the validity period, and the validity time will be refreshed if the interface is continuously used. Do not refresh the voucher frequently.

**NOTE:**
The second login of the same account will not invalidate the jsession generated during the first login.

### Request Example
```
https://v7.cmsv8.com/StandardApiAction_login.action?account=cmsv6&password=cmsv6
```

### Request Parameters
| Parameter name | Parameter type | Must | Example | Description     |
|---------------|---------------|------|---------|----------------|
| account       | string        | Yes  | cmsv6   | User account   |
| password      | string        | Yes  | cmsv6   | User password  |

### Return Parameters
| Parameter Name | Parameter Type | Meaning of Parameters |
|----------------|---------------|----------------------|
| result         | number        | The interface status code: 0 is normal, other values indicate errors. (See Error Code Description) |
| jsession       | string        | Session Number        |
| pri            | string        | User Permissions      |
| account_name   | string        | User Account          |
| JSESSIONID     | string        | Session Number        |

### Return Example
```
{
  "result": 0,
  "jsession": "66d754dd7f41473dbd2",
  "pri": ",1,2,21,24,25,26,27,28,29,210,211,212,213,214,241,242,215,216,217,676,282,283,284,285,218,219,220,221,222,223,224,225,226,23,227,228,229,230,231,232,233,234,22,235,236,257,258,259,260,237,238,671,672,239,240,243,244,245,248,261,262,264,263,265,266,267,268,272,3,31,32,33,34,35,36,37,38,39,40,41,310,311,318,319,320,321,4,41,42,43,44,5,6,7,611,612,613,621,622,623,624,625,626,627,628,629,630,631,641,651,652,653,654,655,656,657,658,659,660,661,663,664,997,998,7,121,122,123,124,125,126,127,128,129,950,681,607,608,615,616,617,290,291,292,293,294,295,296,297,298,299,130,650,677,831,888,682,683,684,685,686,1018,1019,1020,1100,305,690,1023,810,820,841,300,2000,2001,2002,2003,2004,286,287,288,289,275,276,277,16,161,162,163,165,166,274,800,15,151,152,153,154,155,156,157,158,159,160,270,271,273,278,279,281,1511,5- 1,5- 41,5- 42,5- 4,5- 11,5- 14,5- 5,5- 58,5- 39,5- 40,5- 15,5- 55,5- 56,5- 31,5- 19,5- 25,5- 6,5- 20,5- 21,5- 17,5- 22,5- 2,5- 3,5- 7,5- 12,5- 27,5- 28,5- 29,5- 49,5- 34,5- 36,5- 52,5- 46,5- 47,5- 50,5- 13,5- 16,5- 24,5- 26,5- 32,5- 33,5- 43,5- 44,5- 45,5- 48,5- 53,5- 9,5- 8,5- 10,5- 30,5- 37,5- 18,5- 23,5- 38,5- 51,5- 54,5- 57,5- 60,50,501,502,503,504,505,506,507,100,18,19,20,164,167,168,169,170,171,172,173,174,175,176,177,178,179,900,906,901,902,908,918,919,920,921,922,923,909,910,46,47,48,49,643,644,646,903,904,905,911,850,851,852,853,312,9,91,92,93,94,95,96,97,1000,1001,1002,1003,1004,1005,1006,1007,1008,1009,1010,1011,1012,1013,1014,1015,313,665,666,667,668,669,670,673,674,675,55,551,552,553,554,555,556,557,558,559,560,561,562,1017,60,601,602,603,604,605,606,609,610,619,632,633,634,635,636,637,638,642,640,6001,645,6002,6003,6004,614,618,620,889,639,2050,- 1,- 2,2005,1021,1022 ,",
  "account_name": "cmsv6",
  "JSESSIONID": "66d754dd7f41473dbd2"
}
```

## User Logoff

### Interface Description
User Logoff Login

### Request Example
```
https://v7.cmsv8.com/StandardApiAction_logout.action?jsession=66d754dd7f41473dbd2
```

### Request Parameters
| Parameter name | Parameter type | Must | Example                  | Description     |
|---------------|---------------|------|--------------------------|----------------|
| jsession      | string        | Yes  | 66d754dd7f41473dbd2      | Session number  |

### Return Parameters
| Parameter name | Parameter type | Parameter meaning |
|---------------|---------------|-------------------|
| result        | number        | The interface status code: 0 is normal, other values indicate errors. (See Error Code Description) |

### Return Example
```
{
  "result": 0
}
```

## Get user vehicle information

### Interface Description
Get user vehicle information

### Request Example
```
https://v7.cmsv8.com/StandardApiAction_queryUserVehicle.action?jsession=66d754dd7f41473dbd2
```

### Request Parameters
| Parameter name | Parameter type | Must | Example                  | Parameter meaning |
|---------------|---------------|------|--------------------------|-------------------|
| jsession      | string        | Yes  | 66d754dd7f41473dbd2      | Session number    |

### Return Parameters
| Parameter Name | Parameter Type | Meaning of Parameters |
|----------------|---------------|----------------------|
| result         | number        | Interface status code, 0 is normal, other values indicate errors. (See Error Code Description) |
| companys       | array         | List of companies/fleets |
| id             | number        | Vehicle ID or company ID |
| nm             | string        | License plate number or company name |
| pId            | number        | Company or fleet ID |
| vehicles       | array         | List of vehicles |
| ic             | number        | Number of IO |
| pid            | number        | Equipment company |
| pnm            | string        | Company name |
| abbr           | string        | Abbreviation |
| dl             | array         | Device list |
| ...            | ...           | Many more fields as described above |

### Return Example
```
{
  "result": 0,
  "companys": [
    { "id": 3, "nm": "test11", "pId": 1 },
    { "id": 4, "nm": "testce", "pId": 3 },
    { "id": 1, "nm": "test", "pId": 10 }
  ],
  "vehicles": [
    {
      "id": 28979,
      "nm": "S66666",
      "ic": 6,
      "pid": 1,
      "pnm": "test",
      "abbr": "",
      "dl": [
        {
          "id": "013300000001",
          "pid": 1,
          "dt": null,
          "cc": 4,
          "cn": "CH1,CH2,CH3,CH4",
          "ic": 0,
          "io": "",
          "outc": null,
          "outn": null,
          "tc": 0,
          "tn": "",
          "sim": null,
          "md": 1513,
          "st": null,
          "nflt": null,
          "us": 0,
          "sdc": null,
          "did": 33,
          "vt": null,
          "isb": null,
          "srl": "",
          "ptt": null,
          "gps": null,
          "fp": null,
          "tkc": null,
          "ist": "2024-09-09 16:19:52",
          "ol": null,
          "lt": null
        }
      ],
      "pt": "黄牌",
      "vehiColor": null,
      "status": 0,
      "vehiBand": "",
      "vehiType": null,
      "vehiUse": "",
      "dateProduct": -28800000,
      "icon": 6,
      "chnCount": 4,
      "chnName": "CH1,CH2,CH3,CH4",
      "ioInCount": 0,
      "ioInName": "",
      "ioOutCount": 0,
      "ioOutName": "",
      "tempCount": 1,
      "tempName": "0|TEMP_1",
      "payEnable": null,
      "payBegin": 1725811200000,
      "payEnd": null,
      "payMonth": null,
      "payDelayDay": 0,
      "safeDate": null,
      "drivingNum": "",
      "drivingDate": -28800000,
      "operatingNum": "",
      "operatingDate": -28800000,
      "repairDate": null,
      "stlTm": 1725811200000,
      "moreId": null,
      "vehicleGrade": "",
      "approvedNumber": null,
      "approvedLoad": null,
      "vehicleType": 0,
      "installTire": 0,
      "tireBrand": "",
      "tireModel": "",
      "installAdas": 0,
      "adasBrand": "",
      "adasModel": "",
      "installDsm": 0,
      "dsmBrand": "",
      "dsmModel": "",
      "installBlind": 0,
      "blindBrand": "",
      "blindModel": "",
      "installLca": 0,
      "lcaBrand": "",
      "lcaModel": "",
      "installOM": 0,
      "engineNum": "",
      "frameNum": "",
      "ownerName": null,
      "lineId": null,
      "linkPeople": "",
      "linkPhone": "",
      "datePurchase": -28800000,
      "dateAnnualSurvey": -28800000,
      "speedLimit": 120,
      "linesOperation": "",
      "operatingId": null,
      "industry": null,
      "carType": null,
      "carPlace": null,
      "param1": "",
      "param2": "",
      "param3": "",
      "param4": "",
      "roleId": null,
      "area": "",
      "code": "",
      "nuclearAuthority": "",
      "legal": "",
      "legalPhone": "",
      "legalAddress": "",
      "introduction": "",
      "serialNum": null,
      "loginPwd": null,
      "allowLogin": 0,
      "mileCoefficient": null,
      "remark": "",
      "vehicleModel": "",
      "engineModel": "",
      "axesNumber": null,
      "totalWeight": null,
      "quasiTractionMass": null,
      "longOutlineDimensions": null,
      "wideOutlineDimensions": null,
      "highOutlineDimensions": null,
      "longInsideDimension": null,
      "wideInnerDimensions": null,
      "highInsideDimensions": null,
      "ombrand": "",
      "ommodel": ""
    }
  ]
}
```

## Common Error Codes

### Web Error Code Description
| Error Code | Description |
|------------|-------------|
| 1          | The username or password is invalid |
| 2          | The username or password is invalid |
| 3          | User disabled |
| 4          | The user has expired |
| 5          | Session does not exist |
| 6          | System exception |
| 7          | The request parameters are incorrect |
| 8          | No permission to operate the vehicle or equipment |
| 9          | The start time must not be greater than the end time |
| 10         | Query time out of range |
| 11         | The video download task already exists |
| 12         | Account already exists |
| 13         | No permission to operate |
| 14         | Number of managed devices (maximum number of additions reached) |
| 15         | Device already exists |
| 16         | Vehicle already exists |
| 17         | Device already in use |
| 18         | Vehicle not present |
| 19         | Device does not exist |
| 20         | The device does not belong to the current company |
| 21         | The number of registered devices does not match. Check whether the number of registered devices exceeds the number of registered devices. |
| 24         | Network connection exception |
| 25         | Rule name already exists |
| 26         | Rule does not exist |
| 27         | Information does not exist |
| 28         | Session number already exists |
| 29         | Company does not exist |
| 32         | Device not online |
| 34         | Single sign-on user, already logged in |

### Server Error Code Description (return parameters include: "cmsserver":1)
| Error Code | Description |
|------------|-------------|
| 2          | The username or password is invalid |
| 3          | Invalid username or password |
| 4          | User disabled |
| 5          | Information does not exist |
| 6          | Unknown error |
| 7          | Name already in use |
| 21         | Device does not exist |
| 22         | No feedback received from the device |
| 23         | Device not online |
| 26         | Device connection lost |
| 27         | No storage path defined |

## Get device online status

### Interface Description
Get device online status

### Request Example
```
https://v7.cmsv8.com/StandardApiAction_getDeviceOlStatus.action?jsession=66d754dd7f41473dbd2&vehiIdno=S66666&status=1
```

### Request Parameters
| Parameter name | Parameter type | Required | Example         | Parameter meaning |
|---------------|---------------|----------|-----------------|-------------------|
| jsession      | string        | Yes      | 66d754dd7f41473dbd2 | Session number    |
| devIdno       | string        | No       | 01330000001     | Device number(s), comma-separated. If empty, license plate number is used. |
| vehiIdno      | string        | No       | S66666          | License plate number(s), comma-separated. If empty, queries all authorized equipment for the account. |
| status        | number        | No       | 1               | Online status: 0 not online, 1 online, empty queries all. |

### Return Parameters
| Parameter Name | Parameter Type | Meaning of Parameters |
|---------------|---------------|----------------------|
| result        | number        | Interface status code, 0 is normal, other values indicate errors. (See Error Code Description) |
| onlines       | array         | List of device online status information |
| did           | string        | Equipment number |
| vid           | string        | License plate number (empty if queried by equipment number) |
| online        | number        | Online status: 1 means online, otherwise not online |

### Return Example
```
{
  "result": 0,
  "onlines": [
    {
      "vid": "S66666",
      "online": 1,
      "abbr": "",
      "did": "013300000001"
    }
  ]
}
```

## Get real-time device status

### Interface Description
Get real-time device status

### Request Example
```
https://v7.cmsv8.com/StandardApiAction_getDeviceStatus.action?jsession=66d754dd7f41473dbd2&vehiIdno=S66666
```

### Request Parameters
| Parameter name | Parameter type | Required | Example                | Parameter meaning |
|---------------|---------------|----------|------------------------|-------------------|
| jsession      | string        | Yes      | 66d754dd7f41473dbd2    | Session number    |
| devIdno       | string        | No       | 01330000001,01330000002| Device number(s), comma-separated. If empty, license plate number is used. |
| vehiIdno      | string        | No       | S66666                 | License plate number(s), comma-separated. If both are blank, queries all authorized devices of the account. |
| geoaddress    | number        | No       | 0                      | 1: provide geographic location resolution, otherwise not resolved |
| driver        | number        | No       | 1                      | 1: query driver info, other/not passed: not queried |
| toMap         | number        | No       | 0                      | 0: WGS84, 1: Google (gj02), 2: Baidu (bd09) |
| language      | string        | No       | en                     | Language for longitude/latitude analysis (zh for Chinese) |

### Return Parameters
| Parameter name | Parameter type | Description |
|---------------|---------------|-------------|
| id            | string        | Device number |
| vid           | string        | License plate |
| lng           | number        | Longitude (divide by 1,000,000 for true value) |
| lat           | number        | Latitude (divide by 1,000,000 for true value) |
| ft            | number        | Type of manufacturer |
| sp            | number        | Speed (km/h, divide by 10) |
| ol            | number        | Online status: 1 means online, otherwise not online |
| gt            | string        | Locate upload time |
| pt            | number        | Communication protocol type |
| dt            | number        | Hard disk type: 1 SD card, 2 hard disk, 3 SSD card |
| ac            | number        | Audio type |
| net           | number        | Network type: 0 3G, 1 WIFI, 2 wired, 3 4G, 4 5G |
| gw            | string        | Gateway server number |
| s1-s4         | number        | Device status description (see details) |
| t1-t4         | number        | Temperature sensors 1-4 |
| hx            | number        | North direction (0-360 degrees) |
| mlng          | string        | Converted longitude |
| mlat          | string        | Converted latitude |
| pk            | number        | Parking duration (seconds) |
| lc            | number        | Mileage (meters) |
| yl            | number        | Oil quantity (liters, divide by 100) |
| viceYl        | number        | Secondary oil quantity (liters, divide by 100) |
| ps            | string        | Resolved geographic location or (converted longitude, latitude) |
| tsp           | number        | Tachograph speed (km/h, divide by 10) |
| dn            | string        | Driver name |
| jn            | string        | Driver certificate code |
| lt            | number        | Login type: 0-linux, 1-windows, 2-web, 3-Android, 4-ios |
| ust           | number        | Usage status: 0 normal, 1 maintenance, 2 disabled, 3 overdue |
| sn            | number        | Number of satellites |
| lg            | number        | 2: long positioning (see 808-2019 protocol) |
| rt            | string        | Server time to receive location (valid when lg=2) |
| ...           | ...           | Many more fields as described above |

### Return Example
```
{
  "result": 0,
  "status": [
    {
      "id": "013300000001",
      "net": 3,
      "gw": "G1",
      "ol": 1,
      "s1": 805309827,
      "s2": 20480,
      "s3": 65280,
      "s4": 0,
      "t1": 0,
      "t2": 0,
      "t3": 0,
      "t4": 0,
      "yl": 0,
      "sp": 0,
      "hx": 0,
      "lng": 113712944,
      "lat": 23004510,
      "mlng": "113.718131",
      "mlat": "23.001755",
      "ps": "23.004510,113.712944",
      "pk": 0,
      "lc": 0,
      "gt": "2024-12-07 11:58:30.0",
      "pt": 6,
      "dt": 2,
      "ac": 0,
      "ft": 0,
      "vid": "S66666",
      "lt": 0,
      "tsp": 0,
      "abbr": "",
      "lg": 2,
      "sn": 0,
      "viceYl": 0
      // ...other fields...
    }
  ]
}
```

### Map Example

**URL:**
```
https://v7.cmsv8.com/808gps/open/map/vehicleMap.html?jsession=66d754dd7f41473dbd2&devIdno=013300000001&lang=en
```

**Parameter Description**
| Parameter name | Parameter type | Required | Example           | Parameter meaning |
|---------------|---------------|----------|-------------------|-------------------|
| jsession      | string        | No       | 66d754dd7f41473dbd2 | Session number, if blank, use username and password |
| devIdno       | string        | No       | 0133000000001     | Device number, if empty, use license plate number |
| vehiIdno      | string        | No       | S66666            | License plate number, if device number is empty |
| lang          | string        | No       | en                | Language: en for English, otherwise Chinese |

