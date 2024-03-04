package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type FileParser struct {
    index int64 
    file []byte
}

func (fp FileParser) jumpByte() FileParser {
    fp.index += 1
    return fp
}

func (fp *FileParser) jumpBytes(count int64) *FileParser {
    fp.index += count
    return fp
}

func (fp FileParser) readValue(bytes int64) []byte {
    return fp.file[fp.index: fp.index + bytes]
}

func (fp *FileParser) readValueAndUpdateIndexBy(bytes int64) []byte {
    value := fp.file[fp.index: fp.index + bytes]
    fp.index += bytes
    return value
}

type ConstantPool struct {
    tag int
    constantUtf8 ConstantUtf8
    classInfo ClassInfo
    fieldRef FieldRef
    methodRef MethodRef
    interfaceMethodRef InterfaceMethodRef
    constantStringInfo ConstantStringInfo
    constantIntegerInfo ConstantIntegerInfo
    constantFloatInfo ConstantFloatInfo
    constantLongInfo ConstantLongInfo
    constantDoubleInfo ConstantDoubleInfo 
    constantNameAndTypeInfo ConstantNameAndTypeInfo
    constantMethodHandleInfo ConstantMethodHandleInfo
    constantMethodTypeInfo ConstantMethodTypeInfo
    constantInvokeDynamicInfo ConstantInvokeDynamicInfo
}

type ConstantUtf8 struct {
    length int
    bytes []byte
}

type ClassInfo struct {
    nameIndex int
}

type FieldRef struct {
    classIndex int
    nameAndTypeIndex int
}

type MethodRef struct {
    classIndex int
    nameAndTypeIndex int
}

type InterfaceMethodRef struct {
    classIndex int
    nameAndTypeIndex int
}

type ConstantStringInfo struct {
    stringIndex int
}

type ConstantIntegerInfo struct {
    bytes uint32
}

type ConstantFloatInfo struct {
    bytes uint32
}

type ConstantLongInfo struct {
    highBytes uint64
    lowBytes uint64
}

type ConstantDoubleInfo struct {
    highBytes uint64
    lowBytes uint64
}

type ConstantNameAndTypeInfo struct {
    nameIndex int
    descriptionIndex int
}

type ConstantMethodHandleInfo struct {
    referenceKind int
    referenceIndex int
}

type ConstantMethodTypeInfo struct {
    descriptorIndex int
}

type ConstantInvokeDynamicInfo struct {
    bootstrapMethodAttrIndex int
    nameAndTypeIndex int
}

func ParseFileInfo(file []byte) {
    fileParser := FileParser {
        index: 0,
        file: file,
    }

    fmt.Println("parse file info")
    index := 0;
    if len(fileParser.file) > 0 {
        if !isValidFile(fileParser) {
            return
        }

        fileParser.readValueAndUpdateIndexBy(8)
        // skip the magic number, minor and major versions
        constantPoolCount := binary.BigEndian.Uint16(fileParser.readValueAndUpdateIndexBy(2))
        fmt.Printf("constant pool count %d \n", constantPoolCount)
        index += int(constantPoolCount)

        constantPool := []ConstantPool{}

        for i := 0; i < int(constantPoolCount) - 1; i++ {
            poolItem := ConstantPool{
                    tag: int(file[index]),
                }
            index++;
            updateConstantPoolItem(&poolItem, &fileParser)
            constantPool = append(constantPool, poolItem)
            printConstantPool(poolItem)
        }
    }
}

func updateConstantPoolItem(poolItem *ConstantPool, fileParser *FileParser) {
    switch tag := poolItem.tag; tag {
    case 1:
        length := binary.BigEndian.Uint16(fileParser.readValueAndUpdateIndexBy(2)) 
        string := fileParser.readValueAndUpdateIndexBy(int64(length))
        poolItem.constantUtf8 = ConstantUtf8{
            length: int(length),
            bytes: string,
        }
    case 3:
        numb := binary.BigEndian.Uint32(fileParser.readValueAndUpdateIndexBy(4))
        poolItem.constantIntegerInfo = ConstantIntegerInfo{
            bytes: numb,
        }
    case 4:
        numb := binary.BigEndian.Uint32(fileParser.readValueAndUpdateIndexBy(4))
        poolItem.constantFloatInfo = ConstantFloatInfo {
            bytes: numb,
        }
    case 5:
        high := binary.BigEndian.Uint64(fileParser.readValueAndUpdateIndexBy(4))
        low := binary.BigEndian.Uint64(fileParser.readValueAndUpdateIndexBy(4))
        poolItem.constantLongInfo = ConstantLongInfo {
            highBytes: high,
            lowBytes: low,
        }
    case 6:
        high := binary.BigEndian.Uint64(fileParser.readValueAndUpdateIndexBy(4))
        low := binary.BigEndian.Uint64(fileParser.readValueAndUpdateIndexBy(4))
        poolItem.constantDoubleInfo = ConstantDoubleInfo {
            highBytes: high,
            lowBytes: low,
        }
    case 7:
        nameIndex := binary.BigEndian.Uint64(fileParser.readValueAndUpdateIndexBy(2))    

        poolItem.classInfo = ClassInfo{
            nameIndex: int(nameIndex),
        }

    default:
        fmt.Printf("missing tag %d \n", tag)
    }

}

func isValidFile(file FileParser) bool {
    magic := file.readValue(4)
    str := hex.EncodeToString(magic)
    if (str != "cafebabe") {
        return false
    }
    return true
}

func printConstantPool(cp ConstantPool) {
    fmt.Println("--- Constant Pool ---")
    fmt.Printf("tag %d \n", cp.tag)
    switch tag := cp.tag; tag {
    case 1:
        fmt.Printf("string: %s \n", cp.constantUtf8.bytes)
    case 3:
        fmt.Printf("integer: %d \n", cp.constantIntegerInfo.bytes)
    case 4:
        fmt.Printf("float: %d \n", cp.constantFloatInfo.bytes)
    case 5:
        fmt.Printf("long low: %d \n", cp.constantLongInfo.lowBytes)
        fmt.Printf("long high: %d \n", cp.constantLongInfo.highBytes)
    case 6:
        fmt.Printf("double low: %d \n", cp.constantDoubleInfo.lowBytes)
        fmt.Printf("double high: %d \n", cp.constantDoubleInfo.highBytes)
    case 7:
        fmt.Printf("class %d \n", cp.classInfo.nameIndex)
    }
}
