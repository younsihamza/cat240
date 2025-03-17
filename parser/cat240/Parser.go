package cat240

import (
	"fmt"
	"math"
)

func Parser(data []byte) (*ValidData, error) {
	if data[0] != 240 {
		return nil, fmt.Errorf("Invalid data")
	}
	if (int(data[1]) << 8) + int(data[2])  != len(data) {
		return nil, fmt.Errorf("Invalid data")
	}
	cureentByte := 3
	FSPEC := fmt.Sprintf("%08b", int64(int(data[cureentByte])))
	if FSPEC[7] == '1' {
		cureentByte++
		FSPEC = fmt.Sprintf("%016b", int64(int(data[3]) << 8 + int(data[4])))
	}
	cureentByte++
	dataObject , err := convertToVideoDataItem(FSPEC, data[cureentByte:])
	if err != nil {
		return nil, err
	}
	return dataObject, nil
}

// convertToVideoDataItem is a function that converts the data to a VideoDataItem struct
func convertToVideoDataItem(fspec string, data []byte) (*ValidData, error) {
	var item VideoDataItem
	for i, val := range fspec {
		
		i = i + 1
		if i == 8  && val != '1' {
			return nil, fmt.Errorf("Invalid FSPEC")
		}
		if val == '1' {
			switch i {
			case 1:
				if len(data) < 2 {
					return nil, fmt.Errorf("Invalid data")
				}
				item.DataSourceIndentifier = &DataSourceIndentifier{}
				item.DataSourceIndentifier.SAC = int(data[0])
				item.DataSourceIndentifier.SIC = int(data[1])
				data = data[2:]
			case 2:
				if len(data) < 1 {
					return nil, fmt.Errorf("Invalid data")
				}
				item.MessageType = int(data[0])
				data = data[1:]
			case 3:
				if len(data) < 4 {
					return nil, fmt.Errorf("Invalid data")
				}
				item.RecodeHeader = int(data[0]) << 24 + int(data[1]) << 16 + int(data[2]) << 8 + int(data[3])
				data = data[4:]
			case 4:
				if len(data) < 1 {
					return nil, fmt.Errorf("Invalid data")
				}
				item.VideoSummary = data[1 : int(data[0]) + 1]
				data = data[int(data[0]) + 1:]
			case 5:
				if len(data) < 12 {
					return nil, fmt.Errorf("Invalid data")
				}
				item.VideoHeaderNano = &VideosHeaders{
					StartAzimuth : float64(int(data[0]) << 8 + int(data[1])) * 360.0 / 65535.0,
					EndAzimuth : float64(int(data[2]) << 8 + int(data[3])) * 360.0 / 65535.0,
					StartRange : int(data[4]) << 24 + int(data[5]) << 16 + int(data[6]) << 8 + int(data[7]),
					CellDuration : float64(int(data[8]) << 24 + int(data[9]) << 16 + int(data[10]) << 8 + int(data[11])) * math.Pow(10, -9),
				}
				data = data[12:]
			case 6:
				if len(data) < 12 {
					return nil, fmt.Errorf("Invalid data")
				}
				item.VideoHeaderFemto = &VideosHeaders{
					StartAzimuth : float64(int(data[0]) << 8 + int(data[1])) * 360.0 / 65535.0,
					EndAzimuth : float64(int(data[2]) << 8 + int(data[3])) * 360.0 / 65535.0,
					StartRange : int(data[4]) << 24 + int(data[5]) << 16 + int(data[6]) << 8 + int(data[7]),
					CellDuration : float64(int(data[8]) << 24 + int(data[9]) << 16 + int(data[10]) << 8 + int(data[11])) * math.Pow(10, -15),
				}
				data = data[12:]
			case 7:
				if len(data) < 2 {
					return nil, fmt.Errorf("Invalid data")
				}
				item.VideoCellsResolution = &VideoCellsResolution{
					CompressionIndicator : data[0],
					BitResolution : int(data[1]),
				}
				data = data[2:]
			case 9:
				if len(data) < 5 {
					return nil, fmt.Errorf("Invalid data")
				}
				item.VideoOctetsVideoCellCounters = &VideoOctetsVideoCellCounters{
					ValidOctetsInVideoBlock : int(data[0]) << 8 + int(data[1]),
					ValidCellsInVideoBlock : int(data[2]) << 16 + int(data[3]) << 8 + int(data[4]),
				}
				data = data[5:]
			case 10:
				if len(data) < 1  || int(data[0]) * 4 + 1 > len(data) {
					return nil, fmt.Errorf("Invalid data")
				}
				item.VideoBlockLowDataVolume = data[1 : int(data[0])*4+1]
				data = data[int(data[0]) * 4 + 1:]
			case 11:
				if len(data) < 1  || int(data[0]) * 64 + 1 > len(data) {
					return nil, fmt.Errorf("Invalid data")
				}
				item.VideoBlockMediumDataVolume = data[1 : int(data[0])*64+1]
				data = data[int(data[0]) * 64 + 1:]
			case 12:
				if len(data) < 1  || int(data[0]) * 256 + 1 > len(data) {
					return nil, fmt.Errorf("Invalid data")
				}
				item.VideoBlockHighDataVolume = data[1 : int(data[0]) * 256 + 1]
				data = data[int(data[0]) * 256 + 1:]
			case 13:
				if len(data) < 3 {
					return nil, fmt.Errorf("Invalid data")
				}
				item.TimeOfDay = float64(int(data[0]) << 16 + int(data[1]) << 8 + int(data[2]))
				data = data[3:]
			case 14:
				if len(data) < 1 {
					return nil, fmt.Errorf("Invalid data")
				}
				item.ReservedExpansionField = string(data[0])
				data = data[1:]
			case 15:
				if len(data) < 1  || int(data[0]) > len(data) {
					return nil, fmt.Errorf("Invalid data")
				}
				item.SpecialPurposeField = data[1:int(data[0])]
				data = data[int(data[0]):]
			}
		}
	}

	return validator(item)
}
// validate_DAta is a function that validates the data
func validator(item VideoDataItem) (*ValidData, error) {
	var validData ValidData
	//check the mandaroity fields in the message cat 240 
	if item.DataSourceIndentifier == nil || item.MessageType == 0 || item.RecodeHeader == 0 ||
		(item.VideoHeaderNano == nil && item.VideoHeaderFemto == nil) || item.VideoCellsResolution == nil ||
		item.VideoOctetsVideoCellCounters == nil || 
		(item.VideoBlockLowDataVolume == nil && item.VideoBlockMediumDataVolume == nil && item.VideoBlockHighDataVolume == nil){
		return nil, fmt.Errorf("Invalid data")
	}
	validData.DataSourceIndentifier = item.DataSourceIndentifier
	validData.MessageType = item.MessageType
	validData.RecodeHeader = item.RecodeHeader
	if item.VideoHeaderNano != nil {
		validData.VideoHeader = item.VideoHeaderNano
	} else if item.VideoHeaderFemto != nil {	
		validData.VideoHeader = item.VideoHeaderFemto
	}
	validData.VideoCellsResolution = item.VideoCellsResolution
	validData.VideoOctetsVideoCellCounters = item.VideoOctetsVideoCellCounters
	if item.VideoBlockLowDataVolume != nil {
		validData.VideoBlock = item.VideoBlockLowDataVolume
	} else if item.VideoBlockMediumDataVolume != nil {
		validData.VideoBlock = item.VideoBlockMediumDataVolume
	} else if item.VideoBlockHighDataVolume != nil {
		validData.VideoBlock = item.VideoBlockHighDataVolume
	}
	validData.TimeOfDay = item.TimeOfDay
	validData.ReservedExpansionField = item.ReservedExpansionField
	validData.SpecialPurposeField = item.SpecialPurposeField
	return &validData, nil
}
