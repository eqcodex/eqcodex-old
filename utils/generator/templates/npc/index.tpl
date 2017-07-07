{{define "content"}}
<link href="/css/icons.css" rel="stylesheet">
        <div id="page-wrapper">
            <div class="row">
                <div class="col-lg-12">
                    <h1 class="page-header">{{.Npc.Name}}</h1>
                </div>
                <!-- /.col-lg-12 -->
            </div>
            <!-- /.row -->
            <div class="row">
                <div class="col-lg-6">
                    <div class="panel panel-default">
                        <div class="panel-body"> 

                            <table>
                            <tr><td>Level:</td><td>{{.Npc.Level}}</td></tr>
                            <tr><td>HP:</td><td>{{.Npc.Hp}}</td></tr>
                            <tr><td>Damage:</td><td>{{.Npc.Mindmg}}-{{.Npc.Maxdmg}}</td></tr>
                            <tr><td>Zone:</td><td>{{.Npc.Zone_long_name}}</td></tr>
                            </table>
                        </div>
                    </div>
                </div>
                <div class="col-lg-6">
                    <div class="panel panel-default">
                        <div class="panel-body"> 
                    <div id="map_container" style="width: 100%; height: 200px; border:1px solid black; ">
                        <svg id="map_svg" xmlns="http://www.w3.org/2000/svg" style="display: inline; width: inherit; min-width: inherit; max-width: inherit; height: inherit; min-height: inherit; max-height: inherit;"></svg>
                    </div>
                        </div>
                    </div>
                </div>
                <!-- /.col-lg-12 -->
            </div>
            <!-- /.row -->


             <!-- /.row -->
            <div class="row">
                <div class="col-lg-12">                
                    <div class="panel panel-default">
                        <div class="panel-body">                            
                            <p>{{.Npc.Name}} Drops</p>
                            <input type="button" value="Filter" id="advanced">

                            <div class="table-responsive">
                                <table class="table table-bordered table-striped table-hover" id="item-table">
                                    <thead>
                                        <tr>
                                            <th>Item</th>
                                            <th>Category</th>
                                            <th>Era</th>
                                            <th>Quest</th>
                                        </tr>                                
                                        </thead>
                            <tbody>
                            {{ range $key, $value := .Items }}
                            <tr>
                                <td><a href="{{ $value.Url }}">{{ $value.Name }}</a></td>
                                <td>{{ $value.Category }}</td>
                                <td>{{ $value.Era }}</td>
                                <td>{{ $value.Quest }}</td>
                            </div>
                            {{ end }}
                            </tbody>
                            </table>
                                   
                            <p></p>
                        </div>
                    </div>
                </div>
                <!-- /.col-lg-12 -->
            </div>
            <!-- /.row -->

            <script type="text/javascript" src="/js/map/svg-pan-zoom.js"></script>
    <script type="text/javascript" src="/js/map/snap.svg.js"></script>
            <script>
$(document).ready(function() {
    $('#filtermob').click(function() {
        $('#item-table').DataTable({
            responsive: true
        });
        $('#advanced').hide();
    });


    var spawn_points;
    var s = Snap("#map_svg");
    var map_group;

    var spawnData = unescape("{{.Npc.MapData}}");
    
    var json = jQuery.parseJSON(spawnData);
    map_group = s.g();

    if (json.spawn_points === null) {
        s.text(0,0,"No spawn data");
        panZoomInstance = svgPanZoom('#map_svg', {
             zoomEnabled: true,
             controlIconsEnabled: true,
             fit: true,
             center: true,
             //zoom: 0.1,
             //minZoom: 0.1
           //  beforePan: beforePan,
         });
        panZoomInstance.zoom(0.1)
        return
    }

    if (json.lines === null) {
        s.text(0,0,"No map data");
        panZoomInstance = svgPanZoom('#map_svg', {
             zoomEnabled: true,
             controlIconsEnabled: true,
             fit: true,
             center: true,
             //zoom: 0.1,
             //minZoom: 0.1
           //  beforePan: beforePan,
         });

        panZoomInstance.zoom(0.1)
        return
    }

   /* var beforePan
    beforePan = function(oldPan, newPan){
      var stopHorizontal = false
        , stopVertical = false
        , gutterWidth = 100
        , gutterHeight = 100
          // Computed variables
        , sizes = this.getSizes()
        , leftLimit = -((sizes.viewBox.x + sizes.viewBox.width) * sizes.realZoom) + gutterWidth
        , rightLimit = sizes.width - gutterWidth - (sizes.viewBox.x * sizes.realZoom)
        , topLimit = -((sizes.viewBox.y + sizes.viewBox.height) * sizes.realZoom) + gutterHeight
        , bottomLimit = sizes.height - gutterHeight - (sizes.viewBox.y * sizes.realZoom)
      customPan = {}
      customPan.x = Math.max(leftLimit, Math.min(rightLimit, newPan.x))
      customPan.y = Math.max(topLimit, Math.min(bottomLimit, newPan.y))
      return customPan
    }*/


    if (json.lines !== null && json.lines.length > 0) {
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

    if (json.spawn_points !== null &&  json.spawn_points.length > 0) {
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
        //https://github.com/ariutta/svg-pan-zoom
        
    }

    panZoomInstance = svgPanZoom('#map_svg', {
         zoomEnabled: true,
         controlIconsEnabled: true,
         fit: true,
         center: true,
         //zoom: 0.1,
         //minZoom: 0.1
       //  beforePan: beforePan,
     });
    //if (json.spawn_points.length == 1) {
     //   panZoomInstance.zoom(0.2);
    //}
});

            </script>
{{end}}