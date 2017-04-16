var linker,
	schedulerElm,
	isNew = true,
	pipeline = {},
	isNotSaved = false,
	schedulerOpened = false,
	relations = { inputs: {}, outputs: {}, links: [] },
	base = window.location.href.split("?")[0].split("#")[0] + "api/";

/** API
 *
 * @param method
 * @param resource
 * @param callback
 * @param data
 */
function api(method, resource, callback, data)
{
	return $.ajax({
		method: method,
		url: base + resource,
		data: data,
		dataType: "json",
		success: function(data) {
			if(callback){
				callback(data);
			}
		}
	});
}

/** Not Saved
 *
 */
function markAsNotSaved() {
	$("#notSaved").fadeIn();
	isNotSaved = true;
}

/** Saved
 *
 */
function markAsSaved() {
	$("#notSaved").fadeOut();
	isNotSaved = false;
}

/** Reset Board
 *
 */
function resetBoard()
{
	// reset the board
	$(".linker_board", linker).remove();
	$("ul", buildListElm).html("");
	$("#overlay").hide();
	$("#pipelines").removeClass("opened");

	relations = { inputs: {}, outputs: {}, links: [] };
	linker = $("#linker").linker();
	pipeline = {};
}

/** Empty Board
 *
 */
function newBoard()
{
	resetBoard();
	isNew = true;

	$("#pipeline_title").val("New Pipeline");
	buildListElm.hide();

	pipeline = {LastBuild: 0, Builds: [], Schedule: [], Tasks:{"start": {Name: "Start",X:100,Y:100,Inputs:{},Outputs:{"run": {Name: "Run",Destination:[]}}}}};
	buildTasksNodes();

	window.location.hash = "";
}

/** Load Pipeline
 *
 * @param id
 */
function loadPipeline(id)
{
	resetBoard();

	api("GET", "pipelines/" + id, function(data) {
		// if the pipeline is not valid
		if(!data.ID) {
			alert("The pipeline is not valid");
			newBoard();
			return;
		}

		// build the board
		pipeline = data;
		buildTasksNodes();
		window.location.hash = "!" + pipeline.ID;

		$("#pipeline_title").val(pipeline.Name);
		buildListElm.show();

		buildList();

		isNew = false;
	});
}

/** Save Pipeline
 *
 */
function savePipeline()
{
	var btn = $(this).text("Loading");

	pipeline.Name = $("#pipeline_title").val();

	api(isNew ? "POST" : "PUT", "pipelines", function( data ) {
		// if the save fail
		if(!data.id) {
			alert("Pipeline save fail");
			return;
		}

		if(isNew) {
			// assign the new id
			pipeline.ID = data.id;
			window.location.hash = "!" + pipeline.ID;

			// append the new pipeline to the list
			$("#pipelines").append("<div class=\"pipeline item\" data-id=\"" + pipeline.ID + "\"><button class=\"delete material-icons\">delete</button>" + pipeline.Name + "</div>");

			isNew = false;
		}

		markAsSaved();
		btn.text("Save");
		buildListElm.show();
	}, { pipeline: JSON.stringify(pipeline)});
}

/** Delete Pipeline
 *
 */
function deletePipeline(elm)
{
	if(!confirm("Delete this pipeline?")){
		return;
	}

	var id = elm.data("id");
	api("DELETE", "pipelines/" + id, function( data ) {
		if(!data.success){
			alert("Delete the pipeline fail");
		}

		// if the deleted pipeline is in edit mode start new board
		if(id === pipeline.ID){
			newBoard();
		}

		// delete the pipeline list item
		elm.remove();
	});
}

/** Saved Pipeline
 *
 */
function loadSavedPipelines()
{
	api("GET", "pipelines", function(data) {
		$.each(data, function (i, p) {
			$("#pipelines").append("<div class=\"pipeline item\" data-id=\"" + p.ID + "\"><button class=\"delete material-icons\">delete</button>" + p.Name + "</div>");
		});
	});
}

/** Scheduler Init
 *
 */
function schedulerInit() {
	schedulerElm = $("#scheduler");

	// initial the form
	$(".cron-ui", schedulerElm).cron({
		initial: "0 0 * * 0",
		onChange: function() {
			$("button", schedulerElm).data("expression", $(this).cron("value"));
		}
	});

	// add new cron
	$("button", schedulerElm).click(function () {
		var exp = $(this).data("expression");

		pipeline.Schedule.push(exp);
		$("ul", schedulerElm).append("<li data-expression=\"" + exp + "\">" + exp + "<span class=\"remove material-icons\">close</span></li>");
	});

	// remove a cron
	schedulerElm.on("click", ".remove", function () {
		var li = $(this).parent(),
			exp = li.data("expression");

		pipeline.Schedule.splice(pipeline.Schedule.indexOf(exp), 1);
		li.remove();
	});

	// remove a cron
	$(".close", schedulerElm).click(function () {
		schedulerOpened = false;
		markAsNotSaved();

		schedulerElm.fadeOut();
		$("#overlay").fadeOut();
	});
}

/** Scheduler Open
 *
 */
function schedulerOpen()
{
	schedulerOpened = true;

	var ul = $("ul", schedulerElm).html("");

	// append the old schedule
	$.each(pipeline.Schedule, function (i, exp) {
		ul.append("<li data-expression=\"" + exp + "\">" + exp + "<span class=\"remove material-icons\">close</span></li>");
	});

	// the url trigger
	$(".url_trigger").val(pipeline.ID ? base + "pipelines/" + pipeline.ID + "/build" : "");

	// show scheduler panel
	schedulerElm.fadeIn();
	$("#overlay").fadeIn();
}

/** Initialize Board
 *
 */
function initBoard()
{
	linker = $("#linker").linker();

	var hashID = window.location.hash;
	if(hashID.indexOf("#!") === 0){
		loadPipeline(hashID.substring(2));
	} else {
		newBoard();
	}
}

$(document).ready(function () {
	initBoard();
	loadSavedPipelines();

	$("#pipelines").on("click", ".pipeline", function () {
		markAsSaved();

		var id = $(this).data("id");

		// check if it new pipeline or edit old one
		if(id){
			loadPipeline(id);
		} else{
			newBoard();
		}
	})
	// delete pipeline
	.on("click", ".pipeline .delete", function (e) {
		e.stopPropagation();
		deletePipeline($(this).parent());
	});

	$("#save").click(function () {
		savePipeline();
	});

	$("#pipelines_open").click(function () {
		$("#pipelines").addClass("opened");
		$("#overlay").fadeIn();
	});

	$("#overlay").click(function () {
		$(this).fadeOut();
		$("#pipelines").removeClass("opened");

		if(settingOpened){
			$(".close", settingElm).click();
		}

		if(schedulerOpened){
			$(".close", schedulerElm).click();
		}
	});

	$(".sidebar_close").click(function () {
		$("#overlay").fadeOut();
		$(".sidebar").removeClass("opened");
	});

	// alert on close unsaved pipeline
	$(window).bind("beforeunload",function(){
		if(isNotSaved){
			return false;
		}
	});

	schedulerInit();
});