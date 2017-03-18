var relations = {
	inputs: {},
	outputs: {},
	links: []
};
var pipeline = {
	ID: '-', Name: '-',
	Tasks: {
		"-": {
			"Name": '-', "X": 10, "Y": 40,
			"Command": '-',
			"Inputs": {},
			"Outputs": {
				"-": {
					"Name": '-',
					"Destination": [{"Task": '-', "Input": '-'}]
				}
			}
		}
    }
};

$(document).ready(function () {
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
	NEditor.init('board');

	for(var taskID in pipeline.Tasks)
	{
		var task = pipeline.Tasks[taskID];
		var node = new NEditor.Node(task.Name, taskID)
			node.setPosition(task.X, task.Y);
			node.onDrag = function (x, y) {
				var task = pipeline.Tasks[this.id];
				task.X = x;
				task.Y = y;
			};

		// inputs
		for(var inputID in task.Inputs) {
			relations.inputs[taskID + '_' + inputID] = node.addInput(task.Inputs[inputID], inputID, taskID)
		}

		// outputs
		for(var outputID in task.Outputs)
		{
			var output = task.Outputs[outputID];
			var o = node.addOutput(output.Name, outputID, taskID);

			// remove the old input
			o.onRemove = function (index, output)
			{
				var destination = pipeline.Tasks[ output.node_id ].Outputs[ output.id ].Destination;
				destination.splice(index, 1);

				/*var oldDestinations = pipeline.Tasks[old.output.node_id].Outputs[old.output.id].Destination;
				for(var idx in oldDestinations) {
					var d = oldDestinations[idx];
					if(d.Task == old.input.node_id && d.Input == old.input.id)
						oldDestinations.splice(idx, 1);
				}*/
			};

			// update the new input
			o.onAdd = function ()
			{
				var destinations = pipeline.Tasks[ this.node_id ].Outputs[ this.id ].Destination = [];
				for(var i in this.paths) {
					var input = this.paths[i].input;
					destinations.push({Task: input.node_id, input: input.id})
				}
			};

			relations.outputs[ taskID+'_'+outputID ] = o;

			for(var destID in output.Destination) {
				relations.links.push([taskID+'_'+outputID, output.Destination[destID].Task+'_'+output.Destination[destID].Input]);
			}
		}

	}

	for(var i in relations.links)
	{
		var link = relations.links[i];
		relations.outputs[ link[0] ].connectTo(relations.inputs[ link[1] ])
	}

}

function l() {
	console.log(...arguments)
}