
        {{define "nav"}}
        <a class="navbar-brand" href="#">
            <img alt="" src="/static/img/shortcut-icon.png" width="25px" height="25px">
        </a>
        <a class="navbar-brand" href="/">Allen Blog</a>
        <ul class="nav navbar-nav">
            <li {{if .IsHome}} class="active" {{end}}><a href="/">首页</a></li>
            <li {{if .IsCategory}} class="active" {{end}}><a href="/category">分类</a></li>
            <li {{if .IsTopic}} class="active" {{end}}><a href="/topic">文章</a></li>
        </ul>
    <div class="pull-right" style="display: inline-block">
        <ul class="nav navbar-nav">
            {{if .IsLogin}}
            <li><a href="/login?exit=true">退出</a> </li>
            {{else}}
            <li><a href="/login">登录</a> </li>
            {{end}}
        </ul>
    </div>
        {{end}}
    </div>
</nav>
