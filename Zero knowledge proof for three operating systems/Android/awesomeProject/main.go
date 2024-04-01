package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"

	"github.com/consensys/gnark-crypto/ecc"
	bn254 "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

var (
	images      *glutil.Images
	fps         *debug.FPS
	program     gl.Program
	position    gl.Attrib
	scan        gl.Uniform
	color       gl.Attrib
	positionbuf gl.Buffer
	colorbuf    gl.Buffer
	touchLocX   float32
	touchLocY   float32
)

func main() {
	app.Main(func(a app.App) {
		var glctx gl.Context
		var sz size.Event
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
					onStart(glctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					onStop(glctx)
					glctx = nil
				}
			case size.Event:
				sz = e
				touchLocX = float32(sz.WidthPt / 1.5)
				touchLocY = float32(sz.HeightPt / 1.5)
			case paint.Event:
				if glctx == nil || e.External {
					// As we are actively painting as fast as
					// we can (usually 60 FPS), skip any paint
					// events sent by the system.
					continue
				}
				onPaint(glctx, sz)
				a.Publish()
				// Drive the animation by preparing to paint the next frame
				// after this one is shown.
				a.Send(paint.Event{})
			case touch.Event:
				touchLocX = e.X / sz.PixelsPerPt
				touchLocY = e.Y / sz.PixelsPerPt
			case key.Event:
				for i := 0; i < 100; i++ {
					startTime := time.Now()
					proof, respro := Makeproof()
					elapsedTime := time.Since(startTime)
					Proof = proof
					if respro {
						file, err := os.OpenFile("output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
						if err != nil {
							fmt.Println("could not open file:", err)
							return
						}
						defer file.Close()
						// 将字符串写入文件
						content := elapsedTime.String() + " \n"
						_, err = file.WriteString(content)
						if err != nil {
							fmt.Println("Unable to write to file:", err)
							return
						}
					}
				}
			}

		}
	})
}

func onStart(glctx gl.Context) {
	var err error
	program, err = glutil.CreateProgram(glctx, vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}

	/*
		Three Variables in Opengl
		The uniform variable is a variable passed by an external application program to the (vertex and fragment) shader. Therefore, it is assigned by the application through the function glUniform().
		Within the (vertex and fragment) shader program, the uniform variable is like a constant (const) in C language, which cannot be modified by the shader program. (Shader can only be used and cannot be changed)

		The attribute variable is a variable that can only be used in the vertex shader. (It cannot declare attribute variables in the fragment shader, nor can it be used in the fragment shader)
		Generally, attribute variables are used to represent data for certain vertices, such as vertex coordinates, normals, texture coordinates, vertex colors, etc.
		In an application, the function glBindAttribLocation() is generally used to bind the position of each attribute variable, and then the function glVertexAttribPointer() is used to assign values to each attribute variable.

		The varying variable is used for data transfer between vertex and fragment shader. Generally, the vertex shader modifies the value of the variable variable, and then the fragment shader uses the value of the variable.
		Therefore, the declaration of varying variables between vertex and fragment shader must be consistent. Application cannot use this variable.
	*/

	position = glctx.GetAttribLocation(program, "position")
	color = glctx.GetAttribLocation(program, "color")
	scan = glctx.GetUniformLocation(program, "scan")

	/*
		VBO allows the usage identifier to take the following 9 values:

		Gl.STATIC-DRAW
		Gl.STATIC-READ
		Gl.STATIC-COPY

		Gl.DYNAMIC-DRAW
		Gl.DYNAMIC-READ
		Gl.DYNAMIC-COPY

		Gl.STREAM-DRAW
		Gl.STREAM-READ
		Gl.STREAM-COPY

		"Static" means that the data in VBO will not be changed (modified once, used multiple times),
		"Dynamic" means that data can be frequently modified (modified multiple times, used multiple times),
		"Stream" means that the data is different for each frame (once modified, once used).

		"Draw" means that the data will be sent to the GPU for drawing,
		"Read" means that the data will be read by the user's application,
		"Copy" means that the data will be used for drawing and reading.

		Note that when using VBO, only draw is effective, while copy and read will mainly play a role in pixel buffer (PBO) and frame buffer (FBO).
		The system will allocate the optimal storage location for buffer objects based on the usage identifier, such as gl STATIC-DRAW and gl STREAM-DRAW allocates graphics memory,
		Gl.DYNAMIC-DRAW allocates AGP, and any buffer objects related to READ_ will be stored in the system or AGP because this makes the data easier to read and write to
	*/
	positionbuf = glctx.CreateBuffer()
	glctx.BindBuffer(gl.ARRAY_BUFFER, positionbuf)
	glctx.BufferData(gl.ARRAY_BUFFER, triangleData, gl.STATIC_DRAW)

	colorbuf = glctx.CreateBuffer()
	glctx.BindBuffer(gl.ARRAY_BUFFER, colorbuf)
	glctx.BufferData(gl.ARRAY_BUFFER, colorData, gl.STATIC_DRAW)

	images = glutil.NewImages(glctx)
	fps = debug.NewFPS(images)

	// fmt.Println(position.String(),color.String(),offset.String())//Attrib(0) Uniform(1) Uniform(0)
	// TODO(crawshaw): the debug package needs to put GL state init here
	// Can this be an event.Register call now??

	respvk := Makepvk()
	if respvk {
		fmt.Println("Zero knowledge proof initialization successful")
	}
}

func onStop(glctx gl.Context) {
	glctx.DeleteProgram(program)
	glctx.DeleteBuffer(positionbuf)
	glctx.DeleteBuffer(colorbuf)
	fps.Release()
	images.Release()
}

func onPaint(glctx gl.Context, sz size.Event) {
	//gl.Enable(gl.DEPTH_TEST)
	//gl.DepthFunc(gl.LESS)
	glctx.ClearColor(1, 1, 1, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)
	glctx.Clear(gl.DEPTH_BUFFER_BIT)
	glctx.UseProgram(program)
	glctx.UniformMatrix4fv(scan, []float32{
		touchLocX/float32(sz.WidthPt)*4 - 2, 0, 0, 0,
		0, touchLocY/float32(sz.HeightPt)*4 - 2, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 1,
	})
	/*The glVertexAttribPointer specifies the data format and position of the vertex attribute array with an index value of index during rendering. Call the gl. vertexAttribPointer() method,
	Connect the universal attribute index corresponding to a certain attribute in the vertex shader to the bound webGLBuffer object.
	Index specifies the index value of the vertex attribute to be modified
	Size specifies the number of components for each vertex attribute. Must be 1, 2, 3, or 4. The initial value is 4. (For example, position is composed of 3 (x, y, z), while colors are 4 (r, g, b, a))
	Type specifies the data type of each component in the array. The available symbol constants are GL-BYTE, GL_UNSIGNED-BYTE, GL-SHORT, GL_UNSIGNED-SHORT, GL-FIXED,
	And GL-FLOAT, with an initial value of GL-FLOAT.
	Normalized specifies whether fixed point data values should be normalized (GL-TRUE) or directly converted to fixed point values (GL-FALSE) when accessed.
	Strude specifies the offset between consecutive vertex attributes. If it is 0, then the vertex attributes will be understood as: they are tightly arranged together. The initial value is 0.
	Pointer specifies the offset of the first component in the first vertex attribute of the array. This array is bound to GL-ARRAY-BUFFER and stored in a buffer. The initial value is 0;
	*/

	glctx.BindBuffer(gl.ARRAY_BUFFER, positionbuf)
	glctx.EnableVertexAttribArray(position)
	defer glctx.DisableVertexAttribArray(position)
	glctx.VertexAttribPointer(position, coordsPerVertex, gl.FLOAT, false, 0, 0)
	glctx.DrawArrays(gl.TRIANGLES, 0, vertexCount)

	glctx.BindBuffer(gl.ARRAY_BUFFER, colorbuf)
	glctx.EnableVertexAttribArray(color)
	defer glctx.DisableVertexAttribArray(color)
	glctx.VertexAttribPointer(color, colorsPerVertex, gl.FLOAT, false, 0, 0)
	glctx.DrawArrays(gl.TRIANGLES, 0, vertexCount)

	fps.Draw(sz)
}

var triangleData = f32.Bytes(binary.LittleEndian,
	0.0, 0.5, 0.0, // top left
	-0.5, -0.5, 0.0, // bottom left
	0.5, -0.5, 0.0, // bottom right
)

var colorData = f32.Bytes(binary.LittleEndian,
	1.0, 0.0, 0.0, 1, // red
	0.0, 1.0, 0.0, 1, // green
	0.0, 0.0, 1.0, 1, // blue
)

const (
	coordsPerVertex = 3
	vertexCount     = 3
	colorsPerVertex = 4
)

const vertexShader = `#version 100
uniform mat4 scan;
attribute vec4 position;
attribute vec4 color;
varying vec4 vColor;
void main() {
	gl_Position = position*scan;
	vColor = color;
}`

const fragmentShader = `#version 100
precision mediump float;
varying vec4 vColor;
void main() {
	gl_FragColor = vColor;
}`

type Circuit struct {
	PreImage frontend.Variable //Hkey
	Ht       frontend.Variable `gnark:",public"`
}

func (circuit *Circuit) Define(curveID ecc.ID, api frontend.API) error {
	Htt, _ := mimc.NewMiMC(TKeyt, curveID, api)
	Htt.Write(circuit.PreImage)
	api.AssertIsEqual(Htt.Sum(), circuit.Ht)
	return nil
}

var Key = "544893542849cbded60ce6880112235267e0908f580e9db65e138663fb15f78f" //sha256
var t = "11111111"                                                           //tao
var T = "00000000"                                                           //time
var TKeyt = Key + t + T

func mimcHash(data []byte, str string) string {
	f := bn254.NewMiMC(str)
	f.Write(data)
	hash := f.Sum(nil)
	hashInt := big.NewInt(0).SetBytes(hash)
	return hashInt.String()
}

var R1cs frontend.CompiledConstraintSystem
var Pk groth16.ProvingKey
var Vk groth16.VerifyingKey
var Proof groth16.Proof

func Makepvk() bool {
	var circuit Circuit
	r1cs, err := frontend.Compile(ecc.BN254, backend.GROTH16, &circuit)
	if err != nil {
		fmt.Printf("Compile failed : %v\n", err)
		return false
	}

	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		fmt.Printf("Setup failed\n")
		return false
	}
	R1cs = r1cs
	Pk = pk
	Vk = vk
	return true
}
func Makeproof() (groth16.Proof, bool) {
	preImage := []byte{0x01, 0x02, 0x03}
	preImageHt := []byte{0x01, 0x02, 0x03}
	ht := mimcHash(preImageHt, TKeyt)
	witness := &Circuit{
		PreImage: frontend.Value(preImage),
		Ht:       frontend.Value(ht),
	}
	proof, err := groth16.Prove(R1cs, Pk, witness)
	if err != nil {
		fmt.Printf("Prove failed： %v\n", err)
		return nil, false
	}
	Proof = proof
	fmt.Println("Prove generated successfully!")
	return proof, true
}
func Makeverify() bool {

	preImageHt := []byte{0x01, 0x02, 0x03}
	ht := mimcHash(preImageHt, TKeyt)
	publicWitness := &Circuit{
		Ht: frontend.Value(ht),
	}

	err := groth16.Verify(Proof, Vk, publicWitness)
	if err != nil {
		fmt.Printf("verification failed: %v\n", err)
		return false
	}
	return true
}
