package lib

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Delta struct {
	baseObj string
	data    []byte
}

const (
	ObjCommit   int = 1
	ObjTree     int = 2
	ObjBlob     int = 3
	ObjTag      int = 4
	ObjOfsDelta int = 6
	ObjRefDelta int = 7
)

func CloneRepository(url string, directory string) {
	// Create directory
	err := os.Mkdir(directory, 0755)
	if err != nil {
		HandleError("Error creating clone directory: %s\n", err)
	}

	// Change to directory
	err = os.Chdir(directory)
	if err != nil {
		HandleError("Error changing to clone directory: %s\n", err)
	}

	// Initialize git directory
	err = InitRepository(directory)
	if err != nil {
		HandleError("Error initializing repository: %s\n", err)
	}

	// Fetch packfile
	packfile, commit, err := fetchPackfile(url)
	if err != nil {
		HandleError("Error fetching packfile: %s\n", err)
	}

	// Write packfile
	err = writePackfile(packfile)
	if err != nil {
		HandleError("Error writing packfile: %s\n", err)
	}

	// Checkout commit
	err = checkout(commit)
	if err != nil {
		HandleError("Error checking out commit: %s\n", err)
	}
}

func getPackFileResponse(url string) (*http.Response, error) {
	return http.Get(fmt.Sprintf("%s/info/refs?service=git-upload-pack", url))
}

func getUploadPackResponse(url, objName string) (*http.Response, error) {
	uploadPackRequestBuffer := bytes.NewBufferString(fmt.Sprintf("0032want %s\n00000009done\n", objName))
	return http.Post(fmt.Sprintf("%s/git-upload-pack", url), "application/x-git-upload-pack-request", uploadPackRequestBuffer)
}

func readResponse(response *http.Response) ([]byte, error) {
	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(response.Body)
	if err == nil {
		return buffer.Bytes(), nil
	} else {
		return nil, err
	}
}

func fetchPackfile(url string) ([]byte, string, error) {
	packFileResponse, err := getPackFileResponse(url)
	if err != nil {
		return nil, "", err
	}

	packBytes, err := readResponse(packFileResponse)
	if err != nil {
		return nil, "", err
	}

	packLines := readPackfile(packBytes)

	objName, err := getObjectName(packLines)
	if err != nil {
		return nil, "", err
	}

	uploadPackResponse, err := getUploadPackResponse(url, objName)
	if err != nil {
		return nil, "", err
	}

	packFile, err := readResponse(uploadPackResponse)
	if err != nil {
		return nil, "", err
	}

	lines, _, err := processPackLine(packFile)
	if err != nil {
		return nil, "", err
	}

	packFile = packFile[lines:]
	return packFile, objName, nil
}

func writePackfile(packfile []byte) error {
	err := validatePackfile(packfile)
	if err != nil {
		return err
	}

	byteIndex, objectCount, err := readPackfileHeader(packfile)
	if err != nil {
		return err
	}

	packfile, byteIndex, deltas, err := parsePackfileObjects(packfile, byteIndex, objectCount)
	if err != nil {
		return err
	}

	err = applyDeltas(deltas)

	return err
}

func readPackfileHeader(packfile []byte) (int, uint32, error) {
	byteIndex := 8
	objectCount := decodeBigUint32(packfile[byteIndex:])
	byteIndex += 4

	return byteIndex, objectCount, nil
}

func parsePackfileObjects(packfile []byte, byteIndex int, objectCount uint32) ([]byte, int, []Delta, error) {
	var deltas []Delta
	var objReadCount uint32
	packfile = packfile[:len(packfile)-20]

	for byteIndex < len(packfile) {
		var obj []byte
		objReadCount++
		objSize, objType, bRead, err := readObjectHeader(packfile[byteIndex:])
		byteIndex += bRead
		if err != nil {
			return nil, 0, nil, err
		}

		if objType == ObjCommit || objType == ObjTree || objType == ObjBlob || objType == ObjTag {
			var objTypeString string
			objTypeString, err = getObjectTypeString(objType)
			if err != nil {
				HandleError("Error reading object: invalid object type")
			}

			bRead, obj, err = readPackfileObject(packfile[byteIndex:])
			byteIndex += bRead
			if int(objSize) != len(obj) {
				HandleError("Error reading object: invalid object header size")
			}
			_, err = WriteObjectWithType(obj, objTypeString)
			if err != nil {
				HandleError("Error writing object: %s\n", err)
			}
		} else if objType == ObjOfsDelta || objType == ObjRefDelta {
			if objType == ObjOfsDelta { //for ofs delta
				_, bRead, err = decodeObjectSize(packfile[byteIndex:])
				byteIndex += bRead
				if err != nil {
					HandleError("Error decoding object size: %s\n", err)
				}
				bRead, obj, err = readPackfileObject(packfile[byteIndex:])
			} else { //for ref delta
				objHash := packfile[byteIndex : byteIndex+20]
				byteIndex += 20
				bRead, obj, err = readPackfileObject(packfile[byteIndex:])
				byteIndex += bRead
				deltas = append(deltas, Delta{baseObj: hex.EncodeToString(objHash), data: obj})
			}

			if int(objSize) != len(obj) {
				HandleError("Error reading object: invalid object header size")
			}
		} else {
			HandleError("Error reading object: invalid object type")
		}
	}

	if objReadCount != objectCount {
		HandleError("WritePackFile: Error reading object: object count mismatch") //Error handling
	}

	return packfile, byteIndex, deltas, nil
}

func applyDeltas(deltas []Delta) error {
	for len(deltas) > 0 {
		var newDeltas []Delta
		var deltaApplied bool

		for _, d := range deltas {
			if ObjectFileExists(d.baseObj) {
				deltaApplied = true
				deltaObjData, objType, _, err := ReadObjectFile(d.baseObj)
				if err != nil {
					HandleError("Error reading object: %s\n", err)
				}
				err = writeDeltaObject(deltaObjData, d.data, objType)
				if err != nil {
					HandleError("Error writing delta object: %s\n", err)
				}
			} else {
				newDeltas = append(newDeltas, d)
			}
		}
		if !deltaApplied && len(newDeltas) > 0 {
			HandleError("Error reading object: invalid delta object(s)")
		}
		deltas = newDeltas
	}
	return nil
}

func validatePackfile(packfile []byte) error {
	if len(packfile) < 32 {
		return fmt.Errorf("packfile failed validation: invalid size")
	}

	checksum := packfile[len(packfile)-20:]
	data := packfile[:len(packfile)-20]
	dataChecksum := sha1.Sum(data)

	if !bytes.Equal(checksum, dataChecksum[:]) {
		return fmt.Errorf("packfile failed validation: invalid checksum")
	}
	if !bytes.Equal(data[:4], []byte("PACK")) {
		return fmt.Errorf("packfile failed validation: invalid header")
	}
	packfileVersion := decodeBigUint32(data[4:8])
	if packfileVersion != 2 && packfileVersion != 3 {
		return fmt.Errorf("packfile failed validation: invalid version")
	}

	return nil
}

func readPackfile(packBytes []byte) [][]byte {
	var packLines [][]byte

	for len(packBytes) > 0 {
		line, data, err := processPackLine(packBytes)
		if err != nil {
			HandleError("Error processing packfile line: %s\n", err)
		}
		packBytes = packBytes[line:]
		packLines = append(packLines, data)
	}

	return packLines
}

func processPackLine(line []byte) (int, []byte, error) {
	packLength := line[:4]
	blob := line[4:]
	dest := [2]byte{}

	_, err := hex.Decode(dest[:], packLength)
	if err != nil {
		return 0, nil, err
	}

	size := uint16(dest[0])<<8 | uint16(dest[1])
	if size == 0 {
		return 4, nil, nil
	}

	if len(blob) < int(size)-4 {
		return 4, nil, fmt.Errorf("error processing packfile line")
	}

	data := blob[:size-4]
	if data[len(data)-1] == '\n' {
		data = data[:len(data)-1]
	}

	return int(size), data, nil
}

func getObjectName(packLines [][]byte) (string, error) {
	for _, line := range packLines[1:] {
		if len(line) == 0 {
			continue
		}
		var objHash, objRef string
		fmt.Sscanf(string(line), "%s %s", &objHash, &objRef)
		if objRef == "refs/heads/master" {
			return objHash, nil
		}
	}
	return "", fmt.Errorf("name not found")
}

func decodeBigUint32(bytes []byte) uint32 {
	return uint32(bytes[0])<<24 | uint32(bytes[1])<<16 | uint32(bytes[2])<<8 | uint32(bytes[3])
}

func getObjectTypeString(objType int) (string, error) {
	m := map[int]string{
		ObjCommit:   "commit",
		ObjTree:     "tree",
		ObjBlob:     "blob",
		ObjTag:      "tag",
		ObjOfsDelta: "ofs-delta",
		ObjRefDelta: "ref-delta",
	}

	result, ok := m[objType]
	if !ok {
		return "", fmt.Errorf("unknown object type: %d", objType)
	}
	return result, nil
}

func decodeObjectSize(objData []byte) (int, int, error) {
	byteIndex := bytes.IndexByte(objData, 0)
	var objSize int

	fmt.Sscanf(string(objData[:byteIndex]), "%d", &objSize)

	if byteIndex+objSize+1 != len(objData) {
		return 0, 0, fmt.Errorf("invalid object size")
	}

	return objSize, byteIndex, nil
}

func checkout(commitHash string) error {
	commit, objType, _, err := ReadObjectFile(commitHash)
	if err != nil {
		return err
	}

	if objType != "commit" {
		return fmt.Errorf("error reading commit: invalid object type")
	}

	treeHash := commit[5:45]

	err = checkoutTree(string(treeHash), ".")
	if err != nil {
		return fmt.Errorf("error checking out tree: %s\n", err)
	}

	return nil
}

func checkoutTree(treeHash string, path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	tree, err := ReadTreeObjectFile(treeHash)
	if err != nil {
		return err
	}

	for _, entry := range tree {
		entryHash := hex.EncodeToString(entry.hash[:])
		objPath := fmt.Sprintf("%s/%s", path, entry.name)
		if entry.mode == ModeTree {
			err = checkoutTree(entryHash, objPath)
			if err != nil {
				return err
			}
		} else if entry.mode == ModeBlob || entry.mode == ModeBlobExec {
			obj, objType, _, err := ReadObjectFile(entryHash)
			if err != nil {
				return err
			}

			if objType != "blob" {
				return fmt.Errorf("error reading object: object is not a blob")
			}

			err = WriteFile(objPath, obj)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func readObjectHeader(packfile []byte) (size uint64, objectType int, byteIndex int, err error) {
	data := packfile[byteIndex]
	byteIndex++
	objectType = int((data >> 4) & 0x7)
	size = uint64(data & 0xF)
	shift := 4
	for data&0x80 != 0 {
		if len(packfile) <= byteIndex || 64 <= shift {
			return 0, int(0), 0, errors.New("bad object header")
		}
		data = packfile[byteIndex]
		byteIndex++
		size += uint64(data&0x7F) << shift
		shift += 7
	}
	return size, objectType, byteIndex, nil
}

func readPackfileObject(packfile []byte) (int, []byte, error) {
	b := bytes.NewReader(packfile)
	r, err := zlib.NewReader(b)
	if err != nil {
		return 0, nil, err
	}
	defer r.Close()
	object, err := io.ReadAll(r)
	if err != nil {
		return 0, nil, err
	}
	bytesRead := int(b.Size()) - b.Len()
	return bytesRead, object, nil
}

func writeDeltaObject(baseObject, deltaObject []byte, objectType string) error {
	used := 0
	baseSize, read, err := readSize(deltaObject[used:])
	if err != nil {
		return err
	}
	used += read
	if len(baseObject) != int(baseSize) {
		return errors.New("bad delta header")
	}
	expectedSize, read, err := readSize(deltaObject[used:])
	if err != nil {
		return err
	}
	used += read
	buffer := bytes.Buffer{}
	for used < len(deltaObject) {
		opcode := deltaObject[used]
		used++
		if opcode&0x80 != 0 {
			var argument uint64
			for bit := 0; bit < 7; bit++ {
				if opcode&(1<<bit) != 0 {
					argument += uint64(deltaObject[used]) << (bit * 8)
					used++
				}
			}
			offset := argument & 0xFFFFFFFF
			size := (argument >> 32) & 0xFFFFFF
			if size == 0 {
				size = 0x10000
			}
			buffer.Write(baseObject[offset : offset+size])
		} else {
			size := int(opcode & 0x7F)
			buffer.Write(deltaObject[used : used+size])
			used += size
		}
	}
	objToDelta := buffer.Bytes()
	if int(expectedSize) != len(objToDelta) {
		return errors.New("bad delta header")
	}
	_, err = WriteObjectWithType(objToDelta, objectType)
	if err != nil {
		return err
	}
	return nil
}

func readSize(packfile []byte) (size uint64, used int, err error) {
	data := packfile[used]
	used++
	size = uint64(data & 0x7F)
	shift := 7
	for data&0x80 != 0 {
		if len(packfile) <= used || 64 <= shift {
			return 0, 0, errors.New("bad size")
		}
		data = packfile[used]
		used++
		size += uint64(data&0x7F) << shift
		shift += 7
	}
	return size, used, nil
}
