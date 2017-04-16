var buildListElm;

/** Build List
 *
 */
function buildList()
{
	// append the build list
	$("ul", buildListElm).html("");
	$.each(pipeline.Builds, function (i, b) {
		$("ul", buildListElm).append("<li data-id=\"" + b + "\"> Build #" + b + "<span class=\"remove material-icons\">close</span></li>");
	});
}

/** Add Log
 *
 * @param log
 */
function addLog(log)
{
	var html = "<div class=\"log\" style=\"margin-left: " + ((log.Level * 3) + 1) + "%\">" +
		"<h3>" + log.TaskName + "<span>" + log.Command + "</span></h3>" +
		"<h5>Console Output:</h5>" +
		"<pre>" + log.Console + "</pre>" +
		"</div>";

	$("#build_logs").append(html);
}

/** Run Build
 *
 */
function runBuild()
{
	// check if the pipeline saved or not
	if(isNotSaved && !confirm("The board not saved, Continue?")){
		return;
	}

	var url = (base + "pipelines/" + pipeline.ID + "/build").split("://")[1];

	var ws = new WebSocket("ws://" + url);
	ws.onmessage = function(e) {
		// check if the build finished
		if(e.data === "done") {
			pipeline.LastBuild++;
			pipeline.Builds = pipeline.Builds ? pipeline.Builds : [];
			pipeline.Builds.push(pipeline.LastBuild);
			buildList();
		} else if(e.data === "fail"){
			alert("Build Fail");
		} else{
			addLog(JSON.parse(e.data));
		}
	};

	$("#build_details").fadeIn();
}

/** Load Build
 *
 * @param pipeline
 * @param id
 */
function loadBuild(pipeline, id)
{
	$("#build_details").fadeIn();
	$("ul", buildListElm).slideUp();

	api("GET", "pipelines/" + pipeline + "/build/" + id, function(data) {
		$.each(data.Logs, function (i, log) {
			addLog(log);
		});
	});
}

/** Delete Build
 *
 * @param li
 */
function deleteBuild(li)
{
	var buildID = li.data("id");

	// send the delete request to the api
	api("DELETE", "pipelines/" + pipeline.ID + "/build/" + buildID, function( data ) {
		if(!data.success) {
			alert("Build delete fail");
			return;
		}

		// remove the build id and the element
		pipeline.Builds.splice(pipeline.Builds.indexOf(buildID), 1);
		li.remove();
	});
}

$(document).ready(function () {
	buildListElm = $("#build_list");

	// run new build
	$("#build").click(function () {
		runBuild();
		$("ul", buildListElm).slideUp();
	});

	// open build list
	$("#build_list_toggle").click(function () {
		$("ul", buildListElm).slideToggle();
	});

	// open build details
	$("ul", buildListElm).on("click", "li", function () {
		loadBuild(pipeline.ID, $(this).data("id"));
	});

	// close build details panel
	$("#build_close").click(function () {
		$("#build_details").hide();
		$("#build_logs").html("");
	});

	// remove a build
	$("ul", buildListElm).on("click", ".remove", function (e) {
		e.stopPropagation();
		deleteBuild($(this).parent());
	});
});