var edit_file_page_interval_func = null;
var old_file_version = "";
var old_refer_link = "";
var old_file_desc = "";

// init_edit_page
function handle_delete_file_response(data) {
    if (data.success != true) {
        $.showErr("删除失败");
        return
    }

    $("#file_desc_info_rows").addClass("no-display");
    $("#back_to_home_tips_row").removeClass("no-display");
    $("#back_to_home_btn_row").removeClass("no-display");
    edit_file_page_interval_func = setInterval(changeReloadTips, 1000);
}

// 定时任务
function changeReloadTips() {
    var seconds = $("#jump_after_seconds").text();
    if (seconds == "1") {
        if (edit_file_page_interval_func != null) {
            clearInterval(edit_file_page_interval_func);
            edit_file_page_interval_func = null;

            window.location.replace("/");
            return
        }
    }

    seconds = seconds - 1;
    $("#jump_after_seconds").html(seconds);
}

// 页面加载时
$(document).ready(function () {
    $("#delete_file_btn").click(function (e) {
        $.showConfirm("确定要删除吗?", delete_file_after_confirm);
    });
});

// 确认后删除文件
function delete_file_after_confirm() {
    var file_id = $("#edit_file_id").text();

    // 加载数据
    $.ajax({
            url: '/delete_file_api',
            type: "post",
            data: {
                'file_id': file_id
            },
            dataType: 'json',
            success: function (data) {
                handle_delete_file_response(data);
            },
            error: function (jqXHR, textStatus, errorThrown) {
                if (jqXHR.status == 302) {
                    window.parent.location.replace("/");
                } else {
                    $.showErr("删除失败");
                }
            }
        }
    );
}

// 返回首页
$(document).on("click", "#back_to_home", function (e) {
    window.location.replace("/");
});

// 编辑按钮
$(document).on("click", "#edit_file_btn", function (e) {
    old_file_version = $("#edit_file_version").text();
    $("#edit_file_version").html("<input id='new_file_version' type='text' class='form-control' maxlength='32' value='" + old_file_version + "'/>");

    old_refer_link = $("#a_edit_refer_link").attr("href");
    $("#td_edit_refer_link").html("<input id='new_refer_link' type='text' class='form-control' maxlength='32' value='" + old_refer_link + "'/>");

    old_file_desc = $("#edit_file_desc").text();
    $("#edit_file_desc").html("<textarea id='new_file_desc' class='form-control' rows='10' maxlength='1024'>" + old_file_desc + "</textarea>")

    $("#edit_file_btn").addClass("no-display");
    $("#submit_file_btn").removeClass("no-display");
});

// 提交按钮
$(document).on("click", "#submit_file_btn", function (e) {
    var new_file_version = $("#new_file_version").val();
    var new_refer_link = $("#new_refer_link").val();
    var new_file_desc = $("#new_file_desc").val();

    if (old_file_version === new_file_version
        && old_refer_link === new_refer_link
        && old_file_desc === new_file_desc) {
        $("#edit_file_version").html(old_file_version);
        $("#td_edit_refer_link").html('<a id="a_edit_refer_link" target="_blank" href="' + old_refer_link + '">' + old_refer_link + '</a>');
        $("#edit_file_desc").html(old_file_desc);
        $("#edit_file_btn").removeClass("no-display");
        $("#submit_file_btn").addClass("no-display");
        return
    }

    $.ajax({
            url: '/edit_file_api',
            type: "post",
            data: {
                'file_id': $("#edit_file_id").text(),
                'file_version': new_file_version,
                'file_refer_link': new_refer_link,
                'file_desc': new_file_desc
            },
            dataType: 'json',
            success: function (data) {
                edit_file_call_back(data);
            },
            error: function (jqXHR, textStatus, errorThrown) {
                if (jqXHR.status == 302) {
                    window.parent.location.replace("/");
                } else {
                    edit_file_failed();
                }
            }
        }
    );
});

// 编辑文件回调函数
function edit_file_call_back(data) {
    if (data.success !== true) {
        edit_file_failed();
        return;
    }

    window.location.reload();
}

// 编辑失败
function edit_file_failed() {
    $("#edit_file_version").html(old_file_version);
    $("#td_edit_refer_link").html('<a id="a_edit_refer_link" target="_blank" href="' + old_refer_link + '">' + old_refer_link + '</a>');
    $("#edit_file_desc").html(old_file_desc);
    $("#edit_file_btn").removeClass("no-display");
    $("#submit_file_btn").addClass("no-display");
    $.showErr("提交失败");
}