$(document).ready(function () {

    function getTableElement(element, tag) {
        return element.parent().parent().find(tag);
    }

    function bind() {
        $(".file-rename").click(function () {
            var name = getTableElement($(this), ".file-name").attr("data-name");
            $("#win-rename-modal #module-name").val(name);
        });

        $(".file-remove").click(function () {
            var name = getTableElement($(this), ".file-name").attr("data-name");
            console.info(name);
        });

        $(".load-path").click(function () {
            var name = getTableElement($(this), ".file-name").attr("data-name");
            load(name);
        });
    }

    function make(url, name, date, dir) {
        html = "";
        html += "                    <tr>\n";
        html += "                      <td class=\"file-name\" data-name=\"" + url + ((url.substring(url.length - 1) != "/" && name != "..") ? "/" : "") + (name != ".." ? name : "") + "\">";
        if (dir) {
            html += "<a href=\"javascript:void(0);\" class=\"load-path\">"
        } else {
            html += "<a href=\"http://download.heroin.so" + url + (url.substring(url.length - 1) != "/" ? "/" : "" ) + name + "\" target=\"_blank\">";
        }
        html += name + (dir ? "/" : "");
        html += "</a></td>\n";
        html += "                      <td class=\"center\">" + date + "</td>\n";
        html += "                      <td class=\"center\">\n";
        html += "                        <a class=\"file-rename\" title=\"rename\" href=\"#win-rename-modal\" data-toggle=\"modal\"><i class=\"icon-edit\"></i></a>\n";
        html += "                        <a class=\"file-remove\" title=\"remove\" href=\"javascript:void(0);\"><i class=\"icon-remove\"></i></a>\n";
        html += "                      </td>\n";
        html += "                    </tr>\n";
        return html;
    }

    function load(url) {
        url = url == "" ? "/" : url;
        $.get("list?path=" + url, function (result) {
            $("#tab-content tbody tr").remove();
            html = "";
            if (url != "/") {
                parent_url = url.substring(0, url.lastIndexOf("/") + 1);
                if (parent_url.substring(parent_url.length - 1) == "/") {
                    parent_url = parent_url.substring(0, parent_url.length - 1);
                }
                html += make(parent_url, "..", "", true);
            }
            $(result).find("root").find("file").each(function (i) {
                dir = eval($(this).children("dir").text());
                date = $(this).children("date").text();
                name = $(this).children("name").text();
                html += make(url, name, date, dir)
            });
            $("#tab-content tbody").append(html);
            bind();
        });
    }

    load("");

    $("#win-rename-modal input").focus(function () {
        $(this).parent().parent().removeClass("error").removeClass("success");
    }).blur(function () {
            $(this).parent().parent().addClass($.trim($(this).val()) === '' ? "error" : "success");
        });

    $("#btn-download").click(function () {
        $.get("download",
            {"url": $("#down-file").val(), "name": $("#save-name").val(), "path": $("#save-path").val()},
            function (data) {
                var result = eval("(" + data + ")")
                console.info(result)
                if (result.code > 0) {
                    alert("下载进行中....");
                    load("");
                }
            }
        );
    });
});