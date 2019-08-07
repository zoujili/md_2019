//Copyright (c) 2016, Bruce
//All rights reserved.

// Package code provides a very simple Twitter code generator and parser.

/* updated by jili zou .
* add opt for default generate strategy
* from  https://github.com/baidu/uid-generator.  default start epoch "2018.08.01"
* <pre>{@code
* +------+----------------------+----------------+-----------+
* | sign |     delta seconds    | worker node id | sequence  |
* +------+----------------------+----------------+-----------+
*   1bit          28bits              22bits         13bits
* }</pre>

* sign(1bit)
  The highest bit is always 0.

* delta seconds (28 bits)
  The next 28 bits, represents delta seconds since a customer epoch(2016-05-20). The maximum time will be 8.7 years.

* worker id (22 bits)
  The next 22 bits, represents the worker node id, maximum value will be 4.2 million. UidGenerator uses a build-in
  database based ```worker id assigner``` when startup by default, and it will dispose previous work node id after
  reboot. Other strategy such like 'reuse' is coming soon.

* sequence (13 bits)
  the last 13 bits, represents sequence within the one second, maximum is 8192 per second by default.



* form  https://github.com/bwmarrin/snowflake.   default start epoch "2012.01.01"
* <pre>{@code
* +------+----------------------+----------------+-----------+
* | sign |     delta seconds    | worker node id | sequence  |
* +------+----------------------+----------------+-----------+
*   1bit          41bits              10bits         12bits
* }</pre>

* The ID as a whole is a 63 bit integer stored in an int64
* 41 bits are used to store a timestamp with millisecond precision, using a custom epoch.
* 10 bits are used to store a node id - a range from 0 through 1023.
* 12 bits are used to store a sequence number - a range from 0 through 4095.

*/

package code

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	Epoch int64

	// NodeBits holds the number of bits to use for Node
	NodeBits uint8

	// StepBits holds the number of bits to use for Step
	StepBits uint8

	// DEPRECATED: the below four variables will be removed in a future release.
	mu        sync.Mutex
	nodeMax   int64 = -1 ^ (-1 << NodeBits)
	nodeMask        = nodeMax << StepBits
	stepMask  int64 = -1 ^ (-1 << StepBits)
	timeShift       = NodeBits + StepBits
	nodeShift       = StepBits
)

const encodeBase32Map = "ybndrfg8ejkmcpqxot1uwisza345h769"

var decodeBase32Map [256]byte

const encodeBase58Map = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

var decodeBase58Map [256]byte

// A JSONSyntaxError is returned from UnmarshalJSON if an invalid ID is provided.
type JSONSyntaxError struct{ original []byte }

func (j JSONSyntaxError) Error() string {
	return fmt.Sprintf("invalid code ID %q", string(j.original))
}

// ErrInvalidBase58 is returned by ParseBase58 when given an invalid []byte
var ErrInvalidBase58 = errors.New("invalid base58")

// ErrInvalidBase32 is returned by ParseBase32 when given an invalid []byte
var ErrInvalidBase32 = errors.New("invalid base32")

// Create maps for decoding Base58/Base32.
// This speeds up the process tremendously.
func init() {

	for i := 0; i < len(encodeBase58Map); i++ {
		decodeBase58Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase58Map); i++ {
		decodeBase58Map[encodeBase58Map[i]] = byte(i)
	}

	for i := 0; i < len(encodeBase32Map); i++ {
		decodeBase32Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase32Map); i++ {
		decodeBase32Map[encodeBase32Map[i]] = byte(i)
	}
}

// A Node struct holds the basic information needed for a code generator
// node
type Node struct {
	mu    sync.Mutex
	epoch time.Time
	time  int64
	node  int64
	step  int64

	nodeMax   int64
	nodeMask  int64
	stepMask  int64
	timeShift uint8
	nodeShift uint8

	opts options
}

// An ID is a custom type used for a code ID.  This is used so we can
// attach methods onto the ID.
type ID int64

// NewNode returns a new code node that can be used to generate code
// IDs
func NewNode(node int64, opts ...Option) (*Node, error) {
	// re-calc in case custom NodeBits or StepBits were set
	// DEPRECATED: the below block will be removed in a future release.

	mu.Lock()
	nodeMax = -1 ^ (-1 << NodeBits)
	nodeMask = nodeMax << StepBits
	stepMask = -1 ^ (-1 << StepBits)
	timeShift = NodeBits + StepBits
	nodeShift = StepBits
	mu.Unlock()

	n := Node{}

	if len(opts) > 0 {
		n.Option(opts...)
	} else {
		//Default
		n.Option(MySQLAssignerOptions...)
	}

	Epoch = n.opts.Epoch
	//fmt.Println(Epoch)
	NodeBits = n.opts.NodeBits
	StepBits = n.opts.StepBits

	n.node = node
	n.nodeMax = -1 ^ (-1 << NodeBits)
	n.nodeMask = n.nodeMax << StepBits
	n.stepMask = -1 ^ (-1 << StepBits)
	n.timeShift = NodeBits + StepBits
	n.nodeShift = StepBits

	if n.node < 0 || n.node > n.nodeMax {
		return nil, errors.New("Node number must be between 0 and " + strconv.FormatInt(n.nodeMax, 10))
	}

	var curTime = time.Now()
	// add time.Duration to curTime to make sure we use the monotonic clock if available
	n.epoch = curTime.Add(time.Unix(Epoch, 0).Sub(curTime))

	return &n, nil
}

func (n *Node) Option(opts ...Option) *Node {
	for _, opt := range opts {
		opt(&n.opts)
	}
	return n
}

// Generate creates and returns a unique code ID
// To help guarantee uniqueness
// - Make sure your system is keeping accurate system time
// - Make sure you never have multiple nodes running with the same node ID
func (n *Node) Generate() ID {
	n.mu.Lock()
	var r ID
	if n.opts.IsMySQLAssigner {
		now := time.Now().Unix()
		if now == n.time {
			n.step = (n.step + 1) & n.stepMask

			if n.step == 0 {
				for now <= n.time {
					now = time.Now().Unix()
				}
			}
		} else {
			n.step = 0
		}

		n.time = now

		r = ID((now-Epoch)<<n.timeShift |
			(n.node << n.nodeShift) |
			(n.step),
		)

	} else {
		now := time.Since(n.epoch).Nanoseconds() / 1000000
		if now == n.time {
			n.step = (n.step + 1) & n.stepMask

			if n.step == 0 {
				for now <= n.time {
					now = time.Since(n.epoch).Nanoseconds() / 1000000
				}
			}
		} else {
			n.step = 0
		}

		n.time = now

		r = ID((now)<<n.timeShift |
			(n.node << n.nodeShift) |
			(n.step),
		)
	}

	n.mu.Unlock()
	return r
}

// Int64 returns an int64 of the code ID
func (f ID) Int64() int64 {
	return int64(f)
}

// ParseInt64 converts an int64 into a code ID
func ParseInt64(id int64) ID {
	return ID(id)
}

// String returns a string of the code ID
func (f ID) String() string {
	return strconv.FormatInt(int64(f), 10)
}

// ParseString converts a string into a code ID
func ParseString(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	return ID(i), err

}

// Base2 returns a string base2 of the code ID
func (f ID) Base2() string {
	return strconv.FormatInt(int64(f), 2)
}

// ParseBase2 converts a Base2 string into a code ID
func ParseBase2(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 2, 64)
	return ID(i), err
}

// Base32 uses the z-base-32 character set but encodes and decodes similar
// to base58, allowing it to create an even smaller result string.
// NOTE: There are many different base32 implementations so becareful when
// doing any interoperation.
func (f ID) Base32() string {

	if f < 32 {
		return string(encodeBase32Map[f])
	}

	b := make([]byte, 0, 12)
	for f >= 32 {
		b = append(b, encodeBase32Map[f%32])
		f /= 32
	}
	b = append(b, encodeBase32Map[f])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

// ParseBase32 parses a base32 []byte into a code ID
// NOTE: There are many different base32 implementations so becareful when
// doing any interoperation.
func ParseBase32(b []byte) (ID, error) {

	var id int64

	for i := range b {
		if decodeBase32Map[b[i]] == 0xFF {
			return -1, ErrInvalidBase32
		}
		id = id*32 + int64(decodeBase32Map[b[i]])
	}

	return ID(id), nil
}

// Base36 returns a base36 string of the code ID
func (f ID) Base36() string {
	return strconv.FormatInt(int64(f), 36)
}

// ParseBase36 converts a Base36 string into a code ID
func ParseBase36(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 36, 64)
	return ID(i), err
}

// Base58 returns a base58 string of the code ID
func (f ID) Base58() string {

	if f < 58 {
		return string(encodeBase58Map[f])
	}

	b := make([]byte, 0, 11)
	for f >= 58 {
		b = append(b, encodeBase58Map[f%58])
		f /= 58
	}
	b = append(b, encodeBase58Map[f])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

// ParseBase58 parses a base58 []byte into a code ID
func ParseBase58(b []byte) (ID, error) {

	var id int64

	for i := range b {
		if decodeBase58Map[b[i]] == 0xFF {
			return -1, ErrInvalidBase58
		}
		id = id*58 + int64(decodeBase58Map[b[i]])
	}

	return ID(id), nil
}

// Base64 returns a base64 string of the code ID
func (f ID) Base64() string {
	return base64.StdEncoding.EncodeToString(f.Bytes())
}

// ParseBase64 converts a base64 string into a code ID
func ParseBase64(id string) (ID, error) {
	b, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return -1, err
	}
	return ParseBytes(b)

}

// Bytes returns a byte slice of the code ID
func (f ID) Bytes() []byte {
	return []byte(f.String())
}

// ParseBytes converts a byte slice into a code ID
func ParseBytes(id []byte) (ID, error) {
	i, err := strconv.ParseInt(string(id), 10, 64)
	return ID(i), err
}

// IntBytes returns an array of bytes of the code ID, encoded as a
// big endian integer.
func (f ID) IntBytes() [8]byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(f))
	return b
}

// ParseIntBytes converts an array of bytes encoded as big endian integer as
// a code ID
func ParseIntBytes(id [8]byte) ID {
	return ID(int64(binary.BigEndian.Uint64(id[:])))
}

// Time returns an int64 unix timestamp in milliseconds of the code ID time
// DEPRECATED: the below function will be removed in a future release.
func (f ID) Time() int64 {
	return (int64(f) >> timeShift) + Epoch
}

// Node returns an int64 of the code ID node number
// DEPRECATED: the below function will be removed in a future release.
func (f ID) Node() int64 {
	return int64(f) & nodeMask >> nodeShift
}

// Step returns an int64 of the code step (or sequence) number
// DEPRECATED: the below function will be removed in a future release.
func (f ID) Step() int64 {
	return int64(f) & stepMask
}

// MarshalJSON returns a json byte array string of the code ID.
func (f ID) MarshalJSON() ([]byte, error) {
	buff := make([]byte, 0, 22)
	buff = append(buff, '"')
	buff = strconv.AppendInt(buff, int64(f), 10)
	buff = append(buff, '"')
	return buff, nil
}

// UnmarshalJSON converts a json byte array of a code ID into an ID type.
func (f *ID) UnmarshalJSON(b []byte) error {
	if len(b) < 3 || b[0] != '"' || b[len(b)-1] != '"' {
		return JSONSyntaxError{b}
	}

	i, err := strconv.ParseInt(string(b[1:len(b)-1]), 10, 64)
	if err != nil {
		return err
	}

	*f = ID(i)
	return nil
}
