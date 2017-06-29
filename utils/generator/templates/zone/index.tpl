{{define "content"}}
        <div id="page-wrapper">
            <div class="row">
                <div class="col-lg-12">
                    <h1 class="page-header">{{.Zone.Long_name.String}}</h1>
                </div>
                <!-- /.col-lg-12 -->
            </div>
            <!-- /.row -->
            <div class="row">
                <div class="col-lg-12">                
                    <div class="panel panel-default">
                        <div class="panel-body">                            
                            <p>{{.Zone.Description}}</p>
                            <input type="button" value="Filter" id="advanced">

                            <div class="table-responsive">
                                <table class="table table-bordered table-striped table-hover" id="item-table">
                                    <thead>
                                        <tr>
                                            <th>Item</th>
                                            <th>Category</th>
                                            <th>Era</th>
                                            <th>Quest</th>
                                            <th>NPC</th>
                                        </tr>                                
                                        </thead>
                            <tbody>
                            {{ range $key, $value := .Items }}
                            <tr>
                                <td><a href="{{ $value.Url }}">{{ $value.Name }}</a></td>
                                <td>{{ $value.Category }}</td>
                                <td>{{ $value.Era }}</td>
                                <td>{{ $value.Quest }}</td>
                                <td>{{ $value.NPC }}</td>
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
        $('#advanced').click(function() {
            $('#item-table').DataTable({
                responsive: true
            });
            $('#advanced').hide();
        });
    });
    </script>
{{end}}