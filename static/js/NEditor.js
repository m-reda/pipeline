//Copyright 2016 Sketchpunk Labs

//###########################################################################
//Main Static Object
//###########################################################################
var NEditor = {};
NEditor.dragMode = 0;
NEditor.dragItem = null;    //reference to the dragging item
NEditor.startPos = null;    //Used for starting position of dragging lines
NEditor.offsetX = 0;        //OffsetX for dragging nodes
NEditor.offsetY = 0;        //OffsetY for dragging nodes
NEditor.svg = null;         //SVG where the line paths are drawn.

NEditor.pathColor = "#999999";
NEditor.pathColorA = "#86d530";
NEditor.pathWidth = 2;
NEditor.pathDashArray = "20,5,5,5,5,5";

NEditor.init = function(id){
	NEditor.svg = document.getElementById(id);
	NEditor.svg.ns = NEditor.svg.namespaceURI;
};

/*--------------------------------------------------------
Global Function */

//Trail up the parent nodes to get the X,Y position of an element
NEditor.getOffset = function(elm){
	var pos = {x:0,y:0};
	while(elm){
	  pos.x += elm.offsetLeft;
	  pos.y += elm.offsetTop;
	  elm = elm.offsetParent;
	}
	return pos;
};

//Gets the position of one of the connection points
NEditor.getConnPos = function(elm){
	var pos = NEditor.getOffset(elm);
	pos.x += (elm.offsetWidth / 2) + 1.5; //Add some offset so its centers on the element
	pos.y += (elm.offsetHeight / 2) + 0.5;
	return pos;
};

//Used to reset the svg path between two nodes
NEditor.updateConnPath = function(o){
	var pos1 = o.output.getPos(),
		pos2 = o.input.getPos();
	NEditor.setQCurveD(o.path,pos1.x,pos1.y,pos2.x,pos2.y);
};

//Creates an Quadratic Curve path in SVG
NEditor.createQCurve = function (x1, y1, x2, y2) {
	var elm = document.createElementNS(NEditor.svg.ns,"path");
	elm.setAttribute("fill", "none");
	elm.setAttribute("stroke", NEditor.pathColor);
	elm.setAttribute("stroke-width", NEditor.pathWidth);
	elm.setAttribute("stroke-dasharray", NEditor.pathDashArray);

	NEditor.setQCurveD(elm,x1,y1,x2,y2);
	return elm;
}

//This is seperated from the create so it can be reused as a way to update an existing path without duplicating code.
NEditor.setQCurveD = function(elm,x1,y1,x2,y2){
	var dif = Math.abs(x1-x2) / 1.5,
		str = "M" + x1 + "," + y1 + " C" +	//MoveTo
			(x1 + dif) + "," + y1 + " " +	//First Control Point
			(x2 - dif) + "," + y2 + " " +	//Second Control Point
			(x2) + "," + y2;				//End Point

	elm.setAttribute('d', str);
};

NEditor.setCurveColor = function(elm,isActive){ elm.setAttribute('stroke', (isActive)? NEditor.pathColorA : NEditor.pathColor); };


/*--------------------------------------------------------
Dragging Nodes */
NEditor.beginNodeDrag = function(n, x, y){
	if(NEditor.dragMode != 0) return;

	NEditor.dragMode = 1;
	NEditor.dragItem = n;
	this.offsetX = n.offsetLeft - x;
	this.offsetY = n.offsetTop - y;

	window.addEventListener("mousemove", NEditor.onNodeDragMouseMove);
	window.addEventListener("mouseup",   NEditor.onNodeDragMouseUp);
};

NEditor.onNodeDragMouseUp = function(e){
	e.stopPropagation(); e.preventDefault();
	NEditor.dragItem = null;
	NEditor.dragMode = 0;

	window.removeEventListener("mousemove", NEditor.onNodeDragMouseMove);
	window.removeEventListener("mouseup",   NEditor.onNodeDragMouseUp);
};

NEditor.onNodeDragMouseMove = function(e){
	e.stopPropagation(); e.preventDefault();

	if(NEditor.dragItem){
		var x = e.pageX + NEditor.offsetX,
			y = e.pageY + NEditor.offsetY;

		NEditor.dragItem.style.left = x + "px";
		NEditor.dragItem.style.top  = y + "px";
		NEditor.dragItem.ref.updatePaths();
		NEditor.dragItem.ref.onDrag(x, y)
	}
};

/*--------------------------------------------------------
Dragging Paths */
NEditor.beginConnDrag = function(path){
	if(NEditor.dragMode != 0) return;

	NEditor.dragMode = 2;
	NEditor.dragItem = path;
	NEditor.startPos = path.output.getPos();

	NEditor.setCurveColor(path.path,false);
	window.addEventListener("click",NEditor.onConnDragClick);
	window.addEventListener("mousemove",NEditor.onConnDragMouseMove);
};

NEditor.endConnDrag = function(){
	NEditor.dragMode = 0;
	NEditor.dragItem = null;

	window.removeEventListener("click",NEditor.onConnDragClick);
	window.removeEventListener("mousemove",NEditor.onConnDragMouseMove);
};

NEditor.onConnDragClick = function(e){
	e.stopPropagation(); e.preventDefault();
	NEditor.dragItem.output.removePath(NEditor.dragItem);
	NEditor.endConnDrag();
};

NEditor.onConnDragMouseMove = function(e){
	e.stopPropagation(); e.preventDefault();
	if(NEditor.dragItem) NEditor.setQCurveD(NEditor.dragItem.path,NEditor.startPos.x,NEditor.startPos.y,e.pageX,e.pageY);
};

/*--------------------------------------------------------
Connection Event Handling */
NEditor.onOutputClick = function(e){
	e.stopPropagation(); e.preventDefault();
	var path = e.target.parentNode.ref.addPath();

	NEditor.beginConnDrag(path);
};

NEditor.onInputClick = function(e){
	e.stopPropagation(); e.preventDefault();
	var o = this.parentNode.ref;

	// trigger on remove
	if(o.OutputConn != null)
		o.OutputConn.output.onRemove(o.OutputConn.output.paths.indexOf(o.OutputConn), o.OutputConn.output);

	switch(NEditor.dragMode){
		case 2: //Path Drag
			o.applyPath(NEditor.dragItem);
			NEditor.dragItem.output.onAdd(); // trigger on remove
			NEditor.endConnDrag();
			break;

		case 0: //Not in drag mode
		  var path = o.clearPath();
		  if(path != null) NEditor.beginConnDrag(path);
		  break;
	}

};


//###########################################################################
// Connector Object
//###########################################################################

//Connector UI Object. Ideally this should be an abstract class as a base for an output and input class, but save time
//I wrote this object to handle both types. Its a bit hokey but if it becomes a problem I'll rewrite it in a better OOP way.
NEditor.Connector = function(pElm, isInput, name, id, node_id){
	this.name   = name;
	this.id = id;
	this.node_id = node_id;
	this.root   = document.createElement("li");
	this.dot    = document.createElement("i");
	this.label  = document.createElement("span");

	//Input/Output Specific values
	if(isInput) this.OutputConn = null;		//Input can only handle a single connection.
	else this.paths = [];    				//Outputs can connect to as many inputs is needed

	//Create Elements
	pElm.appendChild(this.root);
	this.root.appendChild(this.dot);
	this.root.appendChild(this.label);

	//Define the Elements
	this.root.className = (isInput)?"Input":"Output";
	this.root.ref = this;
	this.label.innerHTML = name;
	this.dot.innerHTML = "&nbsp;";

	this.dot.addEventListener("click", (isInput) ? NEditor.onInputClick : NEditor.onOutputClick );
};

/*--------------------------------------------------------
Common Methods */

//Get the position of the connection ui element
NEditor.Connector.prototype.getPos = function(){ return NEditor.getConnPos(this.dot); }

//Just updates the UI if the connection is currently active
NEditor.Connector.prototype.resetState = function(){
	var isActive = (this.paths && this.paths.length > 0) || (this.OutputConn != null);

	if(isActive) this.root.classList.add("Active");
	else this.root.classList.remove("Active");
};

//Used mostly for dragging nodes, so this allows the paths to be redrawn
NEditor.Connector.prototype.updatePaths = function(){
	if(this.paths && this.paths.length > 0) for(var i=0; i < this.paths.length; i++) NEditor.updateConnPath(this.paths[i]);
	else if( this.OutputConn ) NEditor.updateConnPath(this.OutputConn);
};

NEditor.Connector.prototype.onAdd = function () {};
NEditor.Connector.prototype.onRemove = function () {};


/*--------------------------------------------------------
Output Methods */

//This creates a new path between nodes
NEditor.Connector.prototype.addPath = function(){
	var pos = NEditor.getConnPos(this.dot),
		dat = {
			path: NEditor.createQCurve(pos.x,pos.y,pos.x,pos.y),
			input:null,
			output:this
		};

	NEditor.svg.appendChild(dat.path);
	this.paths.push(dat);
	return dat;
};

//Remove Path
NEditor.Connector.prototype.removePath = function(o){
	var i = this.paths.indexOf(o);

	if(i > -1)
	{
		// this.onRemove(i, o.output);
		NEditor.svg.removeChild(o.path);
		this.paths.splice(i,1);
		this.resetState();
	}
};

NEditor.Connector.prototype.connectTo = function(o){
	if(o.OutputConn === undefined){
		console.log("connectTo - not an input");
		return;
	}

	var conn = this.addPath();
	o.applyPath(conn);
};

/*--------------------------------------------------------
Input Methods */

//Applying a connection from an output
NEditor.Connector.prototype.applyPath = function(o){
	//If a connection exists, disconnect it.
	if(this.OutputConn != null) this.OutputConn.output.removePath(this.OutputConn);

	//If moving a connection to here, tell previous input to clear itself.
	if(o.input != null) o.input.clearPath();

	o.input = this;			//Saving this connection as the input reference
	this.OutputConn = o;	//Saving the path reference to this object
	this.resetState();		//Update the state on both sides of the connection, TODO some kind of event handling scheme would work better maybe
	o.output.resetState();

	NEditor.updateConnPath(o);
	NEditor.setCurveColor(o.path,true);

	// o.output.onChange()
};

//clearing the connection from an output
NEditor.Connector.prototype.clearPath = function(){
	if(this.OutputConn != null){
		var tmp = this.OutputConn;
		tmp.input = null;

		this.OutputConn = null;
		this.resetState();
		return tmp;
	}
};


//###########################################################################
// Node Object
//###########################################################################
NEditor.Node = function(sTitle, id){
	this.Title = sTitle;
	this.Inputs = [];
	this.Outputs = [];

	//.........................
	this.eRoot = document.createElement("div");
	document.body.appendChild(this.eRoot);
	this.eRoot.className = "NodeContainer";
	this.eRoot.ref = this;
	this.id = id;

	if(id != undefined)
		this.eRoot.id = 'node_' + id;

	//.........................
	this.eHeader = document.createElement("header");
	this.eRoot.appendChild(this.eHeader);
	this.eHeader.innerHTML = this.Title;
	this.eHeader.addEventListener("mousedown", this.onHeaderDown);

	//.........................
	this.eList = document.createElement("ul");
	this.eRoot.appendChild(this.eList);
};


NEditor.Node.prototype.addInput = function(name, id, node_id){
	var o = new NEditor.Connector(this.eList, true, name, id, node_id) ;
	this.Inputs.push(o);
	return o;
};

NEditor.Node.prototype.addOutput = function(name, id, node_id){
	var o = new NEditor.Connector(this.eList, false, name, id, node_id);
	this.Outputs.push(o);
	return o;
};

NEditor.Node.prototype.getInputPos = function(i){ return NEditor.getConnPos(this.Inputs[i].dot); }
NEditor.Node.prototype.getOutputPos = function(i){ return NEditor.getConnPos(this.Outputs[i].dot); }

NEditor.Node.prototype.updatePaths = function(){
	var i;
	for(i=0; i < this.Inputs.length; i++) this.Inputs[i].updatePaths();
	for(i=0; i < this.Outputs.length; i++) this.Outputs[i].updatePaths();
};

//Handle the start node dragging functionality
NEditor.Node.prototype.onHeaderDown = function(e){
	e.stopPropagation();
	NEditor.beginNodeDrag(e.target.parentNode,e.pageX,e.pageY);
};

NEditor.Node.prototype.setPosition = function(x,y){
	this.eRoot.style.left = x + "px";
	this.eRoot.style.top = y + "px";
};

NEditor.Node.prototype.setWidth = function(w){ this.eRoot.style.width = w+"px"; };

NEditor.Node.prototype.onDrag = function () {};