package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

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
    bytes []byte
}

type ConstantFloatInfo struct {
    bytes []byte
}

type ConstantLongInfo struct {
    highBytes int64
    lowBytes int64
}

type ConstantDoubleInfo struct {
    highBytes int64
    lowBytes int64
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
    index := 0;
    if len(file) > 0 {
        if !isValidFile(file) {
            return
        }

        // skip the magic number, minor and major versions
        index += 8
        constantPoolCount := binary.BigEndian.Uint16(file[index:index + 2])
        fmt.Printf("constant pool count %d \n", constantPoolCount)
        index += int(constantPoolCount)

        constantPool := []ConstantPool{}

        for i := 0; i < int(constantPoolCount) - 1; i++ {
            poolItem := ConstantPool{
                    tag: int(file[index]),
                }
            index++;
            updateConstantPoolItem(&poolItem, file, index)
            constantPool = append(constantPool, poolItem)
        }
        fmt.Println(constantPool)
    }
}

func updateConstantPoolItem(poolItem *ConstantPool, bytes []byte, index int) {
    switch tag := poolItem.tag; tag {
    case 1:
        length := binary.BigEndian.Uint16(bytes[index:index + 2]) 
        index += 2;
        string := bytes[index: index + int(length)]
        index += int(length)
        poolItem.constantUtf8 = ConstantUtf8{
            length: int(length),
            bytes: string,
        }
    case 3,4:
        
    default:
        fmt.Println("missing tag")
    }

}

func isValidFile(file []byte) bool {
    magic := file[:4]
    str := hex.EncodeToString(magic)
    if (str != "cafebabe") {
        return false
    }
    return true
}
