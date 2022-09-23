package main

import (
	"github.com/theirish81/gowalker"
	"os"
	"testing"
)

func TestTransformTemplate(t *testing.T) {
	if res, subs, _ := TransformTemplate("foobar"); res != "foobar" && len(subs) != 0 {
		t.Error("simple template is not preserved")
	}
	if _, _, err := TransformTemplate("file://bananas/foo.templ"); err == nil {
		t.Error("access to non existent template should return an error")
	}
	_ = os.Mkdir("test_data", os.ModePerm)
	defer os.RemoveAll("test_data")
	_ = os.WriteFile("test_data/t1.templ", []byte("f1"), os.ModePerm)
	_ = os.WriteFile("test_data/t2.templ", []byte("f2"), os.ModePerm)

	if res, subs, _ := TransformTemplate("file://test_data/t1.templ"); res != "f1" || len(subs) != 1 || subs["t2"] != "f2" {
		t.Error("loading templates from file did not work")
	}
}

func TestRender(t *testing.T) {
	templ := "foo ${data.render(t2)}"
	subs := gowalker.NewSubTemplates()
	subs["t2"] = "${.}"
	data := []byte(`{"data":"bar"}`)
	if res, _ := Render(templ, subs, false, false, false, data); string(res) != "foo bar" {
		t.Error("render did not work")
	}
	templ = `{"d1":"${data}"}`
	if res, _ := Render(templ, subs, true, false, false, data); string(res) != "{\"d1\":\"foo\"}" {

	}
	if res, _ := Render(templ, subs, false, true, false, data); string(res) != "d1: bar\n" {
		t.Error("beautifyYAML did not work")
	}
}
