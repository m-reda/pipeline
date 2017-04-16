var settingElm,
	settingTask,
	settingOpened = false;

/** Setting Open
 * add inputs overwrites and setting fields
 * fade setting panel in
 * @param taskID
 */
function settingOpen(taskID)
{
	settingOpened = true;
	settingTask= taskID;

	var task = pipeline.Tasks[taskID];

	// reset the panel elements
	$("h1 span", settingElm).text(task.Name);
	$(".inputs", settingElm).html("");
	$(".setting", settingElm).html("");

	// append the inputs overwrite fields
	$.each(task.Inputs, function (id, name) {
		var val = task.Overwrite[id];
		if(!val) {
			val = "";
		}
		$(".inputs", settingElm).append("<h4>" + name + "</h4><input id=\"setting_overwrite_" + id + "\" value=\"" + val + "\">");
	});

	// append the setting fields
	$.each(task.Setting, function (id, s) {
		var val = task.Setting[id]["Value"];
		if(!val){
			val = "";
		}

		var field;

		switch (s.Type) {
			case "checkbox":
				field = "<input type=\"checkbox\" id=\"setting_options_" + id + "\" " + (val ? "checked" : "") + ">";
				break;

			default:
				field = "<input type=\"" + s.Type + "\" id=\"setting_options_" + id + "\" value=\"" + val + "\">";
				break;
		}

		// add the environment variable name
		var env = (s.EnvVar) ? "<span class=\"env_var\">Environment: " + s.EnvVar + "</span>" : "";

		$(".setting", settingElm).append("<h4>" + s.Name + env + "</h4>" + field);
	});

	// show the setting panel
	settingElm.fadeIn();
	$("#overlay").fadeIn();
}

/** Setting Close
 * update the node's overwrite and setting values
 * fade setting panel out
 */
function settingClose()
{
	var task = pipeline.Tasks[settingTask];

	// update the node's overwrite values
	$.each(task.Inputs, function (id) {
		task.Overwrite[id] = $("#setting_overwrite_" + id).val();
	});

	// update the node's setting values
	$.each(task.Setting, function (id) {
		var el = $("#setting_options_" + id),
			val = el.val();

		if(task.Setting[id].Type === "checkbox") {
			val = el.is(":checked") ? "on" : null;
		}

		task.Setting[id]["Value"] = val;
	});

	// hide the setting panel
	settingElm.fadeOut();
	$("#overlay").fadeOut();
}

/** Task Output
 *
 * @param taskID
 * @param outputID
 * @param output
 * @param node
 */
function taskOutput(taskID, outputID, output, node)
{
	var o = node.output(outputID, output.Name);

	// add new output to outputs map
	relations.outputs[ taskID + "_" + outputID ] = o;

	// add new link to the links map
	$.each(output.Destination, function (destID, dest) {
		relations.links.push([taskID + "_" + outputID, dest.Task + "_" + dest.Input]);
	});

	// update the new input
	o.onConnect = function (input) {
		pipeline.Tasks[ this.node.id ].Outputs[ this.id ].Destination.push({Task: input.node.id, input: input.id});
		markAsNotSaved();
	};

	// remove the old input
	o.onRemove = function (index) {
		pipeline.Tasks[ this.node.id ].Outputs[ this.id ].Destination.splice(index, 1);
		markAsNotSaved();
	};
}

/** Task Node
 *
 * @param taskID
 * @param task
 */
function taskNode(taskID, task)
{
	var node = linker.node({id: taskID, name: task.Name, x: task.X, y: task.Y});

	// update the task's node position
	node.onDragFinish = function (x, y) {
		var task = pipeline.Tasks[this.id];
		task.X = x;
		task.Y = y;

		markAsNotSaved();
	};

	// remove the old input
	node.onRemove = function () {
		delete pipeline.Tasks[ this.id ];
		markAsNotSaved();
	};

	// on setting icon clicked
	node.onSetting = function () {
		if(node.id === "start"){
			schedulerOpen();
		} else {
			settingOpen(this.id);
		}
	};

	// inputs
	$.each(task.Inputs, function (id, name) {
		relations.inputs[taskID + "_" + id] = node.input(id, name);
	});

	// add the outputs to the node
	$.each(task.Outputs, function (outputID, output) {
		taskOutput(taskID, outputID, output, node);
	});
}

/** Build Task Node
 *
 */
function buildTasksNodes()
{
	// build all the tasks nodes
	$.each(pipeline.Tasks, function (taskID, task) {
		taskNode(taskID, task);
	});

	// connect all the links
	$.each(relations.links, function (i, link) {
		relations.outputs[ link[0] ].connect(relations.inputs[ link[1] ], true);
	});

	// close setting
	$(".close", settingElm).click(function () {
		settingClose();
		markAsNotSaved();
	});
}

$(document).ready(function () {
	settingElm = $("#setting");
});