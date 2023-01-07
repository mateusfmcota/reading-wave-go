package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

type RIFF struct {
	ChunkID     []byte
	ChunkSize   []byte
	ChunkFormat []byte
}

type FMT struct {
	SubChunk1ID   []byte
	SubChunk1Size []byte
	AudioFormat   []byte
	NumChannels   []byte
	SampleRate    []byte
	ByteRate      []byte
	BlockAlign    []byte
	BitsPerSample []byte
}

type LIST struct {
	ChunkID  []byte
	size     []byte
	listType []byte
	data     []byte
}

type DATA struct {
	SubChunk2Id   []byte
	SubChunk2Size []byte
	data          []byte
}

func readNBytes(file *os.File, n int) []byte {
	temp := make([]byte, n)

	_, err := file.Read(temp)
	if err != nil {
		panic(err)
	}

	return temp
}

func printRiff(rf RIFF) {
	fmt.Println("ChunkId: ", string(rf.ChunkID))
	fmt.Println("ChunkSize: ", binary.LittleEndian.Uint32(rf.ChunkSize)+8)
	fmt.Println("ChunkFormat: ", string(rf.ChunkFormat))

}

func printFMT(fm FMT) {
	fmt.Println("SubChunk1Id: ", string(fm.SubChunk1ID))
	fmt.Println("SubChunk1Size: ", binary.LittleEndian.Uint32(fm.SubChunk1Size))
	fmt.Println("AudioFormat: ", binary.LittleEndian.Uint16(fm.AudioFormat))
	fmt.Println("NumChannels: ", binary.LittleEndian.Uint16(fm.NumChannels))
	fmt.Println("SampleRate: ", binary.LittleEndian.Uint32(fm.SampleRate))
	fmt.Println("ByteRate: ", binary.LittleEndian.Uint32(fm.ByteRate))
	fmt.Println("BlockAlign: ", binary.LittleEndian.Uint16(fm.BlockAlign))
	fmt.Println("BitsPerSample: ", binary.LittleEndian.Uint16(fm.BitsPerSample))
}

func printLIST(list LIST) {
	fmt.Println("ChunkId: ", string(list.ChunkID))
	fmt.Println("size: ", binary.LittleEndian.Uint32(list.size))
	fmt.Println("listType: ", string(list.listType))
	fmt.Println("data: ", string(list.data))
}

func printData(data DATA) {
	fmt.Println("SubChunk2Id: ", string(data.SubChunk2Id))
	fmt.Println("SubChunk2Size: ", binary.LittleEndian.Uint32(data.SubChunk2Size))
	fmt.Println("data", data.data)
}

func main() {
	file, err := os.Open("thank.wav")

	if err != nil {
		panic(err)
	}

	// RIFF Chunk
	RIFFChunk := RIFF{}

	RIFFChunk.ChunkID = readNBytes(file, 4)
	RIFFChunk.ChunkSize = readNBytes(file, 4)
	RIFFChunk.ChunkFormat = readNBytes(file, 4)

	// FMT sub-chunk
	FMTChunk := FMT{}

	FMTChunk.SubChunk1ID = readNBytes(file, 4)
	FMTChunk.SubChunk1Size = readNBytes(file, 4)
	FMTChunk.AudioFormat = readNBytes(file, 2)
	FMTChunk.NumChannels = readNBytes(file, 2)
	FMTChunk.SampleRate = readNBytes(file, 4)
	FMTChunk.ByteRate = readNBytes(file, 4)
	FMTChunk.BlockAlign = readNBytes(file, 2)
	FMTChunk.BitsPerSample = readNBytes(file, 2)

	// https://www.recordingblogs.com/wiki/list-chunk-of-a-wave-file
	subChunk := readNBytes(file, 4)
	var listChunk *LIST

	if string(subChunk) == "LIST" {
		listChunk = new(LIST)
		listChunk.ChunkID = subChunk
		listChunk.size = readNBytes(file, 4)
		listChunk.listType = readNBytes(file, 4)
		listChunk.data = readNBytes(file, int(binary.LittleEndian.Uint32(listChunk.size))-4)
	}

	// Data sub-chunk
	data := DATA{}

	data.SubChunk2Id = readNBytes(file, 4)
	data.SubChunk2Size = readNBytes(file, 4)
	data.data = readNBytes(file, int(binary.LittleEndian.Uint32(data.SubChunk2Size)))

	printRiff(RIFFChunk)
	fmt.Println("")
	printFMT(FMTChunk)
	fmt.Println("")
	if listChunk != nil {
		printLIST(*listChunk)
	}
	fmt.Println("")
	printData(data)

}
