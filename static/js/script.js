var pipeline = {ID:'',Name:'',Tasks:{'':{Name:'',X:0,Y:0,Command:'',Inputs:{},Outputs:{'':{Name:'',Destination:[{Task:'',Input:''}]}}}}};
var toolbox, units = {};
var relations = {
	inputs: {},
	outputs: {},
	links: []
};

$(document).ready(function () {
	toolbox = $('#toolbox');

	$.getJSON('http://localhost:3000/api/pipelines/1', function(data) {
		pipeline = data;
		drawNodes();
	});

	$.getJSON('http://localhost:3000/api/units', function(data) {
		$.each(data, function (i, unit) {
			units[unit.ID] = unit;
			toolbox.append('<div class="unit" data-unit="' + unit.ID + '">' + unit.Name + '</div>')
		});
	});

	$('#save').click(function () {
		var btn = $(this).text("Loading");

		$.ajax({
			method: "PUT",
			data: { pipeline: JSON.stringify(pipeline)},
			url: 'http://localhost:3000/api/pipelines/1',
			dataType: 'json'
		}).done(function( data ) {
			btn.text("Save");
			l(data);
		});
	});

	$('#add').click(function () {
		toolbox.toggle()
	});

	toolbox.on('click', '.unit', function () {
		var unit = units[$(this).data('unit')];
		var taskID = Math.random().toString(36);
		var task = {
			Name: unit.Name,
			Setting: unit.Setting,
			X: 0, Y: 0,
			Command: unit.Command,
			Inputs: unit.Inputs,
			Outputs: {}
		};

		$.each(unit.Outputs, function (id, name) {
			task.Outputs[id] = {Name:name,Destination:[]}
		});

		pipeline.Tasks[taskID] = task;
		drawNode(taskID, task)
	});
});

function drawNodes()
{
	NEditor.init('board');
	$.each(pipeline.Tasks, function (taskID, task) {
		drawNode(taskID, task);
	});

	$.each(relations.links, function (i, link) {
		relations.outputs[ link[0] ].connectTo(relations.inputs[ link[1] ])
	});
}

function drawNode(taskID, task)
{
	var node = new NEditor.Node(task.Name, taskID);
		node.setPosition(task.X, task.Y);
		node.onDrag = function (x, y) {
			var task = pipeline.Tasks[this.id];
			task.X = x;
			task.Y = y;
		};

	// inputs
	$.each(task.Inputs, function (inputID, input) {
		relations.inputs[taskID + '_' + inputID] = node.addInput(input, inputID, taskID)
	});

	// outputs
	$.each(task.Outputs, function (outputID, output)
	{
		var o = node.addOutput(output.Name, outputID, taskID);

		// remove the old input
		o.onRemove = function (index, output) {
			pipeline.Tasks[ output.node_id ].Outputs[ output.id ].Destination.splice(index, 1);
		};

		// update the new input
		o.onAdd = function () {
			var destinations = pipeline.Tasks[ this.node_id ].Outputs[ this.id ].Destination = [];

			$.each(this.paths, function (i, conn) {
				destinations.push({Task: conn.input.node_id, input: conn.input.id})
			});
		};

		relations.outputs[ taskID+'_'+outputID ] = o;

		$.each(output.Destination, function (destID, dest) {
			relations.links.push([
				taskID+'_'+outputID,
				dest.Task+'_'+dest.Input
			]);
		});
	});
}