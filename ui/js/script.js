var pipeline = {ID:'',Name:'',Tasks:{'':{Name:'',X:0,Y:0,Command:'',Inputs:{},Outputs:{'':{Name:'',Destination:[{Task:'',Input:''}]}}}}};
var toolbox, units = {}, linker;
var relations = {
	inputs: {},
	outputs: {},
	links: []
};

$(document).ready(function () {
	linker = $('#linker').linker();
	toolbox = $('#toolbox');

	$.getJSON('http://localhost:3000/api/pipelines/1', function(data) {
		pipeline = data;
		drawNodes();
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
});

function drawNodes()
{
	$.each(pipeline.Tasks, function (taskID, task) {
		drawNode(taskID, task);
	});

	$.each(relations.links, function (i, link) {
		relations.outputs[ link[0] ].connect(relations.inputs[ link[1] ], true)
	});
}

function drawNode(taskID, task)
{
	var node = linker.node({id: taskID, name: task.Name, x: task.X, y: task.Y});
		node.onDrag = function (x, y) {
			var task = pipeline.Tasks[this.id];
			task.X = x;
			task.Y = y;
		};

	// inputs
	$.each(task.Inputs, function (id, name) {
		relations.inputs[taskID + '_' + id] = node.input(id, name);
	});

	// outputs
	$.each(task.Outputs, function (outputID, output)
	{
		var o = node.output(outputID, output.Name);

		relations.outputs[ taskID+'_'+outputID ] = o;

		$.each(output.Destination, function (destID, dest) {
			relations.links.push([
				taskID+'_'+outputID,
				dest.Task+'_'+dest.Input
			]);
		});

		// update the new input
		o.onConnect = function (input) {
			pipeline.Tasks[ this.node.id ].Outputs[ this.id ].Destination.push({Task: input.node.id, input: input.id});
		};

		// remove the old input
		o.onRemove = function (index) {
			l(pipeline.Tasks[ this.node.id ].Outputs[ this.id ].Destination, index)
			pipeline.Tasks[ this.node.id ].Outputs[ this.id ].Destination.splice(index, 1);
			l(pipeline.Tasks[ this.node.id ].Outputs[ this.id ].Destination)
		};
	});
}