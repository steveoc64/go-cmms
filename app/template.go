package main

// This package has been automatically generated with temple.
// Do not edit manually!

import (
	"github.com/go-humble/temple/temple"
)

var (
	GetTemplate     func(name string) (*temple.Template, error)
	GetPartial      func(name string) (*temple.Partial, error)
	GetLayout       func(name string) (*temple.Layout, error)
	MustGetTemplate func(name string) *temple.Template
	MustGetPartial  func(name string) *temple.Partial
	MustGetLayout   func(name string) *temple.Layout
)

func init() {
	var err error
	g := temple.NewGroup()

	if err = g.AddTemplate("gridform", `<div class="container center-align" id="gridform" style="display:inline-block;">
<img src="img/features.png" alt="">
<form class="grid-form">
    <fieldset>
        <legend>Form Section</legend>
        <div data-row-span="2">
            <div data-field-span="1">
                <label>Field 1</label>
                <input type="text">
            </div>
            <div data-field-span="1">
                <label>Field 2</label>
                <input type="text">
            </div>
        </div>
    </fieldset>
</form>
</div>
`); err != nil {
		panic(err)
	}

	GetTemplate = g.GetTemplate
	GetPartial = g.GetPartial
	GetLayout = g.GetLayout
	MustGetTemplate = g.MustGetTemplate
	MustGetPartial = g.MustGetPartial
	MustGetLayout = g.MustGetLayout
}
