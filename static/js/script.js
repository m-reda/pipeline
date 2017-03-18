var relations = {
	inputs: {},
	outputs: {},
	links: []
};
var pipeline =
{
	start:   {
		name: "Start",
		x: 10, y: 150,
		inputs: {},
		outputs: {
			output1: {name: 'Run', dist: {unit: 'example1', input: 'input1'}},
		}
	},

	example1:   {
		name: "Example 1",
		x: 300, y: 150,
		inputs: {
			input1: 'Input 1'
		},
		outputs: {
			output1: {name: 'Output 1', dist: {unit: 'example2', input: 'input1'}},
			output2: {name: 'Output 2', dist: {unit: 'example3', input: 'input2'}}
		}
	},
	example2:   {
		name: "Example 2",
		x: 700, y: 100,
		inputs: {
			input1: 'Input 1',
			input2: 'Input 2'
		},
		outputs: {
			output1: {name: "Output 1"},
			output2: {name: "Output 2", dist: {unit: 'example3', input: 'input1'}}
		}
	},
	example3:   {
		name: "Example 3",
		x: 950, y: 300,
		inputs: {
			input1: 'Input 1',
			input2: 'Input 2'
		},
		outputs: {
			output1: {name: "Output 1"},
			output2: {name: "Output 2"}
		}
	}
};


window.addEventListener("load",function(e)
{
	for(var unitID in pipeline)
	{
		var unit = pipeline[unitID];
		var node = new NEditor.Node(unit.name)
			node.setPosition(unit.x, unit.y);

		// inputs
		for(var inputID in unit.inputs) {
			relations.inputs[unitID + '_' + inputID] = node.addInput(unit.inputs[inputID])
		}

		// outputs
		for(var outputID in unit.outputs)
		{
			var output = unit.outputs[outputID];
			relations.outputs[unitID + '_' + outputID] = node.addOutput(output.name);

			if(output.dist)
				relations.links.push([unitID + '_' + outputID, output.dist.unit + '_' + output.dist.input]);
		}

	}

	for(var i in relations.links)
	{
		var link = relations.links[i];
		relations.outputs[link[0]].connectTo(relations.inputs[link[1]])
	}

});

function l(v) {
	console.log(v)
}