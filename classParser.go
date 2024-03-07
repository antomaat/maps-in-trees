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
    //fmt.Println(value)
    fp.index += bytes
    return value
}

type ParseResult struct {
    className string
    superClassName string
    fields []string
}

type ConstantPool struct {
    tag uint 
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
    classIndex uint16 
    nameAndTypeIndex uint16 
}

type MethodRef struct {
    classIndex uint16 
    nameAndTypeIndex uint16 
}

type InterfaceMethodRef struct {
    classIndex uint16 
    nameAndTypeIndex uint16 
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
    nameIndex uint16 
    descriptionIndex uint16 
}

type ConstantMethodHandleInfo struct {
    referenceKind uint16 
    referenceIndex uint16
}

type ConstantMethodTypeInfo struct {
    descriptorIndex uint16 
}

type ConstantInvokeDynamicInfo struct {
    bootstrapMethodAttrIndex uint16 
    nameAndTypeIndex uint16 
}


type FieldInfo struct {
    accessFlags uint16 
    name string 
    descriptor string 
    attributeCount uint16 
    attributes []Attribute
}

type Attribute struct {
    name string
    attributeLength uint32
    bytes []byte
}

func ParseFileInfo(file []byte) {
    fileParser := FileParser {
        index: 0,
        file: file,
    }

    result := ParseResult{}

    fmt.Println("parse file info")
    if len(fileParser.file) > 0 {
        if !isValidFile(fileParser) {
            return
        }

        fileParser.readValueAndUpdateIndexBy(8)
        // skip the magic number, minor and major versions
        constantPoolCount := binary.BigEndian.Uint16(fileParser.readValueAndUpdateIndexBy(2))
        fmt.Printf("constant pool count %d \n", constantPoolCount)
        constantPool := []ConstantPool{}

        //for i := 0; i < 1; i++ {
        for i := 0; i < int(constantPoolCount) - 1; i++ {
            nextTag := fileParser.readValueAndUpdateIndexBy(1)
            poolItem := ConstantPool{
                tag: uint(nextTag[0]),
            }
            updateConstantPoolItem(&poolItem, &fileParser)
            constantPool = append(constantPool, poolItem)
            //printConstantPool(poolItem)
        } 
        for i := 0; i < len(constantPool); i++ {
            parseFinishedConstantPool(constantPool[i], constantPool)
        }

        //access flags
        fileParser.readValueAndUpdateIndexBy(2)
        // add class name
        result.className = parseNameFromNextBytes(&fileParser, constantPool) 

        // add super class name
        result.superClassName = parseNameFromNextBytes(&fileParser, constantPool) 
        fmt.Printf("class name %s \n", result.className)
        fmt.Printf("super class name %s \n", result.superClassName)
        // interfaces
        // interface count
        interfaceCount := binary.BigEndian.Uint16(fileParser.readValueAndUpdateIndexBy(2))
        fmt.Printf("interface count %d \n", interfaceCount)
        if interfaceCount > 0 {
            panic("interface support is not present")
        }
        // add fields
        fieldsCount := binary.BigEndian.Uint16(fileParser.readValueAndUpdateIndexBy(2)) 
        fmt.Printf("fields count %d \n", fieldsCount)
        for i:= 0; i < int(fieldsCount); i++ {
            fieldInfo := parseFieldInfo(&fileParser, constantPool)
            fmt.Printf("field info name %s \n", fieldInfo.name)
        }
    }
}

func parseFieldInfo(fp *FileParser, constantPool []ConstantPool) FieldInfo {
    fieldInfo := FieldInfo{}

    fieldInfo.accessFlags = binary.BigEndian.Uint16(fp.readValueAndUpdateIndexBy(2))
    name := parseNameIndex(constantPool, uint(binary.BigEndian.Uint16(fp.readValueAndUpdateIndexBy(2))))
    descriptor := parseNameIndex(constantPool, uint(binary.BigEndian.Uint16(fp.readValueAndUpdateIndexBy(2))))

    fieldInfo.name = name
    fieldInfo.descriptor = descriptor 
    fieldInfo.attributeCount = binary.BigEndian.Uint16(fp.readValueAndUpdateIndexBy(2))
    for i := 0; i < int(fieldInfo.attributeCount); i++ {
        attribute := Attribute{}
        nameIndex:= binary.BigEndian.Uint16(fp.readValueAndUpdateIndexBy(2))
        attribute.name = parseNameIndex(constantPool, uint(nameIndex)) 
        attributeLenBytes := fp.readValueAndUpdateIndexBy(4)
        fmt.Printf("attr len bytes %v \n", attributeLenBytes)
        fmt.Printf("attr len bytes %d \n", binary.BigEndian.Uint32(attributeLenBytes))
        attribute.attributeLength = binary.BigEndian.Uint32(attributeLenBytes)
        attribute.bytes = fp.readValueAndUpdateIndexBy(int64(attribute.attributeLength))
        fieldInfo.attributes = append(fieldInfo.attributes, attribute)
    }
    return fieldInfo
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
        nameIndex := binary.BigEndian.Uint16(fileParser.readValueAndUpdateIndexBy(2))    

        poolItem.classInfo = ClassInfo{
            nameIndex: int(nameIndex),
        }
    case 8:
        stringIndex := binary.BigEndian.Uint16(fileParser.readValueAndUpdateIndexBy(2))    
        poolItem.constantStringInfo = ConstantStringInfo {
            stringIndex: int(stringIndex),
        }
    case 9:
        classIndex := binary.BigEndian.Uint16(fileParser.readValueAndUpdateIndexBy(2))
        nameAndTypeIndex := binary.BigEndian.Uint16(fileParser.readValueAndUpdateIndexBy(2)) 

        poolItem.fieldRef = FieldRef {
            classIndex: classIndex,
            nameAndTypeIndex: nameAndTypeIndex,
        }
    case 10:
        classIndex := binary.BigEndian.Uint16(fileParser.readValueAndUpdateIndexBy(2))
        nameAndTypeIndex := binary.BigEndian.Uint16(fileParser.readValueAndUpdateIndexBy(2)) 

        poolItem.methodRef = MethodRef {
            classIndex: classIndex,
            nameAndTypeIndex: nameAndTypeIndex,
        }
    case 11:
        classIndex := binary.BigEndian.Uint16(fileParser.readValueAndUpdateIndexBy(2))
        nameAndTypeIndex := binary.BigEndian.Uint16(fileParser.readValueAndUpdateIndexBy(2)) 

        poolItem.interfaceMethodRef = InterfaceMethodRef {
            classIndex: classIndex,
            nameAndTypeIndex: nameAndTypeIndex,
        }
    case 12:
        poolItem.constantNameAndTypeInfo = ConstantNameAndTypeInfo{
            nameIndex: binary.BigEndian.Uint16(fileParser.readValueAndUpdateIndexBy(2)),
            descriptionIndex: binary.BigEndian.Uint16(fileParser.readValueAndUpdateIndexBy(2)),
        }
    case 15:
        poolItem.constantMethodHandleInfo = ConstantMethodHandleInfo{
            referenceKind: binary.BigEndian.Uint16(fileParser.readValueAndUpdateIndexBy(2)),
            referenceIndex: binary.BigEndian.Uint16(fileParser.readValueAndUpdateIndexBy(2)),
        }
    case 16:
        poolItem.constantMethodTypeInfo = ConstantMethodTypeInfo{
            descriptorIndex: binary.BigEndian.Uint16(fileParser.readValueAndUpdateIndexBy(2)),
        }
    case 18:
        poolItem.constantInvokeDynamicInfo = ConstantInvokeDynamicInfo{
            bootstrapMethodAttrIndex: binary.BigEndian.Uint16(fileParser.readValueAndUpdateIndexBy(2)),
            nameAndTypeIndex: binary.BigEndian.Uint16(fileParser.readValueAndUpdateIndexBy(2)),
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
    case 8:
        fmt.Printf("string %d \n", cp.constantStringInfo.stringIndex)
    case 9:
        fmt.Printf("fieldref class index %d \n", cp.fieldRef.classIndex)
    case 10:
        fmt.Printf("method ref class index %d \n", cp.methodRef.classIndex)
        fmt.Printf("method ref name and type index %d \n", cp.methodRef.nameAndTypeIndex)
    }
}

func parseFinishedConstantPool(poolItem ConstantPool, constantPool []ConstantPool) {
    //fmt.Printf("item tag %d \n", poolItem.tag)
    switch tag := poolItem.tag; tag {
    case 7:
        fmt.Println(parseNameIndex(constantPool, uint(poolItem.classInfo.nameIndex -1)))
    case 9:
        fmt.Println("tag is 9")
        fieldRef := poolItem.fieldRef
        classIndex:= constantPool[fieldRef.classIndex -1].classInfo
        name := parseNameIndex(constantPool, uint(classIndex.nameIndex))
        fmt.Printf("string %s\n", name)
        parseNameAndType(constantPool, uint(fieldRef.nameAndTypeIndex - 1))
    case 10:
        fmt.Println("tag is 10")
        methodRef := poolItem.methodRef
        classIndex:= constantPool[methodRef.classIndex -1].classInfo
        name := parseNameIndex(constantPool, uint(classIndex.nameIndex))
        fmt.Printf("string %s\n", name)
        parseNameAndType(constantPool, uint(methodRef.nameAndTypeIndex - 1))
    case 11:
        fmt.Println("tag is 11")
        methodRef := poolItem.interfaceMethodRef
        classIndex:= constantPool[methodRef.classIndex -1].classInfo
        name := parseNameIndex(constantPool, uint(classIndex.nameIndex))
        fmt.Printf("string %s\n", name)
        parseNameAndType(constantPool, uint(methodRef.nameAndTypeIndex - 1))
    }

}

func parseNameIndex(cp []ConstantPool, index uint) string {
    className := cp[index].constantUtf8.bytes
    return string(className)
}

func parseNameAndType(cp []ConstantPool, index uint) string {
    nameAndType := cp[index].constantNameAndTypeInfo
    name := string(cp[nameAndType.nameIndex - 1].constantUtf8.bytes)
    description:= string(cp[nameAndType.descriptionIndex -1].constantUtf8.bytes)
    fmt.Printf("name: %s \n", name)
    fmt.Printf("description: %s \n", description)
    return ""
}

func parseNameFromNextBytes(fileParser *FileParser, constantPool []ConstantPool) string {
        superClassNameIndex := binary.BigEndian.Uint16(fileParser.readValueAndUpdateIndexBy(2)) 
        superClassInfo := constantPool[superClassNameIndex - 1].classInfo
        return parseNameIndex(constantPool, uint(superClassInfo.nameIndex - 1))
}
