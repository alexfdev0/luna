package video

import (
	"runtime"
	"fmt"
	"os"
	"time"
	"luna_l2/shared"
	"luna_l2/keyboard"
	"luna_l2/component"
	"luna_l2/proxy"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/ncruces/zenity"
	"image"
)

var CommonComponentPathPrefix string = "/usr/local/lib/l2/video/"
var WindowsComponentPathPrefix string = "C:\\Program Files (x86)\\Luna L2\\lib\\l2\\video\\"
var VideoComponent component.Component


// Function definitions
var ReturnFramebuffer func() *image.RGBA



const VertexShaderSrc = `
#version 150

in vec2 inPos;
in vec2 inUV;

out vec2 uv;

void main() {
    uv = inUV;
    gl_Position = vec4(inPos, 0.0, 1.0);
}
` + "\x00"

const FragmentShaderSrc = `
#version 150

in vec2 uv;
out vec4 color;

uniform sampler2D tex;

void main() {
    color = texture(tex, uv);
}
` + "\x00"

func compileShader(src string, t uint32) uint32 {
    shader := gl.CreateShader(t)
    csrc, free := gl.Strs(src)
    gl.ShaderSource(shader, 1, csrc, nil)
    free()
    gl.CompileShader(shader)

    var status int32
    gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
    if status == gl.FALSE {
        var logLen int32
        gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLen)
        log := make([]byte, logLen)
        gl.GetShaderInfoLog(shader, logLen, nil, &log[0])
        panic(string(log))
    }
    return shader
}

func CreateProgram() uint32 {
    vs := compileShader(VertexShaderSrc, gl.VERTEX_SHADER)
    fs := compileShader(FragmentShaderSrc, gl.FRAGMENT_SHADER)

    program := gl.CreateProgram()
    gl.AttachShader(program, vs)
    gl.AttachShader(program, fs)
	gl.BindAttribLocation(program, 0, gl.Str("inPos\x00"))
	gl.BindAttribLocation(program, 1, gl.Str("inUV\x00"))
    gl.LinkProgram(program)

    var status int32
    gl.GetProgramiv(program, gl.LINK_STATUS, &status)
    if status == gl.FALSE {
        var logLen int32
        gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLen)
        log := make([]byte, logLen)
        gl.GetProgramInfoLog(program, logLen, nil, &log[0])
        panic(string(log))
    }

    gl.DeleteShader(vs)
    gl.DeleteShader(fs)

    return program
}

// Frontend code
var Ready bool

var Vertices = []float32 {
	-1, -1, 0, 1,
     1, -1, 1, 1,
     1,  1, 1, 0,

    -1, -1, 0, 1,
     1,  1, 1, 0,
    -1,  1, 0, 0,	
}

func FileOpenDialogue(title string, drive int) {
    ZOpen := func(title string) {
        _path, err := zenity.SelectFile(
            zenity.Title(title),
        )
        if err != nil {
            return
        }
        switch drive {
        case 0:
            shared.Filename = _path
        case 1:
            shared.SDFilename = _path
        case 2:
            shared.OpticalFilename = _path
        }
    }

    switch runtime.GOOS {
    case "darwin":
        ZOpen(title)
    default:
        go ZOpen(title)
    }
}

func ToggleGrab(window *glfw.Window, Grab bool) {
	if Grab == true {
		window.SetTitle("Luna L2 - Press Ctrl+Alt+G to release grab")
		window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	} else {
		window.SetTitle("Luna L2")
		window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	}
}

func ResetAspectRatio(window *glfw.Window, width int, height int) {
	aspect := float32(640) / float32(400)
	actual := float32(width) / float32(height)

	var H int
	var W int
	var X int
	var Y int

	if actual > aspect {
		H = height
		W = int(float32(height) * aspect)
		X = (width - W) / 2
		Y = 0
	} else {
		W = width
		H = int(float32(width) / aspect)
		X = 0
		Y = (height - H) / 2
	}

	gl.Viewport(int32(X), int32(Y), int32(W), int32(H))
}

var FS bool
func ToggleFullscreen(window *glfw.Window) {
	if FS == false {
		window.SetFramebufferSizeCallback(ResetAspectRatio)
		window.SetMonitor(glfw.GetPrimaryMonitor(), 0, 0, 640, 400, 60)
		FS = true
	} else {
		window.SetMonitor(nil, 960, 540, 640, 400, 0)
		FS = false
	}
}

var Grab bool
func InitializeWindow(ComponentName string) {
	wd, _ := os.Getwd()
	prefix := CommonComponentPathPrefix
	ext := ".so"
	if runtime.GOOS == "windows" {
		prefix = WindowsComponentPathPrefix
		ext = ".dll"
	}
	VideoComponent = component.InitializeComponent(prefix + ComponentName + ext)
	component.ReturnComponentFunction(VideoComponent, "InitializePalette").(func())()

	ReturnFramebuffer = component.ReturnComponentFunction(VideoComponent, "ReturnFramebuffer").(func() *image.RGBA)
	proxy.VideoPrintChar = component.ReturnComponentFunction(VideoComponent, "PrintChar").(func(rune, byte, byte))
	proxy.VideoSetCursor = component.ReturnComponentFunction(VideoComponent, "SetCursor").(func(int, int))
	proxy.VideoGetCursor = component.ReturnComponentFunction(VideoComponent, "GetCursor").(func() (int, int))
	proxy.VideoClearVideoMemory = component.ReturnComponentFunction(VideoComponent, "ClearVideoMemory").(func())
	proxy.VideoReadVideoMemory = component.ReturnComponentFunction(VideoComponent, "ReadVideoMemory").(func(uint32) byte)
	proxy.VideoWriteVideoMemory = component.ReturnComponentFunction(VideoComponent, "WriteVideoMemory").(func(uint32, byte))

	err := glfw.Init()
	if err != nil {
		fmt.Println("luna-l2: could not initialize window: ", err)
		os.Exit(1)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.ScaleToMonitor, glfw.True)

	window, err := glfw.CreateWindow(640, 400, "Luna L2", nil, nil)
	if err != nil {
		fmt.Println("luna-l2: could not initialize window: ", err)
		os.Exit(1)
	}
	window.MakeContextCurrent()

	window.SetFramebufferSizeCallback(ResetAspectRatio)

	err = gl.Init()
	if err != nil {
		fmt.Println("luna-l2: could not initialize window: ", err)
		os.Exit(1)
	}

	fbWidth, fbHeight := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(fbWidth), int32(fbHeight))
	gl.ClearColor(0, 0, 0, 1)	

	program := CreateProgram()
	gl.UseProgram(program)

	loc := gl.GetUniformLocation(program, gl.Str("tex\x00"))	
	gl.Uniform1i(loc, 0)

	var vao, vbo uint32

	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(Vertices) * 4,
		gl.Ptr(Vertices),
		gl.STATIC_DRAW,
	)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press || action == glfw.Repeat {
			shift := (mods & glfw.ModShift) != 0
			alt := (mods & glfw.ModAlt) != 0
			ctrl := (mods & glfw.ModControl) != 0

			if ctrl && alt && key == glfw.KeyG {
				if Grab == true {
					ToggleGrab(window, false)
					Grab = false
					return
				}	
			}

			switch key {	
			case glfw.KeyF1:
				// Insert into HDD slot
				if shared.Filename == "" {
					FileOpenDialogue("Select hard disk file", 0)
				} else {
					shared.Filename = ""
				}
				return
			case glfw.KeyF2:
				// Insert into SD slot
				if shared.SDFilename == "" {
					FileOpenDialogue("Select SD/USB file", 1)
				} else {
					shared.SDFilename = ""
				}	
				return
			case glfw.KeyF3:
				// Insert into CD/DVD slot
				if shared.OpticalFilename == "" {
					FileOpenDialogue("Select CD/DVD file", 2)
				} else {
					shared.OpticalFilename = ""
				}	
				return
			case glfw.KeyF4:
				if shared.Debug == true {
					shared.LogOn = false
					shared.Debug = false
				} else {
					shared.LogOn = true
					shared.Debug = true
				}
				return
			case glfw.KeyF5:
				shared.RaiseInterrupt(0xF)
			case glfw.KeyF6:
				f, _ := os.Create("memory_dump.bin")
				f.Write((*shared.Memory)[:])
				f.Close()
			case glfw.KeyF7:
				f, _ := os.Create("audio_dump.bin")
				f.Write((*shared.MemoryAudio)[:])
				f.Close()
			case glfw.KeyF11:
				ToggleFullscreen(window)
				return	
			}	
	

			var char string
			switch key {
			case glfw.KeySpace:
				char = string(byte(0x20))
			case glfw.KeyEnter:
				char = string(byte(0x0A))
			case glfw.KeyBackspace:
				char = string(byte(0xC3))	
			default:
				char = glfw.GetKeyName(key, scancode)	
			}

			if shift {
				char = keyboard.Upper(char)
			} else {
				char = keyboard.Lower(char)
			}

			if len(char) > 0 {
				keyboard.MemoryKeyboard[0] = byte(char[0])
				shared.RaiseInterrupt(0x5)
				shared.SetRegister(0x001b, uint32(char[0]))
			}
		}
	})

	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if button == glfw.MouseButtonLeft && action == glfw.Press {
			if Grab == false {
				ToggleGrab(window, true)
				Grab = true
				return
			}
		}
	})

	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		if Grab == false {
			return
		}

		if xpos > 320 {
			xpos = 320
		} else if xpos < 0 {
			xpos = 0
		}

		if ypos > 200 {
			ypos = 200
		} else if ypos < 0 {
			ypos = 0
		}
	
		ixh := int(xpos) >> 8
		ixl := int(xpos) & 0xFF

		iyh := int(ypos) >> 8
		iyl := int(ypos) & 0xFF

		keyboard.MemoryMouse[2] = byte(ixh)
		keyboard.MemoryMouse[3] = byte(ixl)
		keyboard.MemoryMouse[6] = byte(iyh)
		keyboard.MemoryMouse[7] = byte(iyl)
		
		shared.RaiseInterrupt(0x12)
	})

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA8,
		320, 200,
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		nil,
	)	

	os.Chdir(wd)
	next := time.Now()	
	for !window.ShouldClose() {
		Ready = true
		img := ReturnFramebuffer() 

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.TexSubImage2D(
			gl.TEXTURE_2D,
			0,
			0, 0,
			320, 200,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(img.Pix),
		)

		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.UseProgram(program)
		gl.BindVertexArray(vao)

		gl.DrawArrays(gl.TRIANGLES, 0, 6)

		next = next.Add(time.Second / 70)
		sleep := time.Until(next)
		if sleep > 0 {
			time.Sleep(sleep)
		} else {
			next = time.Now()
		}

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
