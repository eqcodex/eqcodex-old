{{define "content"}}
<link href="/css/icons.css" rel="stylesheet">
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
                            <tr><td style="padding: 20px 0px;"><span class="item-slot"></span><span class="item-icon icon-{{.Item.Icon}}"></span> {{.Item.Name}}</td></tr>
                                <tr><td colspan=7>{{.Item.Header_line}}</td></tr>
                                <tr><td colspan=7>Class: {{.Item.Class_line}}</td></tr>
                                <tr><td colspan=7>Race: {{.Item.Race_line}}</td></tr>
                                <tr><td colspan=7>{{.Item.Slot_line}}</td></tr>
                                <tr><td colspan=7></td></tr>
                                <tr><td>Size:</td><td>{{.Item.Size_line}}</td><td></td><td>AC:</td><td>{{.Item.Ac}}</td><td></td><td></td></tr>
                                <tr><td>Weight:</td><td>{{.Item.Weight_line}}</td><td></td><td>HP:</td><td>{{.Item.Hp}}</tr><td></td><td></td></tr>
                                <tr><td>Rec Level:</td><td>{{.Item.Reclevel}}</td><td></td><td>Mana:</td><td>{{.Item.Mana}}</tr><td></td><td></td></tr>
                                <tr><td>Req Level:</td><td>{{.Item.Reqlevel}}</td><td></td><td>Endur:</td><td>{{.Item.Endur}}</tr><td></td><td></td></tr>
                                <tr><td colspan=3></td>                                    <td>Purity:</td><td>{{.Item.Purity}}</tr><td></td><td></td></tr>
                                <tr><td colspan=3></td>                                    <td>Haste:</td><td>{{.Item.Haste}}%</tr><td></td><td></td></tr>
                                <tr><td>Strength:</td><td>{{.Item.Astr}}</td><td>+{{.Item.Heroic_str}}</td><td>Magic:</td><td>{{.Item.Mr}}</tr><td>Attack:</td><td>{{.Item.Attack}}</td></tr>
                                <tr><td>Stamina:</td><td>{{.Item.Asta}}</td><td>+{{.Item.Heroic_sta}}</td><td>Fire:</td><td>{{.Item.Fr}}</tr><td>HP Regen:</td><td>{{.Item.Regen}}</td></tr>
                                <tr><td>Intelligence:</td><td>{{.Item.Aint}}</td><td>+{{.Item.Heroic_int}}</td><td>Cold:</td><td>{{.Item.Cr}}</tr><td>Mana Regen:</td><td>{{.Item.Manaregen}}</td></tr>
                                <tr><td>Wisdom:</td><td>{{.Item.Awis}}</td><td>+{{.Item.Heroic_wis}}</td><td>Disease:</td><td>{{.Item.Dr}}</tr><td>Heal Amount:</td><td>{{.Item.Healamt}}</td></tr>
                                <tr><td>Agility:</td><td>{{.Item.Aagi}}</td><td>+{{.Item.Heroic_agi}}</td><td>Poison:</td><td>{{.Item.Pr}}</tr><td>Spell Dmg:</td><td>{{.Item.Spelldmg}}</td></tr>
                                <tr><td>Dexterity:</td><td>{{.Item.Adex}}</td><td>+{{.Item.Heroic_dex}}</td><td>Corruption:</td><td>{{.Item.Svcorruption}}</tr><td>Clairvoyance:</td><td>{{.Item.Clairvoyance}}</td></tr>
                                <tr><td>Charisma:</td><td>{{.Item.Acha}}</td><td>+{{.Item.Heroic_cha}}</td><td colspan=4></td></tr>
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