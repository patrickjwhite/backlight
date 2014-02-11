package main

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/niemeyer/qml"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	UI = `
	import QtQuick 2.2
	import QtQuick.Window 2.0
	import QtQuick.Controls 1.1
	import QtQuick.Controls.Styles 1.1

	ApplicationWindow {
		id: backlight
		flags: Qt.FramelessWindowHint
		visible: true
		title: qsTr("backlight")
		width: 500
		height: 50
		x: (Screen.width - width) / 2
		y: (Screen.height - height) / 2
		color: "transparent"

		property real slide
		signal onSlide(real value)

		Rectangle {
			anchors.centerIn: parent
			width: parent.width
			height: 50
			color: "transparent"

			Rectangle {
				anchors.fill: parent
				radius: 25
				opacity: 0.3
				color: "gray"
			}

			Slider {
				anchors.centerIn: parent
				width: backlight.width - 16
				height: backlight.height
				value: backlight.slide
				focus: true
				onValueChanged: backlight.onSlide(value)
				Keys.onSpacePressed: ctrl.close()
				Keys.onEscapePressed: ctrl.close()

				style: SliderStyle {
					groove: Rectangle {
						implicitHeight: 8
						radius: 4
						color: "gray"
					}
					handle: Rectangle {
						anchors.centerIn: parent
						color: control.pressed ? "white" : "lightgray"
						border.color: "gray"
						border.width: 2
						width: 34
						height: 34
						radius: 17
					}
				}
			}
		}
	}`
	BRIGHTNESS     = "/sys/class/backlight/intel_backlight/brightness"
	MAX_BRIGHTNESS = "/sys/class/backlight/intel_backlight/max_brightness"
)

func main() {
	qml.Init(nil)
	engine := qml.NewEngine()

	component, err := engine.LoadString("ui.qml", UI)
	if err != nil {
		panic(err)
	}

	window, err := NewWindow(component)
	if err != nil {
		panic(err)
	}
	defer window.Close()

	engine.Context().SetVar("ctrl", window)
	window.Show()
	window.AlwaysOnTop()
	window.Wait()
}

type Window struct {
	*qml.Window
	brightnessFile *os.File
	maxBrightness  float64
}

func NewWindow(component qml.Object) (window *Window, err error) {
	window = new(Window)

	contents, err := ioutil.ReadFile(MAX_BRIGHTNESS)
	if err != nil {
		return nil, err
	}

	window.maxBrightness, err = strconv.ParseFloat(strings.TrimSpace(string(contents)), 64)
	if err != nil {
		return nil, err
	}

	window.brightnessFile, err = os.OpenFile(BRIGHTNESS, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	defer func() { if err != nil { window.brightnessFile.Close() } }()

	contents, err = ioutil.ReadAll(window.brightnessFile)
	if err != nil {
		return nil, err
	}

	brightness, err := strconv.ParseFloat(strings.TrimSpace(string(contents)), 64)
	if err != nil {
		return nil, err
	}

	window.Window = component.CreateWindow(nil)
	window.Set("slide", math.Pow(brightness/window.maxBrightness, 0.5))
	window.On("onSlide", window.OnSlide)

	return window, nil
}

func (window *Window) SetBrightness(value float64) error {
	window.brightnessFile.Truncate(0)
	brightness := int(math.Pow(value, 2.0) * window.maxBrightness)
	_, err := window.brightnessFile.WriteString(strconv.Itoa(brightness))
	return err
}

func (window *Window) OnSlide(value float64) {
	window.SetBrightness(value)
}

func (window *Window) Close() {
	if window.Window != nil {
		window.Window.Destroy()
		window.Window = nil
	}

	if window.brightnessFile != nil {
		window.brightnessFile.Close()
		window.brightnessFile = nil
	}
}

const (
	_NET_WM_STATE_REMOVE = 0
	_NET_WM_STATE_ADD    = 1
	_NET_WM_STATE_TOGGLE = 2
)

func (window *Window) AlwaysOnTop() {
	xid := xproto.Window(window.PlatformId())
	X, err := xgb.NewConn()
	if err != nil {
		log.Println(err)
		return
	}
	defer X.Close()

	state, err := xproto.InternAtom(X, false, uint16(len("_NET_WM_STATE")),
		"_NET_WM_STATE").Reply()
	if err != nil {
		log.Println(err)
		return
	}

	stateAbove, err := xproto.InternAtom(X, false,
		uint16(len("_NET_WM_STATE_ABOVE")), "_NET_WM_STATE_ABOVE").Reply()
	if err != nil {
		log.Println(err)
		return
	}

	evt := xproto.ClientMessageEvent{
		Window: xid,
		Format: 32,
		Type:   state.Atom,
		Data: xproto.ClientMessageDataUnionData32New([]uint32{
			_NET_WM_STATE_ADD,
			uint32(stateAbove.Atom),
			0,
			0,
			0,
		}),
	}

	err = xproto.SendEventChecked(X, false, xproto.Setup(X).DefaultScreen(X).Root,
		xproto.EventMaskSubstructureRedirect|xproto.EventMaskSubstructureNotify,
		string(evt.Bytes())).Check()
	if err != nil {
		log.Println(err)
	}
}
