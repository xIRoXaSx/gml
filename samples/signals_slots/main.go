/*
 * GML - Go QML
 *
 * The MIT License (MIT)
 *
 * Copyright (c) 2019 Roland Singer <roland.singer[at]desertbit.com>
 * Copyright (c) 2019 Sebastian Borchers <sebastian[at]desertbit.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/desertbit/gml"
	_ "github.com/desertbit/gml/samples/signals_slots/testy"
)

type A struct{} // TODO: test as property.

type Bridge struct {
	gml.Object
	_ struct {
		state   int                                                                            `gml:"property"`
		clicked func(i int, v *gml.Variant)                                                    `gml:"slot"`
		greet   func(i1 uint8, i2 int32, i3 int, s string, r rune, b byte, bb bool, bs []byte) `gml:"signal"`
	}
}

func (b *Bridge) clicked(i int, v *gml.Variant) {
	b.emitGreet(1, 2, 3, "foo", '本', 4, true, []byte{1, 2, 3})
}

func (b *Bridge) stateChanged() {
	println("state changed")
}

type Model struct {
	*gml.ListModel

	data []string
}

func newModel() *Model {
	m := &Model{data: []string{"1", "2", "3"}}
	m.ListModel = gml.NewListModel(m)
	return m
}

func (m *Model) RowCount() int {
	return len(m.data)
}

func (m *Model) Data(row int) interface{} {
	return "Test: " + m.data[row]
}

func (m *Model) Append(s string) {
	m.ListModel.Insert(len(m.data), 1, func() {
		m.data = append(m.data, s)
	})
}

func (m *Model) MoveItem(srcRow, dstRow int) {
	m.ListModel.Move(srcRow, 1, dstRow, func() {
		m.data[srcRow], m.data[dstRow] = m.data[dstRow], m.data[srcRow]
	})
}

func (m *Model) UpdateItem(row int, s string) {
	m.ListModel.Reload(row, 1, func() {
		m.data[row] = s
	})
}

func (m *Model) Pop() (s string) {
	m.ListModel.Remove(len(m.data)-1, 1, func() {
		s = m.data[len(m.data)-1]
		m.data = m.data[:len(m.data)-1]
	})
	return
}

func main() {
	app, err := gml.NewApp()
	if err != nil {
		log.Fatalln(err)
	}

	app.SetApplicationName("AppName")
	app.SetOrganizationName("Desertbit")

	b := &Bridge{}
	b.GMLInit()
	err = app.SetContextProperty("bridge", b)
	if err != nil {
		log.Fatalln(err)
	}

	model := newModel()
	err = app.SetContextProperty("modl", model)
	if err != nil {
		log.Fatalln(err)
	}

	time.AfterFunc(time.Second, func() {
		model.Append("hello")
	})
	time.AfterFunc(time.Second*2, func() {
		println(model.Pop())
	})
	time.AfterFunc(time.Second*3, func() {
		model.MoveItem(2, 0)
	})
	time.AfterFunc(time.Second*4, func() {
		model.UpdateItem(1, "REVOLUTION")
	})

	err = app.AddImageProvider(
		"imgprov",
		gml.NewImageProvider(
			gml.KeepAspectRatio,
			gml.FastTransformation,
			func(id string, img *gml.Image) error {
				data, _ := ioutil.ReadFile("/tmp/a.png")
				return img.LoadFromData(data)
			},
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	err = app.Load("qrc:/qml/main.qml")
	if err != nil {
		log.Fatalln(err)
	}

	gml.Exec(app)
}
