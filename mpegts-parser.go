package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

const syncByte = byte(0x47)
const pidMask = byte(0x1F)
const packetSize = 188

func main() {
	byteArray, err := io.ReadAll(os.Stdin)

	if err != nil {
		log.Fatal("Error reading from stdin: ", err)
	}

	indexArray := findSyncByteIndices(byteArray, syncByte)
	lastSyncByteIndex, err := getLastSyncByteIndex(indexArray, byteArray)

	if err != nil {
		log.Fatalf("Could not get last sync byte index: %v", err)
	}

	startIndex, err := findStartIndex(indexArray, lastSyncByteIndex, packetSize)

	if err != nil {
		log.Fatalf("Could not find start packet index to start parsing: %v", err)
	}

	uniquePIDS, err := processPackets(startIndex, lastSyncByteIndex, byteArray)

	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(uniquePIDS, func(i, j int) bool {
		return uniquePIDS[i] < uniquePIDS[j]
	})

	for _, pid := range uniquePIDS {
		fmt.Printf("0x%x\n", pid)
	}
}

func findSyncByteIndices(byteArray []byte, syncByte byte) []int {
	indexArray := []int{}
	for index, value := range byteArray {
		if value == syncByte {
			indexArray = append(indexArray, index)
		}
	}
	return indexArray
}

func getLastSyncByteIndex(syncByteIndices []int, byteArray []byte) (int, error) {
	if len(syncByteIndices) < 2 {
		return 0, fmt.Errorf("not enough sync byte indices to proceed")
	}
	return syncByteIndices[len(syncByteIndices)-1], nil
}

func findStartIndex(indexArray []int, lastSyncByteIndex int, packetSize int) (int, error) {
	for _, indexValue := range indexArray {
		check := (lastSyncByteIndex - indexValue) % packetSize

		if check == 0 && indexValue != 0 {
			return indexValue, nil
		}
	}
	return 0, fmt.Errorf("could not find a valid start index")
}

func processPackets(startIndex int, lastSyncByteIndex int, byteArray []byte) ([]uint16, error) {
	packetNumber := 1
	pidMap := make(map[uint16]bool)
	uniquePids := []uint16{}

	for i := startIndex; i <= lastSyncByteIndex; i += packetSize {
		if i+packetSize > len(byteArray) {
			continue // skip packet
		}

		packetBytes := byteArray[i : i+packetSize]

		pid, err := validatePacket(packetBytes, packetNumber, i)

		if err != nil {
			return nil, err
		}

		recordUniquePid(pid, pidMap, &uniquePids)

		packetNumber += 1
	}

	return uniquePids, nil
}

func validatePacket(packetBytes []byte, packetNumber int, bytesIndex int) (uint16, error) {
	if packetBytes[0] != syncByte {
		return 0, fmt.Errorf("error: No sync byte present in packet %v, offset %v", packetNumber, bytesIndex)
	}

	last5Bits := packetBytes[1] & pidMask

	shiftedBits := uint16(last5Bits) << 8
	packetPid := shiftedBits | uint16(packetBytes[2])

	return packetPid, nil
}

func recordUniquePid(pid uint16, pidMap map[uint16]bool, uniquePIDS *[]uint16) {
	if pidMap[pid] {
		return // not unique
	}
	*uniquePIDS = append(*uniquePIDS, pid)
	pidMap[pid] = true
}
