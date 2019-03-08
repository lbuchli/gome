package graphics

import (
	"bufio"
	"fmt"
	"gitlocal/gome"
	"image"
	"image/draw"
	_ "image/png"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

var OBJ_VERTEX_LAYOUT = VertexLayout{layout: []ElementType{FVEC3, FVEC2, FVEC3}}

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
func (ofr *OBJFileReader) Data(file io.Reader) (data VertexArray, texture uint32, err error) {

	reader := bufio.NewReader(file)

	// store original data temporary so it can be further processed
	tempPositionIndices := []uint32{}
	tempUVIndices := []uint32{}
	tempNormalIndices := []uint32{}
	tempPositions := []gome.FloatVector{}
	tempUVs := []gome.FloatVector{}
	tempNormals := []gome.FloatVector{}

	material := ""

	for {
		// read the file line by line
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// stop at the end of the file
				break
			}

			return data, texture, err
		}

		// clip the newline char from the line
		line = line[:len(line)-1]

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
		case "usemtl": // texture / material
			material = words[1] // TODO this will bug out when file name contains a space
		case "f": // indices
			// for every word except the first one
			for _, word := range words[1:] {
				// append the indices for one vertex
				split := strings.Split(word, "/")
				tempPositionIndices = append(tempPositionIndices, convertToUint32(split[0])-1)
				tempUVIndices = append(tempUVIndices, convertToUint32(split[1])-1)
				tempNormalIndices = append(tempNormalIndices, convertToUint32(split[2])-1)
			}
		}
	}

	vertices := []objVertex{}
	indices := []uint32{}

	// process temporary data
	// TODO fix indexing
	for i := 0; i < len(tempPositionIndices); i++ {

		var matched bool

		// only use the uv if there are textures
		uv := gome.FloatVector2{}
		if len(tempUVs) > 0 {
			uv = tempUVs[tempUVIndices[i]].(gome.FloatVector2)
		}

		// check if there is a match
		for j, vertex := range vertices {
			if vertex.IsSimilarTo(
				tempPositions[tempPositionIndices[i]].(gome.FloatVector3),
				uv,
				tempNormals[tempNormalIndices[i]].(gome.FloatVector3),
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
				position: tempPositions[tempPositionIndices[i]].(gome.FloatVector3),
				uv:       uv,
				normal:   tempNormals[tempNormalIndices[i]].(gome.FloatVector3),
			})
		}
	}

	rawData := []float32{}

	// convert vertex data to a raw float32 array
	for _, vertex := range vertices {
		rawData = append(rawData,
			vertex.position.X, vertex.position.Y, vertex.position.Z,
			vertex.uv.X, vertex.uv.Y,
			vertex.normal.X, vertex.normal.Y, vertex.normal.Z,
		)
	}

	// set vertex data
	data.SetLayout(OBJ_VERTEX_LAYOUT)
	data.SetData(rawData)
	data.SetIndexData(indices)

	// if the material was set, set the texture
	if len(material) > 0 {
		texture, err = newTexture(material)
		if err != nil {
			return data, texture, err
		}
	} else {
		texture = 0
	}

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

// newTexture generates a new OpenGL texture from a file.
// Source: https://gist.github.com/errcw/e3311a0ed1a1c0113a92
func newTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, err
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil
}
