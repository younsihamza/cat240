package cat240

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"math"
)

type BlockData struct {
	Longtitude 		float64 	"json:longtitude"
	Latitude 		float64		"json:latitude"
	Intencity 		int			"json:intencity"
	StartAzimuth 	float64		"json:start_azimuth"
	EndAzimuth 		float64		"json:end_azimuth"
	StartRange 		float64		"json:start_range"
}

func Decode (data *ValidData) map[string]interface{}{
	if checkHightOrderBit(data.VideoCellsResolution.CompressionIndicator) {
		data.VideoBlock = decompresData(data.VideoBlock)
	}
	return toGeoJson(coordinateTransformation(data),data.VideoHeader.StartAzimuth,  data.VideoHeader.EndAzimuth, data.VideoHeader.StartRange) 
}


	func coordinateTransformation(data *ValidData) *[]BlockData {
		var coordinateHold = []BlockData{}
		speedOfLight := 299792458.0 // speed of light in meters  per second
		rangeCell := data.VideoHeader.CellDuration * speedOfLight / 2.0 
		azimuthIncrement := (data.VideoHeader.EndAzimuth - data.VideoHeader.StartAzimuth) / float64(data.VideoOctetsVideoCellCounters.ValidCellsInVideoBlock)
		fmt.Println(rangeCell * float64(data.VideoHeader.StartRange))
		currentRange := rangeCell * float64(len(data.VideoBlock) - 1 + data.VideoHeader.StartRange)
		currentAzimuth := data.VideoHeader.StartAzimuth + azimuthIncrement * float64(len(data.VideoBlock)-1)
		x, y := polarToCartesian(currentRange, currentAzimuth)
		lat, longtitud := CartesianToGeo(51.754245,-1.356208, x, y)
		coordinateHold = append(coordinateHold, BlockData{Longtitude:longtitud, Latitude:lat, Intencity:0, StartAzimuth:data.VideoHeader.StartAzimuth, EndAzimuth:data.VideoHeader.EndAzimuth, StartRange:currentRange})
		// fmt.Println(coordinateHold)
		for i := 0; i < len(data.VideoBlock)-1; i++ {
			if int(data.VideoBlock[i]) < 50 {
				continue
			}
			currentRange := rangeCell * float64(i + data.VideoHeader.StartRange)
			currentAzimuth := data.VideoHeader.StartAzimuth + azimuthIncrement * float64(i)
			x, y := polarToCartesian(currentRange, currentAzimuth)
			lat, longtitud := CartesianToGeo(51.754245,-1.356208, x, y)
			coordinateHold = append(coordinateHold, BlockData{longtitud, lat, int(data.VideoBlock[i]), data.VideoHeader.StartAzimuth, data.VideoHeader.EndAzimuth, currentRange})
		}
		return &coordinateHold
	}
	
	
	func CartesianToGeo(originLat, originLon, x, y float64) (float64, float64) {
		EarthRadius :=  6378137.0 // Earth radius in meters .
		// Convert the origin latitude to radians.
		originLatRad := originLat * math.Pi / 180.0
		// Calculate the angular offsets in radians.
		deltaLatRad := y / EarthRadius
		deltaLonRad := x / (EarthRadius * math.Cos(originLatRad))
		
		// Convert the angular offsets from radians to degrees.
		deltaLatDeg := deltaLatRad * 180.0 / math.Pi
		deltaLonDeg := deltaLonRad * 180.0 / math.Pi
		
		// Calculate the new geographic coordinates.
		newLat := originLat + deltaLatDeg
		newLon := originLon + deltaLonDeg
		return newLat, newLon
	}
	
	func polarToCartesian(rangeCell float64, azimuth float64) (float64, float64) {
		azimuthRad := azimuth * math.Pi / 180
		var x float64 = rangeCell * math.Sin(azimuthRad)
		var y float64 = rangeCell * math.Cos(azimuthRad)
		return x, y
	}
	
	
	func decompresData(data []byte) ([]byte){
		r, err := zlib.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil
		}
		defer r.Close()
		data ,err  = io.ReadAll(r)
		if err != nil {
			return nil
		}
		return data
	}
	
	func checkHightOrderBit(b byte) bool {
		return (b >> 7) == 0x01
	}
	
	func bitResolution(data *ValidData, bit_per_cell int) {
		var video_data []byte
		switch bit_per_cell {
		case 1, 2, 4, 8:
			for i := 0; i < len(data.VideoBlock); i++ {
				for j := 0; j < 8; j += bit_per_cell {
					video_data = append(video_data, data.VideoBlock[i] >> (8 - bit_per_cell - j) & (1 << bit_per_cell - 1))
				}
			}
		}
		data.VideoBlock = video_data
	}
	
	func toGeoJson(data *[]BlockData, start_azimuth,  end_azimuth float64, StartRange int)  (map[string]interface{}){
		hold :=  []interface{}{}
		for i , block := range *data {
			if i == 0 {
				continue
			}
			hold = append(hold, map[string]interface{}{
				"type": "Feature",
				"properties": map[string]interface{}{
					"intensity": block.Intencity,
					"start_azimuth": block.StartAzimuth,
					"end_azimuth": block.EndAzimuth,
					"start_range": block.StartRange,
				},
				"geometry": map[string]interface{}{
					"type": "Point",
					"coordinates": []float64{block.Longtitude, block.Latitude},
				},
			})
		}
		fmt.Println((*data))
		return map[string]interface{}{
			"start_azimuth": start_azimuth,
			"end_azimuth": end_azimuth,
			"start_range": (*data)[0].StartRange,
			"last_point" : map[string]interface{}{
				"Longtitude": (*data)[0].Longtitude,
				"Latitude": (*data)[0].Latitude,
			},
			"features": hold,
		}
	}