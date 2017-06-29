{{define "content"}}

    
    
        <div id="page-wrapper">
            <div class="row">
                <div class="col-lg-12">
                    <h1 class="page-header">Zone Leveling Chart</h1>
                </div>
                <!-- /.col-lg-12 -->
            </div>
            <!-- /.row -->
            <div class="row">
                <div class="col-lg-12">
                    <div class="panel panel-default">
                        <div class="panel-body"> 
                            <h3>Zone Leveling Chart</h3>
                            <p>Below is a list of zones based on level.</p>

                            <div class="table-responsive">
                                <table class="table table-bordered table-striped table-hover" id="zone-table">
                                    <thead>
                                        <tr>
                                            <th>Name</th>
                                            <th>1</th>
                                            <th>5</th>
                                            <th>10</th>
                                            <th>15</th>
                                            <th>20</th>
                                            <th>25</th>
                                            <th>30</th>
                                            <th>35</th>
                                            <th>40</th>
                                            <th>45</th>
                                            <th>50</th>
                                            <th>55</th>
                                            <th>60</th>
                                            <td><span class="fa fa-question-circle-o"></span></td>
                                        </tr>
                                    </thead>
                                    <tbody>
                                    {{ range $key, $value := .Zones }}
                                        <tr>
                                            <th><a href="{{ $value.Url }}">{{ $value.Long_name.String }}</a></th>
                                            <td>{{ if $value.IsLevel $value.Levels 1}}<span class="fa fa-check"></span>{{ end }}</td>
                                            <td>{{ if $value.IsLevel $value.Levels 5}}<span class="fa fa-check"></span>{{ end }}</td>
                                            <td>{{ if $value.IsLevel $value.Levels 10}}<span class="fa fa-check"></span>{{ end }}</td>
                                            <td>{{ if $value.IsLevel $value.Levels 15}}<span class="fa fa-check"></span>{{ end }}</td>
                                            <td>{{ if $value.IsLevel $value.Levels 20}}<span class="fa fa-check"></span>{{ end }}</td>
                                            <td>{{ if $value.IsLevel $value.Levels 25}}<span class="fa fa-check"></span>{{ end }}</td>
                                            <td>{{ if $value.IsLevel $value.Levels 30}}<span class="fa fa-check"></span>{{ end }}</td>
                                            <td>{{ if $value.IsLevel $value.Levels 35}}<span class="fa fa-check"></span>{{ end }}</td>
                                            <td>{{ if $value.IsLevel $value.Levels 40}}<span class="fa fa-check"></span>{{ end }}</td>
                                            <td>{{ if $value.IsLevel $value.Levels 45}}<span class="fa fa-check"></span>{{ end }}</td>
                                            <td>{{ if $value.IsLevel $value.Levels 50}}<span class="fa fa-check"></span>{{ end }}</td>
                                            <td>{{ if $value.IsLevel $value.Levels 55}}<span class="fa fa-check"></span>{{ end }}</td>
                                            <td>{{ if $value.IsLevel $value.Levels 60}}<span class="fa fa-check"></span>{{ end }}</td>
                                            <td><a target="_blank" href="https://github.com/eqcodex/eqcodex/issues/new?title=zone.html+-+{{ $value.Short_name.String }}&body=I+found+an+issue+with+eqcodex+as+follows:"><span class="fa fa-question-circle-o"></span></a></td>
                                        </tr>
                                    {{ end }}                                        
                                    </tbody>
                                </table>
                            </div>
                            <p></p>
                        </div>
                    </div>
                </div>
                <!-- /.col-lg-12 -->
            </div>
            <!-- /.row -->
             <script>
    $(document).ready(function() {
        $('#zone-table').DataTable({
            responsive: true
        });
    });
    </script>
{{end}}