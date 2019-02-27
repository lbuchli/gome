package graphics

import (
	"bufio"
	"gitlocal/gome"
	"io"
	"strconv"
	"strings"
)

var OBJ_VERTEX_LAYOUT = VertexLayout{layout: []ElementType{VEC3, VEC2, VEC3}}

// A VertexObject contains all data needed by a renderer to draw a single object.
type VertexObject struct {
	Indices            []uint
	Vertices           []float32
	TextureCoordinates []float32
	Normals            []float32
}

// A OBJFileReader reads a .obj file.
type OBJFileReader struct{}

// A objVertex is a vertex generated by a .obj file.
type objVertex struct {
	position gome.FloatVector3
	uv       gome.FloatVector2
	normal   gome.FloatVector3
}

// IsSimilarTo checks if a vertex is similar to another vertex.
func (ov objVertex) IsSimilarTo(position gome.FloatVector3, uv gome.FloatVector2, normal gome.FloatVector3) bool {
	return ov.position.IsSimilarTo(position) &&
		ov.uv.IsSimilarTo(uv) &&
		ov.normal.IsSimilarTo(normal)
}

// Check checks if the file type is obj.
func (ofr *OBJFileReader) Check(file io.Reader) bool {
	// TODO
	return true
}

// Extension returns the default file extention for this file type.
func (ofr *OBJFileReader) Extension() string { return "obj" }

// Data returns the data of the whole file in a more readable format.
func (ofr *OBJFileReader) Data(file io.Reader) (data VertexArray, err error) {

	reader := bufio.NewReader(file)

	// store original data temporary so it can be further processed
	tempPositionIndices := []uint32{}
	tempUVIndices := []uint32{}
	tempNormalIndices := []uint32{}
	tempPositions := []gome.FloatVector{}
	tempUVs := []gome.FloatVector{}
	tempNormals := []gome.FloatVector{}

	for {
		// read the file line by line
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// stop at the end of the file
				break
			}

			return data, err
		}

		words := strings.Split(line, " ")

		// check the first word to identify the line type
		switch words[0] {
		case "#": // comments
			continue
		case "v": // vertices
			tempPositions = append(tempPositions, gome.FloatVector3{
				X: convertToFloat32(words[1]),
				Y: convertToFloat32(words[2]),
				Z: convertToFloat32(words[3]),
			})
		case "vt": // texture coordinates
			tempUVs = append(tempUVs, gome.FloatVector2{
				X: convertToFloat32(words[1]),
				Y: convertToFloat32(words[2]),
			})
		case "vn": // vertex normals
			tempNormals = append(tempNormals, gome.FloatVector3{
				X: convertToFloat32(words[1]),
				Y: convertToFloat32(words[2]),
				Z: convertToFloat32(words[3]),
			})
		case "s":
			continue // TODO
		case "f": // indices
			// for every word except the first one
			for _, word := range words[1:] {
				// append the indices for one vertex
				split := strings.Split(word, "/")
				tempPositionIndices = append(tempPositionIndices, convertToUint32(split[0]))
				tempUVIndices = append(tempUVIndices, convertToUint32(split[1]))
				tempNormalIndices = append(tempNormalIndices, convertToUint32(split[0]))
			}
		}
	}

	vertices := []objVertex{}
	indices := []uint32{}

	// process temporary data
	for i := 0; i < len(tempPositionIndices); i++ {

		var matched bool

		// only use the uv if there are textures
		uv := gome.FloatVector2{}
		if len(tempUVs) > 0 {
			uv = tempUVs[tempUVIndices[i]-1].(gome.FloatVector2)
		}

		// check if there is a match
		for j, vertex := range vertices {
			if vertex.IsSimilarTo(
				tempPositions[tempPositionIndices[i]-1].(gome.FloatVector3),
				uv,
				tempNormals[tempNormalIndices[i]-1].(gome.FloatVector3),
			) {
				// there is a similar vertex, use it instead
				indices = append(indices, uint32(j))
				matched = true
				break
			}
		}

		// if there is no match, add to the vertices and refer to the new entry
		if !matched {
			indices = append(indices, uint32(len(vertices)))
			vertices = append(vertices, objVertex{
				position: tempPositions[tempPositionIndices[i]-1].(gome.FloatVector3),
				uv:       uv,
				normal:   tempNormals[tempNormalIndices[i]-1].(gome.FloatVector3),
			})
		}
	}

	// set vertex data
	data.SetLayout(OBJ_VERTEX_LAYOUT)
	data.SetData(0, tempPositions)
	data.SetData(1, tempUVs)
	data.SetData(2, tempNormals)
	data.SetIndexData(indices)

	return
}

// convertToFloat32 silently converts a string to a float32.
func convertToFloat32(input string) float32 {
	result, err := strconv.ParseFloat(input, 32)
	if err != nil {
		return 0
	}
	return float32(result)
}

// convertToUint silently converts a string to a uint32.
func convertToUint32(input string) uint32 {
	result, err := strconv.ParseUint(input, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(result)
}
