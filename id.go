package bsonid

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"
)

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
	hw := md5.New()
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
func NewFromSeed(c interface{}) string {
	var b [12]byte
	// TimeStamp, 4 bytes, big endian.
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	// Machine, first 3 bytes of md5(hostname)
	for i := 0; i < len(machineID); i++ {
		b[4+i] = machineID[i]
	}

	structHash := []byte(Hash(c))

	// Pid, 2 bytes, specs don't specify endianness, but we use big endian
	b[7] = structHash[0]
	b[8] = structHash[1]

	// increment 3 bytes, big Endian
	i := atomic.AddUint32(&objectIDCounter, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)
	return hex.EncodeToString(b[:])
}

func New() string {
	var b [12]byte
	// TimeStamp, 4 bytes, big endian.
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	// Machine, first 3 bytes of md5(hostname)
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
	return hex.EncodeToString(b[:])
}
