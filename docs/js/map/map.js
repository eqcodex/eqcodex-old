var spawn_points;
var s = Snap("#map_svg");
var map_group;

$.ajax({
		type: "GET",
		url: "/map/kael.json",
		success: function (data) {
			json = data

			//console.log("Data", data)
			//var json = jQuery.parseJSON(data);
			console.log(json);

			map_group = s.g();

			if (json.lines.length > 0) {
				for (var i = 0; i < json.lines.length; i++ ) {
					var curLine = json.lines[i];
					// Fix the line data for our svg
					curLine.x1 = (curLine.x1 + 2000) / 5;
					curLine.y1 = (curLine.y1 + 2000) / 5;
					curLine.x2 = (curLine.x2 + 2000) / 5;
					curLine.y2 = (curLine.y2 + 2000) / 5;

					var line = s.line(curLine.x1, curLine.y1, curLine.x2, curLine.y2).attr({strokeWidth:1, stroke:"green"});
					map_group.add(line);
				}
			}

			if (json.spawn_points.length > 0) {
				for (var i = 0; i < json.spawn_points.length; i++) {
					var spawn = json.spawn_points[i];
					// Fix the point data for our svg
					spawn.x = (spawn.x + 2000) / 5;
					spawn.y = (spawn.y + 2000) / 5;

					var circle = s.circle(spawn.x, spawn.y, 2).attr({fill: 'maroon', stroke: 'red', strokeWidth: 1});
					map_group.add(circle);
					circle.spawn = spawn;
					circle.click(function(event) {
							console.log(this);
					});
				}

				panZoomInstance = svgPanZoom('#map_svg', {
					 zoomEnabled: true,
					 controlIconsEnabled: true,
					 fit: true,
					 center: true,
					 minZoom: 0.1
				 });

			 panZoomInstance.zoom(0.8)
			}
		},
		error: function () {
				console.log("ajax related error occured");
		}
});
