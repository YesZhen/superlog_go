package superlog

import (
	"fmt"
	"os"
	"unsafe"
	"syscall"
	"sync/atomic"
	"errors"
	"time"
	"bytes"
	"runtime"
	"strconv"
	//https://github.com/overtalk/shm
	//"github.com/overtalk/shm"
	"github.com/silentred/gid"
)

const (
	SUPERLOG_OVERRIDE = 2
	MSG_N		= 1*1024*1024
	CREATE_MODE	= "0644"
	KEEPEXIST	= 1
	MAX_DATA_LEN = 32
	MSG_SIZE = 64
	SUPERLOG_DEFAULT_FILE_PATH = "/dev/shm/superlog"
	SUPERLOG_DEFAULT_INIT_FLAGS = 0
	SUPERLOG_RING_SIZE = 8 * 4 + 32 + MSG_N * MSG_SIZE
)

type Msg struct {
	Timestamp	int64
	Pid	int64
	Tid	int64
	Name [8]byte
	Data [MAX_DATA_LEN]byte
}

type Ring struct {
	head uint64
	tail uint64
	max_len uint64
	msg_sz uint64
	//pad for dummy space
	pad [32]byte
	msgs [MSG_N]Msg
}

var Name [8]byte

func init() {
	
}

//
func Min(x, y int64) int64 {
    if x < y {
        return x
    }
    return y
}

//return go routine ID
func getGID() uint64 {
    b := make([]byte, 64)
    b = b[:runtime.Stack(b, false)]
    b = bytes.TrimPrefix(b, []byte("goroutine "))
    b = b[:bytes.IndexByte(b, ' ')]
    n, _ := strconv.ParseUint(string(b), 10, 64)
    return n
}

//Return a Msg struct array
func CreateMsg() ([]Msg) {
	return []Msg{}
}

//Append a Msg struct to Msg struct array
func AppendMsg(msgs []Msg, data string) ([]Msg, error) {
	if len(data) > MAX_DATA_LEN{
		//fmt.Println("data lenght exceed the limit(56 byte)")
        return msgs, errors.New("AppendMsg error, data length exceeds the limit(56 byte)")
	}

	item := Msg{}
	item.Timestamp = time.Now().UnixNano() / int64(time.Microsecond)
	item.Pid = int64(os.Getpid())
	item.Tid = int64(getGID())
	copy(item.Data[:], data)
	msgs = append(msgs,item)

	return msgs, nil
}



func SetName(name string) {

	// nameMap[r] = []byte(name)
	copy(Name[:], name)
}

func InitDefault() (*Ring, error) {
	return Init(SUPERLOG_DEFAULT_FILE_PATH, SUPERLOG_DEFAULT_INIT_FLAGS)
}

//Open shared memory and mmap into virtual memory
func Init(path string, flags uint) (*Ring, error) {

	_, err := os.Stat(path)
	exists := !os.IsNotExist(err)


	if exists && (flags & SUPERLOG_OVERRIDE) != 0 {
		if err = os.Remove(path); err != nil {
			return nil, err
		}
	}

	
	fd, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE, os.ModePerm)
    if err != nil {
		return nil, err
    }

	syscall.Fchmod(int(fd.Fd()), 0777)

	defer fd.Close()

	// Info, err := fd.Stat()
    // if err != nil {
	// 	return nil, err
    // }


	if err = syscall.Ftruncate(int(fd.Fd()), SUPERLOG_RING_SIZE); err != nil {
		return nil, err
	}

	buf, err := syscall.Mmap(int(fd.Fd()), 0, SUPERLOG_RING_SIZE, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	//make buf points to *Ring
	r := (*Ring)(unsafe.Pointer(&buf[0]))

	if !exists || (flags & SUPERLOG_OVERRIDE) != 0 {
		r.head = 0
		r.tail = 0
		r.max_len = MSG_N
		r.msg_sz = MSG_SIZE
	}

	return r, nil
	
}

//Write Msg struct array to memory without lock
func Write(r *Ring, msgs []Msg, n uint64) {
	//alloc slots, __sync_fetch_and_add
	start := atomic.AddUint64(&r.tail, n) - n
	start &= r.max_len - 1
	
	if start + n <= r.max_len {
		copy(r.msgs[start:], msgs)
	}else {
		copy(r.msgs[start:], msgs[0:(r.max_len-start)])
		copy(r.msgs[0:], msgs[r.max_len-start:])
	}

}

// Write a single log
func Log(r *Ring, data string) Msg {

	start := atomic.AddUint64(&r.tail, 1) - 1
	start &= r.max_len - 1

	r.msgs[start].Name = Name
	r.msgs[start].Timestamp = time.Now().UnixNano() / int64(time.Microsecond)
	r.msgs[start].Pid = int64(os.Getpid())
	r.msgs[start].Tid = gid.Get()
	copy(r.msgs[start].Data[:], data)
	return r.msgs[start]
}


func Read(r *Ring, n uint64) ([]Msg){
	
	n = uint64(Min(int64((r.tail + r.max_len - r.head) & (r.max_len - 1)), int64(n)))
	msgs := make([]Msg, n);
	//alloc slots, __sync_fetch_and_add
	start := atomic.AddUint64(&r.head, n) - n
	start &= r.max_len - 1
	
	if start + n <= r.max_len {
		copy(msgs, r.msgs[start:(start + n)])
	}else {
		copy(msgs, r.msgs[start:])
		copy(msgs[(r.max_len-start):], r.msgs[0:(n+start-r.max_len)])
	}	
	return msgs
}

func PrintRing(r *Ring) {
	fmt.Println("r.head = ", r.head)
	fmt.Println("r.tail = ", r.tail)
	fmt.Println("r.max_len = ", r.max_len)
	fmt.Println("r.msg_sz = ", r.msg_sz)
}
