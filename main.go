package main

import (
	// "bufio"
	"fmt"
	"os"
	"time"
	"strconv"
	"encoding/binary"
	// "reflect"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func GetHeader(f *os.File) {
	hf := new(os.File)
	*hf = *f
}

type Entry struct {
	cTime time.Time
	mTime time.Time
	Dev   uint64
	Inode uint64
	Mode  uint64
	Uid   uint64
	Gid   uint64
	Size  uint64
	Hash  string
	Name  string
}

type Index struct {
	Entries []Entry
}

// func GetUnixTime()

// Bytes2uint converts []byte to uint64
func Bytes2uint(bytes []byte) uint64 {
	padding := make([]byte, 8-len(bytes))
	i := binary.BigEndian.Uint64(append(padding, bytes...))
	return i
}

// Bytes2str converts []byte to string("00 00 00 00 00 00 00 00")
func Bytes2str(bytes []byte) string {
	strs := []string{}
	for _, b := range bytes {
		strs = append(strs, fmt.Sprintf("%02x", b))
	}
	return strings.Join(strs, " ")
}

func GetPaddingSize(had uint64) uint64 {
	Rem := had % 8
	return 8 - Rem
}

func GetUnixTime(sTime uint64) (time.Time, error) {
	unixTime := int64(sTime)
	var offsetHour, offsetMinute int
	if _, err := fmt.Sscanf("+0900", "+%02d%02d", &offsetHour, &offsetMinute); err != nil {
		return time.Time{}, err
	}
	location := time.FixedZone(" ", 3600*offsetHour+60*offsetMinute)
	timestamp := time.Unix(unixTime, 0).In(location)
	time.Now().String()
	return timestamp, nil
}

func ten2eight(mode uint64) uint64{
	tmode := int64(mode)
	emode := strconv.FormatInt(tmode, 8)
	rmode, _ := strconv.ParseInt(emode, 10, 64)
	return uint64(rmode)
}

func GetEntry(buffer []byte) (*Index, error) {
	version := Bytes2uint(buffer[4:8])
	if version != 2 {
		err := errors.New("Invalid Version Error")
		return nil, err
	}

	enum := Bytes2uint(buffer[8:12])

	buffer = buffer[12:]

	var index Index

	var count uint64
	count = 0
	for {
		if count >= enum {
			break
		}
		count++

		ctime, err := GetUnixTime(Bytes2uint(buffer[0:4]))
		if err != nil {
			fmt.Println(err)
			continue
		}
		_ = Bytes2str(buffer[4:8])

		mtime, err := GetUnixTime(Bytes2uint(buffer[8:12]))
		if err != nil {
			fmt.Println(err)
			continue
		}
		_ = Bytes2str(buffer[12:16])

		dev := Bytes2uint(buffer[16:20])
		inode := Bytes2uint(buffer[20:24])
		mode := ten2eight(Bytes2uint(buffer[24:28]))
		uid := Bytes2uint(buffer[28:32])
		gid := Bytes2uint(buffer[32:36])
		size := Bytes2uint(buffer[36:40])
		hash := Bytes2str(buffer[40:60])
		nsize := Bytes2uint(buffer[60:62])
		name := string(buffer[62 : 62+nsize])

		entry := Entry{
			cTime: ctime,
			mTime: mtime,
			Dev:   dev,
			Inode: inode,
			Mode:  mode,
			Uid:   uid,
			Gid:   gid,
			Size:  size,
			Hash:  hash,
			Name:  name,
		}

		index.Entries = append(index.Entries, entry)

		padding := GetPaddingSize(62 + nsize)
		offset := 62 + nsize + padding
		buffer = buffer[offset:]
	}

	return &index, nil
}

func LsFiles() {
	fp := "/mnt/c/Users/81701/Desktop/AtCoder/.git/index"
	f, err := os.Open(fp)
	if err != nil {
		fmt.Println("NG")
	}
	defer f.Close()

	buf := make([]byte, 64)
	buffer := make([]byte, 0)
	for {
		n, _ := (*f).Read(buf)
		if n == 0 {
			break
		}
		buffer = append(buffer, buf...)
	}

	index, err := GetEntry(buffer)
	if err != nil {
		fmt.Println(err)
	}

	for _, entry := range (*index).Entries {
		fmt.Println(entry.Mode, entry.Name)
	}
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		path := filepath.Join(dir, file.Name())
		if strings.Contains(path, ".git/") {
			continue
		}
		paths = append(paths, path)
	}

	return paths
}

func UpdateIndex() {
	paths := dirwalk("./EC-CUBE")
	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			fmt.Println(err)
		}
		fstat, err := os.Stat(absPath)
		if err != nil {
			fmt.Println(err)
		}

		var fmode os.FileMode = fstat.Mode()
	
		fmt.Println(fmode, path)

	}
}

func main() {
	LsFiles()
	// UpdateIndex()
}
