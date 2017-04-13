var toolboxElm,
	units = {};

function toolboxUnits(unitsData)
{
	$.each(unitsData, function (i, unit) {
		units[unit.ID] = unit;

		// if group is empty add it to the general group
		if(!unit.Group){
			unit.Group = "general";
		}

		var group = $("#group_" + unit.Group, toolboxElm);

		// check if the group container not exist add new container
		if(group.length === 0){
			group = $("<div id=\"group_" + unit.Group + "\"><h3>" + unit.Group + "</h3></div>").appendTo(toolboxElm);
		}

		// append the new unit to the group container
		group.append("<div class=\"unit item\" data-unit=\"" + unit.ID + "\"><span>" + unit.Name + "</span></div>");
	});
}

/** Add Unit
 * add new unit to the board
 * @param unitID
 */
function addUnit(unitID)
{
	markAsNotSaved();

	var unit   = units[unitID],
		taskID = Date.now().toString(36) + Math.random().toString(36);

	// handle the unit command prefixes
	unit.Command = unit.Command.replace("bin:", "../../../units/bin").replace("unit:", "../../../units/" + unit.ID);

	// add a start and done by default
	unit.Inputs[0] = "Start";
	unit.Outputs[0] = "Done";

	var task = {
		X: 100,
		Y: 200,
		ID: taskID,
		Unit: unit.ID,
		Name: unit.Name,
		Overwrite: {},
		Setting: unit.Setting,
		Command: unit.Command,
		Inputs: unit.Inputs,
		Outputs: {}
	};

	$.each(unit.Outputs, function (id, name) {
		task.Outputs[id] = {Name:name, Destination:[]};
	});

	pipeline.Tasks[taskID] = task;
	taskNode(taskID, task);
}

$(document).ready(function () {
	toolboxElm = $("#toolbox");

	// load the unites
	api("GET", "units", function(data) {
		toolboxUnits(data);
	});

	// add new unit to the board
	toolboxElm.on("click", ".unit", function () {
		addUnit( $(this).data("unit") );
	});

	// open the toolbox
	$("#toolbox_open").click(function () {
		toolboxElm.addClass("opened");
	});
});