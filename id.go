package id

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"
)

type Id struct {
	stringValue    string
	byteSliceValue []byte
	uint32Value    uint32
}

var objectIDCounter = getRandomCounter()

func getRandomCounter() uint32 {
	counter := make([]byte, 4)
	rand.Read(counter)
	return binary.LittleEndian.Uint32(counter)
}

// TODO: Change the pid to a hash of the object data.
var pid = os.Getpid()
var machineID = generateMachineID()

func generateMachineID() [3]byte {
	var sum [3]byte // 3 byte Machine ID
	id := sum[:]
	hostname, err1 := os.Hostname()
	if err1 != nil {
		// if getting hostname failed
		// get a crypto random id and return
		_, err2 := io.ReadFull(rand.Reader, id)
		if err2 != nil {
			panic(fmt.Errorf("Cannot get hostname: %v, %v", err1, err2))
		}
		copy(sum[:], id)
		return sum
	}
	// TODO Why md5? Lets use xxhash!
	hw := NewXXHash32()
	// append hostname to the running hash
	hw.Write([]byte(hostname))
	copy(sum[:], hw.Sum(nil))
	return sum
}

// New returns a new unique ObjectId.
// 4 byte time,
// 3 byte Machine ID
// 2 byte pid
// 3 byte self increased id.
func NewFromSeed(c interface{}) Id {
	var b [12]byte
	// TimeStamp, 4 bytes, big endian.
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	// Machine, first 3 bytes of md5(hostname)
	for i := 0; i < len(machineID); i++ {
		b[4+i] = machineID[i]
	}

	structHash := Hash(c).Bytes()

	// Pid, 2 bytes, specs don't specify endianness, but we use big endian
	b[7] = structHash[0]
	b[8] = structHash[1]

	// increment 3 bytes, big Endian
	i := atomic.AddUint32(&objectIDCounter, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)

	hw := NewXXHash32()
	hw.Write(b[:])
	return Id{
		stringValue:    hex.EncodeToString(b[:]),
		uint32Value:    hw.Sum32(),
		byteSliceValue: b[:],
	}
}

func New() Id {
	var b [12]byte
	// TimeStamp, 4 bytes, big endian.
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	// Machine, first 3 bytes of xxh32(hostname)
	for i := 0; i < len(machineID); i++ {
		b[4+i] = machineID[i]
	}

	// Pid, 2 bytes, specs don't specify endianness, but we use big endian
	b[7] = byte(pid >> 8)
	b[8] = byte(pid)

	// increment 3 bytes, big Endian
	i := atomic.AddUint32(&objectIDCounter, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)

	hw := NewXXHash32()
	hw.Write(b[:])
	return Id{
		stringValue:    hex.EncodeToString(b[:]),
		uint32Value:    hw.Sum32(),
		byteSliceValue: b[:],
	}
}

func NewShort() Id {
	// TODO: Take out the pid , then nanoId could be take out the machine, then
	// pico could be take out both
	var b [8]byte
	// TimeStamp, 4 bytes, big endian.
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	// Machine, first 3 bytes of xxh32(hostname)
	for i := 2; i < len(machineID); i++ {
		b[2+i] = machineID[i]
	}

	// increment 3 bytes, big Endian
	i := atomic.AddUint32(&objectIDCounter, 1)
	b[5] = byte(i >> 16)
	b[6] = byte(i >> 8)
	b[7] = byte(i)

	hw := NewXXHash32()
	hw.Write(b[:])
	return Id{
		stringValue:    hex.EncodeToString(b[:]),
		uint32Value:    hw.Sum32(),
		byteSliceValue: b[:],
	}
}

// TODO: Would be nice to be able to extract the timestamp, and possibly the
// value used as a seed.
func (self Id) Bytes() []byte  { return self.byteSliceValue }
func (self Id) UInt32() uint32 { return self.uint32Value }
func (self Id) String() string { return self.stringValue }
