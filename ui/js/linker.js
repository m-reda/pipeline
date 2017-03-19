(function ( $ ) { $.fn.linker = function(options)
{
    var lk = this,
        settings = $.extend({}, options );

    lk.addClass('linker_container').append('<svg id="linker_svg"></svg>');

    /*
     *  nodes
     */
    this.node = function (data)
    {
        var node = data ? data : {x: 0, y: 0};

        node.__id = gid();
        node.id = node.id ? node.id : node.__id;
        node.el = $('<div class="linker_node node_'+node.id+'" style="left:'+node.x+'px;top: '+node.y+'px;"><h3>'+node.name+'</h3><div class="linker_inputs"></div><div class="linker_outputs"></div></div>');
        node.el.data('obj', node);
        node.paths_ot = {}; // paths out from this node
        node.paths_in = {}; // paths in to this node

        // add input to the node
        node.inputs = [];
        node.input = function (id, name)
        {
            var i = node.inputs.push({
                __id: gid(),
                id: id,
                node: node,
                el: $('<div class="linker_point" data-type="input"></div>')
            });

            var input = node.inputs[i - 1];
            input.el.data('obj', input);
            node.paths_in[input.__id] = [];


            var label = $('<div class="linker_label"><span>'+name+'</span></span></div>').append(input.el);
            $('.linker_inputs', node.el).append(label);
            return input;
        };

        // add output to the node
        node.outputs = [];
        node.output = function (id, name)
        {
            var i = node.outputs.push({
                __id: gid(),
                id: id,
                node: node,
                el: $('<div class="linker_point" data-type="output"></div>')
            });


            var output = node.outputs[i - 1];
            output.el.data('obj', output);
            node.paths_ot[output.__id] = [];

            output.connect = function (input, withoutEvent)
            {
                var path = drawPath(this.el.offset(), input.el.offset()),
                    conn = [path, this, input];

                node.paths_ot[this.__id].push(conn);
                input.node.paths_in[input.__id].push(conn);

                if(!withoutEvent &&  this.onConnect)
                    this.onConnect(input);

                // remove connection
                $(path).click(function () {
                    if(!confirm('Delete this path?'))
                        return;

                    var out_idx = node.paths_ot[output.__id].indexOf(conn);
                    node.paths_ot[output.__id].splice(out_idx, 1);
                    input.node.paths_in[input.__id].splice(1, 1);
                    $(this).remove();

                    if(output.onRemove)
                        output.onRemove(out_idx);
                });
            };

            var label = $('<div class="linker_label"><span>'+name+'</span></div>').append(output.el);
            $('.linker_outputs', node.el).append(label);
            return output;
        };

        // on drug
        node.el.on('drag', function (e)
        {
             if(!node.onDrag)
                 return;

            node.onDrag(
                parseInt($(this).css('left')),
                parseInt($(this).css('top'))
            )
        });

        // add the node to the linker container
        lk.append(node.el);

        return node;
    };


    /*
     *  linking
     */
    var selected_output = null, drag_path, drag_path_pos;
    lk.on('click', function (e) {

        var el = $(e.target),
            isPoint = el.hasClass('linker_point'),
            isOutput = (el.data('type') == 'output');

        // if there is a selected output
        // check if the new on is input point
        if(selected_output)
        {
            // connect the output and the input
            if(isPoint && !isOutput) {
                var output = selected_output.data('obj'),
                    input = el.data('obj');

                // do not connect if the input and the output in the same node
                if(output.node != input.node)
                    output.connect(input)
            }

            // clear the selected output and remove the draggable path
            selected_output.removeClass('selected');
            selected_output = null;
            $(drag_path).remove();
            drag_path = null;
            lk.removeClass('drag_path');

            return
        }

        // if no output selected yet
        // select this and add draggable path
        if(isPoint && isOutput) {
            selected_output = $(e.target).addClass('selected');
            drag_path_pos = selected_output.offset();
            drag_path = drawPath(drag_path_pos, drag_path_pos);
            lk.addClass('drag_path');
        }
    });


    /*
     *  dragging
     */
    var drag_node, drag_width = 0;

    lk.on('mousedown', '.linker_node > h3', function (e) {
        drag_node   = $(e.target).parent();
        drag_width  = drag_node.width() / 2;

    })
    .on('mouseup', function (e) {
        drag_node = null;
    })
    .on("mousemove", function(e) {
        // drag
        if (drag_node) {
            drag_node.offset({top: e.pageY - 10, left: e.pageX - drag_width}).trigger('drag');

            // update paths
            $.each(drag_node.data('obj').paths_ot, function (_, arr) {
                $.each(arr, function (_, p) {
                    var p1 = p[1].el.offset(),
                        p2 = p[2].el.offset();
                    $(p[0]).attr('d', curve(p1.left, p1.top, p2.left, p2.top))
                });
            });
            $.each(drag_node.data('obj').paths_in, function (_, arr) {
                $.each(arr, function (_, p) {
                    var p1 = p[1].el.offset(),
                        p2 = p[2].el.offset();
                    $(p[0]).attr('d', curve(p1.left, p1.top, p2.left, p2.top))
                });
            });
        }

        // path
        if(drag_path) {
            $(drag_path).attr('d', curve(drag_path_pos.left, drag_path_pos.top, e.pageX, e.pageY))
        }
    });

    /*
     *  general
     */
    // draw new path
    function drawPath(p1, p2) {
        var p = document.createElementNS('http://www.w3.org/2000/svg', 'path');
        p.setAttribute('d', curve(p1.left, p1.top, p2.left, p2.top));
        document.getElementById('linker_svg').appendChild(p);

        return p;
    }

    // calculate the path curve
    function curve(x1, y1, x2, y2)
    {
        var d = Math.abs(x1-x2) / 2;
        y1 += 5;
        y2 += 5;

		return " M" +  x1      + "," + y1 +
			   " C" + (x1 + d) + "," + y1 +
			   " "  + (x2 - d) + "," + y2 +
			   " "  + x2       + "," + y2;
    }

    // generate id
    var id_counter = Date.now();
    function gid() {
        id_counter++;
        return Math.random().toString(36) + id_counter;
    }

    return this;
}
}( jQuery ));