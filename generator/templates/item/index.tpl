{{define "content"}}
        <div id="page-wrapper">
            <div class="row">
                <div class="col-lg-12">
                    <h1 class="page-header">{{.Item.Name}}</h1>
                </div>
                <!-- /.col-lg-12 -->
            </div>
            <!-- /.row -->
            <div class="row">
                <div class="col-lg-12">
                    <div class="panel panel-default">
                        <div class="panel-body"> 
                            <table>
                            <th></th>
                                <tr>
                                    <td>CLASS: {{.Item.Classes}}</td>
                                </tr>
                                <tr>
                                    <td>RACE: {{.Item.Races}}</td>
                                </tr>
                                <tr>
                                    <td>AC: {{.Item.Ac}}</td>
                                </tr>
                                <tr>
                                    <td>HP: {{.Item.Hp}}</td>
                                </tr>
                                <tr>
                                    <td>Mana: {{.Item.Mana}}</td>
                                </tr>
                                <tr>
                                    <td>Lore: {{.Item.Lore}}</td>
                                </tr>
                            </table>                            
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
                            <p>Mobs that drop {{.Item.Name}}</p>
                            <input type="button" value="Filter" id="filtermob">

                            <div class="table-responsive">
                                <table class="table table-bordered table-striped table-hover" id="item-table">
                                    <thead>
                                        <tr>
                                            <th>NPC</th>
                                            <th>Level</th>
                                            <th>Quest</th>
                                            <th>Zone</th>
                                        </tr>                                
                                        </thead>
                            <tbody>
                            {{ range $key, $value := .NPCs}}
                            <tr>
                                <td>{{/*<a href="{{ $value.Url }}">*/}}{{ $value.Name }}{{/*</a>*/}}</td>
                                <td>{{ $value.Level }}</td>
                                <td>{{ $value.Quest }}</td>
                                <td><a href="{{ $value.Zone_url }}">{{ $value.Zone_name }}</a></td>
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
            <script>
            $(document).ready(function() {
                $('#filtermob').click(function() {
                    $('#item-table').DataTable({
                        responsive: true
                    });
                    $('#advanced').hide();
                });
            });
            </script>
{{end}}