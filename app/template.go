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

	if err = g.AddTemplate("admin-dashboard", `<div class="container" id="admin_dashboard">
This is the admin dashboard
</div>
`); err != nil {
		panic(err)
	}

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

	if err = g.AddTemplate("login", `<div class="container" id="loginform">
<form>
  <fieldset>
    <label for="l-username">Name</label>
    <input type="text" placeholder="User Name" id="l-username">
    <label for="l-passwd">Password</label>
    <!-- <input type="password" placeholder="Password" id="l-passwd"> -->
    <input type="text" placeholder="Password" id="l-passwd">
    <div class="remember-me">
      <input type="checkbox" id="l-remember">
      <label class="label-inline" for="l-remember">Remember Me</label>
    </div>
    <input class="button" type="submit" value="Login" id="l-loginbtn">
  </fieldset>
</form>
</div>
`); err != nil {
		panic(err)
	}

	if err = g.AddTemplate("machine-diag", `<svg class="svg-panel-nohover" xmlns="http://www.w3.org/2000/svg" 
    	 width="95%"
    	 viewBox="0 0 {{.Machine.SVGWidth1}} 180">
  <defs>
    <linearGradient id="grad1" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" style="stop-color:#eee;stop-opacity:1" />
      <stop offset="100%" style="stop-color:#869ab1;stop-opacity:1" />
    </linearGradient>
    <linearGradient id="bgrad" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" style="stop-color:#fff;stop-opacity:1" />
      <stop offset="100%" style="stop-color:#9ac;stop-opacity:1" />
    </linearGradient>
    <linearGradient id="bgradh" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" style="stop-color:#fff;stop-opacity:1" />
      <stop offset="100%" style="stop-color:#abd;stop-opacity:1" />
    </linearGradient>
    <filter id="shadow" x="-20%" y="-20%" width="140%" height="140%">
      <feGaussianBlur stdDeviation="2 2" result="shadow"/>
      <feOffset dx="3" dy="3"/>
    </filter>      
    <filter id="shadow1" x="0" y="0" width="200%" height="200%">
      <feOffset result="offOut" in="SourceAlpha" dx="1" dy="1" />
      <feGaussianBlur result="blurOut" in="offOut" stdDeviation="10" />
      <feBlend in="SourceGraphic" in2="blurOut" mode="normal" />
    </filter>
    <radialGradient id="GreenBtn">
        <stop offset="10%" stop-color="#4f2"/>
        <stop offset="95%" stop-color="#2a1"/>
    </radialGradient>
    <radialGradient id="DullGreenBtn">
        <stop offset="10%" stop-color="#898"/>
        <stop offset="95%" stop-color="#888"/>
    </radialGradient>
    <radialGradient id="BlueBtn">
        <stop offset="10%" stop-color="#ff0">
          <animate attributeName="stop-color"
          values="#26a;#08e;#0af;#0af;#08e;#26a"
          dur="0.8s"
          repeatCount="indefinite" />
        </stop>
        <stop offset="95%" stop-color="#26a"/>
    </radialGradient>
    <radialGradient id="YellowBtn">
        <stop offset="10%" stop-color="#ff0">
          <animate attributeName="stop-color"
          values="#da2;#ee0;#ff0;#ff0;#ee0;#da2"
          dur="0.8s"
          repeatCount="indefinite" />
        </stop>
        <stop offset="95%" stop-color="#da2"/>
    </radialGradient>
    <radialGradient id="DullYellowBtn">
        <stop offset="10%" stop-color="#998"/>
        <stop offset="95%" stop-color="#888"/>
    </radialGradient>
    <radialGradient id="RedBtn">   <!--  fx="60%" fy="30%"> -->
        <stop offset="10%" stop-color="#fa0">
          <animate attributeName="stop-color"
          values="#f00;#f80;#fa0;#f80;#f00"
          dur="0.8s"
          repeatCount="indefinite" />
        </stop>
        <stop offset="95%" stop-color="#e00">
          <animate attributeName="stop-color"
          values="#800;#a00;#f00;#a00;#800"
          dur="0.8s"
          repeatCount="indefinite" />
        </stop>
    </radialGradient>
  </defs>

  <!-- Add picture for uncoiler -->
  <g stroke="#114" stroke-width="2" class="fillhover">    
    <title>Raise Event on Uncoiler</title>
    <polygon points="30,90 0,148 60,148" fill="url(#grad1)"/>
    <circle cx="30" cy="90" r="30" fill="{{.Machine.NonToolBg .Machine.Uncoiler}}"/>
  </g>

  <!-- Add buttons for elec, etc -->
  <g stroke="#114" stroke-width="1" class="fillhover">
    <title>Electrical</title>
    <rect x="0" y="0" width="40" rx="3" ry="3" height="40" class="bhover"
          fill="{{.Machine.NonToolBg .Machine.Electrical}}"/>
    <image xlink:href="/img/elec.png" x="1" y="2" height="38px" width="38px"/>
  </g>  
  <g stroke="#114" stroke-width="1" class="fillhover">
    <title>Hydraulic</title>
    <rect x="50" y="0" width="40" rx="3" ry="3" height="40" class="bhover"
          fill="{{.Machine.NonToolBg .Machine.Hydraulic}}"/>
    <image xlink:href="/img/hydraulic.png" x="51" y="2" height="38px" width="38px"/>
  </g>
  <g stroke="#114" stroke-width="1" class="fillhover">
    <title>Lube</title>
    <rect x="100" y="0" width="40" rx="3" ry="3" height="40" class="bhover"
          fill="{{.Machine.NonToolBg .Machine.Lube}}"/>
    <image xlink:href="/img/lube.png" x="101" y="2" height="38px" width="38px"/>
  </g>
  <g stroke="#114" stroke-width="1" class="fillhover">
    <title>Printer</title>
    <rect x="150" y="0" width="40" rx="3" ry="3" height="40" class="bhover"
          fill="{{.Machine.NonToolBg .Machine.Printer}}"/>
    <image xlink:href="/img/printer.png" x="151" y="2" height="38px" width="38px"/>
  </g>
  <g stroke="#114" stroke-width="1" class="fillhover">
    <title>Console</title>
    <rect x="200" y="0" width="40" rx="3" ry="3" height="40" class="bhover"
          fill="{{.Machine.NonToolBg .Machine.Console}}"/>
    <image xlink:href="/img/console.png" x="201" y="2" height="38px" width="38px"/>
  </g>

  <!-- Add main rectangle -->
  <rect x="80" y="100" width="{{.Machine.SVGWidth2}}" 
    height="48" stroke="black" stroke-width="2" fill="url(#grad1)"/>
  <text transform="translate(100 135)" style="font-size: 22;">{{.Machine.Name}}</text>

  <!-- Add the status indicators -->
  <g stroke="black" fill="{{.Machine.SVGStatus}}">
    <circle cx="{{.Machine.SVGX}}" cy="125" r="18"/>
  </g>
  <g stroke="#114" stroke-width="2" class="fillhover">
    <title>Roll Bed</title>

  <rect x="80" y="60" width="30" height="40" fill="{{.Machine.NonToolBg .Machine.Rollbed}}"/>
    <circle cx="95" cy="73" r="6" fill="#ccc"/>
    <circle cx="95" cy="87" r="6" fill="#ccc"/>
  <rect x="115" y="60" width="30" height="40" fill="{{.Machine.NonToolBg .Machine.Rollbed}}"/>
    <circle cx="130" cy="73" r="6" fill="#ccc"/>
    <circle cx="130" cy="87" r="6" fill="#ccc"/>
  <rect x="150" y="60" width="30" height="40" fill="{{.Machine.NonToolBg .Machine.Rollbed}}"/>
    <circle cx="165" cy="73" r="6" fill="#ccc"/>
    <circle cx="165" cy="87" r="6" fill="#ccc"/>
  <rect x="185" y="60" width="30" height="40" fill="{{.Machine.NonToolBg .Machine.Rollbed}}"/>
    <circle cx="200" cy="73" r="6" fill="#ccc"/>
    <circle cx="200" cy="87" r="6" fill="#ccc"/>
  </g>

  <!-- Now draw all the tools     -->
  {{$compID := .CompID}}
  {{range $index,$comp := .Machine.Components}}
  <svg x="{{$comp.SVGX $index}}" 
       class="tooltip tooltip--ne"
       ng-click="Machines.raiseIssue(row,comp,comp.ID,'tool')">
       <title>{{$comp.Name}}</title>
    <a>
    <rect x="16" y="0" width="15" height="20" stroke="black" stroke-width="1" fill="#ddd"/>
    <rect y="20" width="45" rx="10" ry="10" height="80" stroke="black" stroke-width="2" 
          fill="{{$comp.SVGFill2 $compID}}"
          class="hoverme"/>
    <text x="5" y="50">{{$comp.SVGName $index}}</text>
    </a>
  </svg>
  {{end}}

</svg>`); err != nil {
		panic(err)
	}

	if err = g.AddTemplate("machines", `<div class="container" id="machines">
This is the Machines List
</div>
`); err != nil {
		panic(err)
	}

	if err = g.AddTemplate("raise-comp-issue", `<div class="md-content">
	<h3>Raise New Issue</h3>	
	<div id="issue-machine-diag" class="row"></div>
	<div>
	<form>
		<fieldset>
	    <label for="desc">Description of Problem</label>
	    {{if .IsTool}}
	    	<textarea id="desc">Problem with {{.Component.Name}} tool on {{.Machine.Name}} machine.</textarea>
	    {{else}}
	    	<textarea id="desc">Problem with {{.NonTool}} on {{.Machine.Name}} machine.</textarea>
	    {{end}}
	    <label for="photo">Upload Photo</label>
	    <input id="photo" name="photo" type="file">
		</fieldset>
	</form>
	</div>

	<div class="row">
		<button class="column button-outline md-close">Cancel</button>
		<button class="column button-primary md-save">Raise Event</button>
	</div>
</div>
`); err != nil {
		panic(err)
	}

	if err = g.AddTemplate("sitelist", `<div class="container" id="sitelist">
This is the Site List
</div>
`); err != nil {
		panic(err)
	}

	if err = g.AddTemplate("sitemachines", `<div class="fluid skipheader">
	<div class="row">
    {{if .MultiSite}}
		<div class="column column-30">	
      <svg xmlns="http://www.w3.org/2000/svg" 
        height="190" width="305">
        <defs>
          <radialGradient id="GreenBtn">
              <stop offset="10%" stop-color="#4f2"/>
              <stop offset="95%" stop-color="#2a1"/>
          </radialGradient>
          <radialGradient id="YellowBtn">
              <stop offset="10%" stop-color="#ff0">
                <animate attributeName="stop-color"
                values="#da2;#ee0;#ff0;#ff0;#ee0;#da2"
                dur="0.8s"
                repeatCount="indefinite" />
              </stop>
              <stop offset="95%" stop-color="#da2"/>
          </radialGradient>
          <radialGradient id="RedBtn">   <!--  fx="60%" fy="30%"> -->
              <stop offset="10%" stop-color="#fa0">
                <animate attributeName="stop-color"
                values="#f00;#f80;#fa0;#f80;#f00"
                dur="0.8s"
                repeatCount="indefinite" />
              </stop>
              <stop offset="95%" stop-color="#e00">
                <animate attributeName="stop-color"
                values="#800;#a00;#f00;#a00;#800"
                dur="0.8s"
                repeatCount="indefinite" />
              </stop>
          </radialGradient>
        </defs>

        <image xlink:href="/img/aust.png" x="1" y="1" height="182px" width="201px" id="austmap"/>
        <text x="60" y="160">Edinburgh</text>
        <g stroke="black" fill="url(#{{.Status.EButton}})">
          <circle cx="125" cy="130" r="7"/>
        </g>
        <text x="202" y="102">Chinderah</text>
        <g stroke="black" fill="url(#{{.Status.CButton}})">
          <circle cx="190" cy="100" r="7"/>
        </g>
        <text x="200" y="122">Tomago</text>
        <g stroke="black" fill="url(#{{.Status.TButton}})">
          <circle cx="190" cy="115" r="7"/>
        </g>
        <text x="190" y="142">Minto</text>
        <g stroke="black" fill="url(#{{.Status.MButton}})">
          <circle cx="180" cy="130" r="7"/>
        </g>
      </svg>
		</div>
    {{end}}
		<div class="column column-70">
	    <h1>{{.Site.Name}}</h1>		
		</div>
	</div>  <!-- Map of Australia -->

<!-- Give it a menu on the right -->
<nav class="cbp-spmenu cbp-spmenu-vertical cbp-spmenu-right" id="machine-menu">
</nav>
<!-- End of menu -->


<!-- Raise issue modal dialog -->
<div class="md-modal md-effect-12" id="raise-comp-issue"></div>
<div class="md-overlay"></div>

<!-- Grid of machines -->
<div class="row row-wrap" style="flex-wrap: wrap">

	{{range .Machines}}
	<div class="column span-{{.Span}}" id="machine-div-{{.ID}}" machine-id="{{.ID}}">
  <svg class="svg-panel" xmlns="http://www.w3.org/2000/svg" 
    	 width="95%"
    	 viewBox="0 0 {{.SVGWidth1}} 200">
    <defs>
      <linearGradient id="grad1" x1="0%" y1="0%" x2="100%" y2="100%">
        <stop offset="0%" style="stop-color:#eee;stop-opacity:1" />
        <stop offset="100%" style="stop-color:#869ab1;stop-opacity:1" />
      </linearGradient>
      <linearGradient id="bgrad" x1="0%" y1="0%" x2="100%" y2="100%">
        <stop offset="0%" style="stop-color:#fff;stop-opacity:1" />
        <stop offset="100%" style="stop-color:#9ac;stop-opacity:1" />
      </linearGradient>
      <linearGradient id="bgradh" x1="0%" y1="0%" x2="100%" y2="100%">
        <stop offset="0%" style="stop-color:#fff;stop-opacity:1" />
        <stop offset="100%" style="stop-color:#abd;stop-opacity:1" />
      </linearGradient>
      <filter id="shadow" x="-20%" y="-20%" width="140%" height="140%">
        <feGaussianBlur stdDeviation="2 2" result="shadow"/>
        <feOffset dx="3" dy="3"/>
      </filter>      
      <filter id="shadow1" x="0" y="0" width="200%" height="200%">
        <feOffset result="offOut" in="SourceAlpha" dx="1" dy="1" />
        <feGaussianBlur result="blurOut" in="offOut" stdDeviation="10" />
        <feBlend in="SourceGraphic" in2="blurOut" mode="normal" />
      </filter>
      <radialGradient id="GreenBtn">
          <stop offset="10%" stop-color="#4f2"/>
          <stop offset="95%" stop-color="#2a1"/>
      </radialGradient>
      <radialGradient id="DullGreenBtn">
          <stop offset="10%" stop-color="#898"/>
          <stop offset="95%" stop-color="#888"/>
      </radialGradient>
      <radialGradient id="YellowBtn">
          <stop offset="10%" stop-color="#ff0">
            <animate attributeName="stop-color"
            values="#da2;#ee0;#ff0;#ff0;#ee0;#da2"
            dur="0.8s"
            repeatCount="indefinite" />
          </stop>
          <stop offset="95%" stop-color="#da2"/>
      </radialGradient>
      <radialGradient id="DullYellowBtn">
          <stop offset="10%" stop-color="#998"/>
          <stop offset="95%" stop-color="#888"/>
      </radialGradient>
      <radialGradient id="RedBtn">   <!--  fx="60%" fy="30%"> -->
          <stop offset="10%" stop-color="#fa0">
            <animate attributeName="stop-color"
            values="#f00;#f80;#fa0;#f80;#f00"
            dur="0.8s"
            repeatCount="indefinite" />
          </stop>
          <stop offset="95%" stop-color="#e00">
            <animate attributeName="stop-color"
            values="#800;#a00;#f00;#a00;#800"
            dur="0.8s"
            repeatCount="indefinite" />
          </stop>
      </radialGradient>
    </defs>

    <!-- Add picture for uncoiler -->
    <g stroke="#114" stroke-width="2" class="fillhover">    
      <title>Raise Event on Uncoiler</title>
      <polygon points="30,90 0,148 60,148" fill="url(#grad1)"/>
      <circle cx="30" cy="90" r="30" fill="{{.NonToolBg .Uncoiler}}"/>
    </g>

    <!-- Add buttons for elec, etc -->
    <g stroke="#114" stroke-width="1" class="fillhover">
      <title>Electrical</title>
      <rect x="0" y="0" width="40" rx="3" ry="3" height="40" class="bhover"
            fill="{{.NonToolBg .Electrical}}"/>
      <image xlink:href="/img/elec.png" x="1" y="2" height="38px" width="38px"/>
    </g>  
    <g stroke="#114" stroke-width="1" class="fillhover">
      <title>Hydraulic</title>
      <rect x="50" y="0" width="40" rx="3" ry="3" height="40" class="bhover"
            fill="{{.NonToolBg .Hydraulic}}"/>
      <image xlink:href="/img/hydraulic.png" x="51" y="2" height="38px" width="38px"/>
    </g>
    <g stroke="#114" stroke-width="1" class="fillhover">
      <title>Lube</title>
      <rect x="100" y="0" width="40" rx="3" ry="3" height="40" class="bhover"
            fill="{{.NonToolBg .Lube}}"/>
      <image xlink:href="/img/lube.png" x="101" y="2" height="38px" width="38px"/>
    </g>
    <g stroke="#114" stroke-width="1" class="fillhover">
      <title>Printer</title>
      <rect x="150" y="0" width="40" rx="3" ry="3" height="40" class="bhover"
            fill="{{.NonToolBg .Printer}}"/>
      <image xlink:href="/img/printer.png" x="151" y="2" height="38px" width="38px"/>
    </g>
    <g stroke="#114" stroke-width="1" class="fillhover">
      <title>Console</title>
      <rect x="200" y="0" width="40" rx="3" ry="3" height="40" class="bhover"
            fill="{{.NonToolBg .Console}}"/>
      <image xlink:href="/img/console.png" x="201" y="2" height="38px" width="38px"/>
    </g>
    
    <!-- Add main rectangle -->
    <rect x="80" y="100" width="{{.SVGWidth2}}" 
      height="48" stroke="black" stroke-width="2" fill="url(#grad1)"/>
    <text transform="translate(100 135)" style="font-size: 22;">{{.Name}}</text>

    <!-- Add the status indicators -->
    <g stroke="black" fill="{{.SVGStatus}}">
      <circle cx="{{.SVGX}}" cy="125" r="18"/>
    </g>
    <g stroke="#114" stroke-width="2" class="fillhover">
      <title>Roll Bed</title>
  
    <rect x="80" y="60" width="30" height="40" fill="{{.NonToolBg .Rollbed}}"/>
      <circle cx="95" cy="73" r="6" fill="#ccc"/>
      <circle cx="95" cy="87" r="6" fill="#ccc"/>
    <rect x="115" y="60" width="30" height="40" fill="{{.NonToolBg .Rollbed}}"/>
      <circle cx="130" cy="73" r="6" fill="#ccc"/>
      <circle cx="130" cy="87" r="6" fill="#ccc"/>
    <rect x="150" y="60" width="30" height="40" fill="{{.NonToolBg .Rollbed}}"/>
      <circle cx="165" cy="73" r="6" fill="#ccc"/>
      <circle cx="165" cy="87" r="6" fill="#ccc"/>
    <rect x="185" y="60" width="30" height="40" fill="{{.NonToolBg .Rollbed}}"/>
      <circle cx="200" cy="73" r="6" fill="#ccc"/>
      <circle cx="200" cy="87" r="6" fill="#ccc"/>
    </g>

    <!-- Now draw all the tools     -->
    {{range $index,$comp := .Components}}
    <svg x="{{$comp.SVGX $index}}" 
         class="tooltip tooltip--ne"
         ng-click="Machines.raiseIssue(row,comp,comp.ID,'tool')">
         <title>{{$comp.Name}}</title>
      <a>
      <rect x="16" y="0" width="15" height="20" stroke="black" stroke-width="1" fill="#ddd"/>
      <rect y="20" width="45" rx="10" ry="10" height="80" stroke="black" stroke-width="2" 
            fill="{{$comp.SVGFill}}"
            class="hoverme"/>
      <text x="5" y="50">{{$comp.SVGName $index}}</text>
      </a>
    </svg>
    {{end}}

    </svg>

		</div>		
	{{end}}
	</div>
</div>

`); err != nil {
		panic(err)
	}

	if err = g.AddTemplate("sitemap", `<div class="container" id="sitemap">
	<div class="row">
		<div class="column" style="justify-content: space-around;">
			<!-- <img src="img/aust.png" alt="Australia" width="402px" height="364px"> -->
		  <svg class="svg-map"
		       viewBox="0 0 490 380"
		       xmlns="http://www.w3.org/2000/svg">
		    <defs>
		      <radialGradient id="GreenBtn">
		          <stop offset="10%" stop-color="#4f2"/>
		          <stop offset="95%" stop-color="#2a1"/>
		      </radialGradient>    
		      <radialGradient id="YellowBtn">
		          <stop offset="10%" stop-color="#ff0">
		            <animate attributeName="stop-color"
		            values="#da2;#ee0;#ff0;#ff0;#ee0;#da2"
		            dur="0.8s"
		            repeatCount="indefinite" />
		          </stop>
		          <stop offset="95%" stop-color="#da2"/>
		      </radialGradient>
		      <radialGradient id="RedBtn">   <!--  fx="60%" fy="30%"> -->
		          <stop offset="10%" stop-color="#fa0">
		            <animate attributeName="stop-color"
		            values="#f00;#f80;#fa0;#f80;#f00"
		            dur="0.8s"
		            repeatCount="indefinite" />
		          </stop>
		          <stop offset="95%" stop-color="#e00">
		            <animate attributeName="stop-color"
		            values="#800;#a00;#f00;#a00;#800"
		            dur="0.8s"
		            repeatCount="indefinite" />
		          </stop>
		      </radialGradient>
		    </defs>

		    <image xlink:href="/img/aust.png" x="1" y="1" height="364px" width="402px"/>
		    <text x="180" y="290">Edinburgh</text>
 		    <g stroke="black" fill="url(#{{.Status.EButton}})">
		      <circle cx="260" cy="250" r="12"/>
		    </g>
		    <text x="410" y="205">Chinderah</text>
		    <g stroke="black" fill="url(#{{.Status.CButton}})">
		      <circle cx="390" cy="200" r="12"/>
		    </g>
		    <text x="400" y="235">Tomago</text>
		    <g stroke="black" fill="url(#{{.Status.TButton}})">
		      <circle cx="380" cy="230" r="12"/>
		    </g>
	      <text x="380" y="265">Minto</text>
		    <g stroke="black" fill="url(#{{.Status.MButton}})">
		      <circle cx="360" cy="260" r="12"/>
		    </g>
		  </defs>
		</svg>

		</div>
	</div>
	<div class="row row-wrap" style="flex-wrap: wrap">
		{{range .Sites}}
		<div class="column row-">
			<input type="button" value="{{.Name}}" id="{{.ID}}">						
		</div>
		{{end}}
	</div>

<div class="pricing pricing--tenzin">
	<div class="pricing__item">
		<h3 class="pricing__title">Startup</h3>
		<div class="pricing__price"><span class="pricing__currency">$</span>9.90</div>
		<p class="pricing__sentence">Small business solution</p>
		<ul class="pricing__feature-list">
			<li class="pricing__feature">Unlimited calls</li>
			<li class="pricing__feature">Free hosting</li>
			<li class="pricing__feature">40MB of storage space</li>
		</ul>
		<button class="pricing__action">Choose plan</button>
	</div>
	<div class="pricing__item">
		<h3 class="pricing__title">Standard</h3>
		<div class="pricing__price"><span class="pricing__currency">$</span>29,90</div>
		<p class="pricing__sentence">Medium business solution</p>
		<ul class="pricing__feature-list">
			<li class="pricing__feature">Unlimited calls</li>
			<li class="pricing__feature">Free hosting</li>
			<li class="pricing__feature">10 hours of support</li>
			<li class="pricing__feature">Social media integration</li>
			<li class="pricing__feature">1GB of storage space</li>
		</ul>
		<button class="pricing__action">Choose plan</button>
	</div>
	<div class="pricing__item">
		<h3 class="pricing__title">Professional</h3>
		<div class="pricing__price"><span class="pricing__currency">$</span>59,90</div>
		<p class="pricing__sentence">Gigantic business solution</p>
		<ul class="pricing__feature-list">
			<li class="pricing__feature">Unlimited calls</li>
			<li class="pricing__feature">Free hosting</li>
			<li class="pricing__feature">Unlimited hours of support</li>
			<li class="pricing__feature">Social media integration</li>
			<li class="pricing__feature">Anaylitcs integration</li>
			<li class="pricing__feature">Unlimited storage space</li>
		</ul>
		<button class="pricing__action">Choose plan</button>
	</div>
</div>

</div>

`); err != nil {
		panic(err)
	}

	if err = g.AddTemplate("user-profile", `<div class="md-content">
	<h3>User Profile - {{.Username}}</h3>
	<div>
	<form id="user-profile-form">
	  <fieldset>
	    <label for="nameField">Name</label>
	    <input type="text" value="{{.Name}}" id="nameField" name="Name">

	    <label for="emailField">Email</label>
	    <input type="text" value="{{.Email}}" id="emailField" name="Email">

	    <label for="smsField">SMS</label>
	    <input type="text" value="{{.SMS}}" id="smsField" name="SMS">

	    <label for="roleField">Role</label>
	    <input type="text" value="{{.Role}}" id="roleField" readonly>

	    <label for="pwField">New Password</label>
	    <input type="password" id="pwField" name="p1" placeholder="Leave Blank to remain unchanged">
	    <input type="password" id="pwcField" name="p2" placeholder="Repeat Password to change">
	  </fieldset>
	</form>
	</div>

	<div class="row">
	<button class="column button-outline md-up-close">Cancel</button>
	<button class="column button-primary md-up-save">Save</button>
	</div>
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
