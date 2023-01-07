package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

type RIFF struct {
	ChunkID     [4]byte
	ChunkSize   [4]byte
	ChunkFormat [4]byte
}

func (r RIFF) print() {
	fmt.Println("ChunkId: ", string(r.ChunkID[:]))
	fmt.Println("ChunkSize: ", binary.LittleEndian.Uint32(r.ChunkSize[:])+8)
	fmt.Println("ChunkFormat: ", string(r.ChunkFormat[:]))
	fmt.Println()
}

type FMT struct {
	SubChunk1ID   [4]byte
	SubChunk1Size [4]byte
	AudioFormat   [2]byte
	NumChannels   [2]byte
	SampleRate    [4]byte
	ByteRate      [4]byte
	BlockAlign    [2]byte
	BitsPerSample [2]byte
}

func (fm FMT) print() {
	fmt.Println("SubChunk1Id: ", string(fm.SubChunk1ID[:]))
	fmt.Println("SubChunk1Size: ", binary.LittleEndian.Uint32(fm.SubChunk1Size[:]))
	fmt.Println("AudioFormat: ", binary.LittleEndian.Uint16(fm.AudioFormat[:]))
	fmt.Println("NumChannels: ", binary.LittleEndian.Uint16(fm.NumChannels[:]))
	fmt.Println("SampleRate: ", binary.LittleEndian.Uint32(fm.SampleRate[:]))
	fmt.Println("ByteRate: ", binary.LittleEndian.Uint32(fm.ByteRate[:]))
	fmt.Println("BlockAlign: ", binary.LittleEndian.Uint16(fm.BlockAlign[:]))
	fmt.Println("BitsPerSample: ", binary.LittleEndian.Uint16(fm.BitsPerSample[:]))
	fmt.Println()
}

type LIST struct {
	ChunkID  [4]byte
	size     [4]byte
	listType [4]byte
	data     []byte
}

func (list LIST) print() {
	fmt.Println("ChunkId: ", string(list.ChunkID[:]))
	fmt.Println("size: ", binary.LittleEndian.Uint32(list.size[:]))
	fmt.Println("listType: ", string(list.listType[:]))
	fmt.Println("data: ", string(list.data))
	fmt.Println()
}

func (list *LIST) read(file *os.File) {

	listCondition := make([]byte, 4)
	file.Read(listCondition)
	file.Seek(-4, 1)

	if string(listCondition) != "LIST" {
		return
	}

	binary.Read(file, binary.BigEndian, &list.ChunkID)
	binary.Read(file, binary.BigEndian, &list.size)
	binary.Read(file, binary.BigEndian, &list.listType)
	list.data = make([]byte, binary.LittleEndian.Uint32(list.size[:])-4)
	binary.Read(file, binary.BigEndian, &list.data)
}

type DATA struct {
	SubChunk2Id   [4]byte
	SubChunk2Size [4]byte
	data          []byte
}

func (data DATA) print() {
	fmt.Println("SubChunk2Id: ", string(data.SubChunk2Id[:]))
	fmt.Println("SubChunk2Size: ", binary.BigEndian.Uint32(data.SubChunk2Size[:]))
	fmt.Println("first 100 samples", data.data[:100])
	fmt.Println()
}

func (data *DATA) read(file *os.File) {
	binary.Read(file, binary.BigEndian, &data.SubChunk2Id)
	binary.Read(file, binary.BigEndian, &data.SubChunk2Size)
	data.data = make([]byte, binary.LittleEndian.Uint32(data.SubChunk2Size[:]))
	binary.Read(file, binary.BigEndian, &data.data)
}

func main() {
	file, err := os.Open("thank.wav")

	if err != nil {
		panic(err)
	}

	// RIFF Chunk
	RIFFChunk := RIFF{}
	binary.Read(file, binary.BigEndian, &RIFFChunk)
	RIFFChunk.print()

	FMTChunk := FMT{}
	binary.Read(file, binary.BigEndian, &FMTChunk)
	FMTChunk.print()

	//LIST Chunk
	listChunk := LIST{}
	listChunk.read(file)
	listChunk.print()

	//DATA Chunk
	dataChunk := DATA{}
	dataChunk.read(file)

	dataChunk.print()

}
