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

	if err = g.AddTemplate("login", `<!-- Login Screen, initially hidden -->
<div class="container" id="loginform">
  <form>
  <div class="row">
    <div class="col s6 offset-s3">
      <h3 class="center-align">Login</h3>        
      <div class="row">
        <div class="input-field col s12">
          <input id="l-username" type="text" class="validate" placeholder="User Name">
          <label for="l-username" class="active">User Name</label>
        </div>
      </div>
      <div class="row">
        <div class="input-field col s12">
          <input id="l-passwd" type="password" class="validate" placeholder="Password">
          <label for="l-passwd" class="active">Password</label>
        </div>
      </div>
      <div class="row">
        <div class="input-field col s6">
          <input id="l-remember" type="checkbox">
          <label for="l-remember">Remember Me ?</label>
        </div>
        <div class="input-field col s6">
          <button id="l-loginbtn" class="btn btn-large waves-effect waves-light" name="Login" type="submit">Login</button>
        </div>
      </div>
    </div>
  </div>
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
