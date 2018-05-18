$(document).ready(function () {
    // 定义全局变量
    if (!window.save_data) {
        reset_save_data();
    }

    // 查询全部用户并更新列表
    query_and_update_view();
});

// 初始化全局变量
function reset_save_data() {
    window.save_data = {
        'item_list': [],
        'db_total_item_count': 0,
        'db_return_item_count': 0,
        'db_max_page_idx': 0,
        'view_max_page_count': 5,
        'view_item_count_per_page': 10,
        'view_start_page_idx': 0,
        'view_current_page_idx': 0,
        'view_current_page_count': 0,
    };
}

// 查询数据并更新页面
function query_and_update_view() {
    var off_set = window.save_data.view_current_page_idx * window.save_data.view_item_count_per_page;
    var limit = window.save_data.view_item_count_per_page;

    $.ajax({
            url: '/list_user_api',
            type: "post",
            data: {
                'user_email': $("#search_user_email").val(),
                'off_set': off_set,
                'limit': limit
            },
            dataType: 'json',
            success: function (response) {
                save_data_and_update_page_view(response);
            },
            error: function (jqXHR, textStatus, errorThrown) {
                if (jqXHR.status == 302) {
                    window.parent.location.replace("/");
                } else {
                    $.showErr("查询失败");
                }
            }
        }
    );
}

// 更新页表
function update_page_view(page_idx) {
    // 删除表格
    $('#user_list_result > tbody  > tr').each(function () {
        $(this).remove();
    });

    // 添加表格
    for (var i = 0; i < window.save_data.item_list.length; i++) {
        var user = window.save_data.item_list[i];
        add_row(user.id, user.user_email, user.manager_right, user.create_time);
    }

    // 更新分页标签
    update_page_partition(page_idx);
}

// 全选多选框
$("#check_all").click(function () {
    if ($(this).prop('checked')) {
        $("#user_list_result").find("input[name='user_list[]']").each(function (i, e) {
            $(e).prop('checked', true);
        })
    } else {
        $("#user_list_result").find("input[name='user_list[]']").each(function (i, e) {
            $(e).prop('checked', false);
        })
    }
});

// 点击删除用户操作
$("#del_user_button").click(function () {
    var count = 0;
    $('#user_list_result > tbody  > tr').each(function () {
        var $check_box = $(this).find("td:eq(0)").find("input[name='user_list[]']");
        if ($check_box.prop('checked')) {
            count += 1;
        }
    });

    if (count > 0) {
        $.showConfirm("确定要删除吗?", query_delete_selected_user);
    }
});

// 删除用户操作
function query_delete_selected_user() {
    var content = '';
    $('#user_list_result > tbody  > tr').each(function () {
        var $check_box = $(this).find("td:eq(0)").find("input[name='user_list[]']");
        if ($check_box.prop('checked')) {
            var user_id = $check_box.val();
            content += user_id + ","
        }
    });

    // 发送请求删除后台数据
    if (content !== '') {
        $.ajax({
                url: '/del_user_api',
                type: "post",
                data: {
                    'user_id_list': content
                },
                dataType: 'json',
                success: function (response) {
                    refresh_view(response);
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
}

// 编辑用户信息
$(document).on("click", ".user-edit-button", function () {
    var $tr = $(this).parent().parent();
    var user_id = $tr.find("td:eq(1)").text();
    var user_email = $tr.find("td:eq(2)").text();
    var manager_right_txt = $tr.find("td:eq(3)").text();

    var manager_right = false;
    if (manager_right_txt === '是') {
        manager_right = true;
    }

    show_edit_dialog(user_id, user_email, manager_right);
});

// 弹出编辑对话框
function show_edit_dialog(user_id, user_email, manager_right) {
    BootstrapDialog.show({
        message: function (dialog) {
            // header
            var content = '<div>';

            // id
            content += '<div style="display: none"><input id="edit_user_id" type="text" class="form-control" value="' + user_id + '" disabled></div>';

            // 账号
            content += '<div><input id="edit_user_email" type="text" class="form-control" value="' + user_email + '" disabled></div>';

            // 用户权限
            content += '<div class="checkbox">';
            content += '<span style="margin-right: 30px;">用户权限:</span>';

            if (manager_right) {
                content += '<label style="margin: 0 10px;"><input id="manager_right_in_dialog" type="checkbox" name="manager_right" value="1" checked/>用户管理</label>';
            } else {
                content += '<label style="margin: 0 10px;"><input id="manager_right_in_dialog" type="checkbox" name="manager_right" value="0"/>用户管理</label>';
            }

            content += '</div>';

            // footer
            content += '</div>';

            return content;
        },
        title: "编辑用户",
        closable: false,
        draggable: true,
        buttons: [{
            label: '确定',
            action: function (dialogItself) {
                // 获取用户添加数据
                var user_id = $("#edit_user_id").val();

                var new_manager_right = false;
                if ($("#manager_right_in_dialog").prop('checked')) {
                    new_manager_right = true;
                }

                // 权限发生变化后发送请求
                if ( manager_right !== new_manager_right) {
                    // 发送请求
                    $.ajax({
                            url: '/edit_user_api',
                            type: "post",
                            data: {
                                'user_id': user_id,
                                'manager_right': new_manager_right
                            },
                            dataType: 'json',
                            success: function (response) {
                                edit_user_page_view(response);
                            },
                            error: function (jqXHR, textStatus, errorThrown) {
                                if (jqXHR.status === 302) {
                                    window.parent.location.replace("/");
                                } else {
                                    $.showErr("更新失败");
                                }
                            }
                        }
                    );
                }

                dialogItself.close();
            }
        },
            {
                label: '取消',
                action: function (dialogItself) {
                    dialogItself.close();
                }
            }]
    });
}

// 根据response更新用户列表
function edit_user_page_view(response) {
    if (!response.success) {
        $.showErr("更新失败");
        return;
    }

    var user = response.content;
    var user_id = user.id;
    var user_email = user.user_email;
    var manager_right = user.manager_right;
    var create_time = user.create_time;

    $('#user_list_result > tbody  > tr').each(function () {
        var $check_box = $(this).find("td:eq(0)").find("input[name='user_list[]']");
        var bind_user_id = $check_box.val();

        if (bind_user_id === user_id.toString()) {
            var user_control = '是';
            if (!manager_right) {
                user_control = '否';
            }

            $(this).find("td:eq(1)").html(user_id);
            $(this).find("td:eq(2)").html(user_email);
            $(this).find("td:eq(3)").html(user_control);
            $(this).find("td:eq(4)").html(create_time);
        }
    });
}

// 增加用户
$("#add_user_button").click(function () {
    BootstrapDialog.show({
        message: function (dialog) {
            // header
            var content = '<div>';

            // 账号
            content += '<div><input id="add_user_email" type="text" class="form-control" placeholder="输入账号"></div>';

            // 权限
            content += '<div class="checkbox">';
            content += '<span style="margin-right: 30px;">权限:</span>';
            content += '<label style="margin: 0 10px;"><input id="manager_control_in_dialog" type="checkbox" name="user_right[]" value="0" />用户管理</label>';
            content += '</div>';

            // footer
            content += '</div>';

            return content;
        },
        title: "增加用户",
        closable: false,
        draggable: true,
        buttons: [{
            label: '确定',
            action: function (dialogItself) {
                // 获取用户添加数据
                var user_email = $("#add_user_email").val();

                var manager_right = false;
                if ($("#manager_control_in_dialog").prop('checked')) {
                    manager_right = true;
                }

                // 发送请求
                $.ajax({
                        url: '/add_user_api',
                        type: "post",
                        data: {
                            'user_email': user_email,
                            'manager_right': manager_right
                        },
                        dataType: 'json',
                        success: function (response) {
                            refresh_view(response);
                        },
                        error: function (jqXHR, textStatus, errorThrown) {
                            if (jqXHR.status === 302) {
                                window.parent.location.replace("/");
                            } else {
                                $.showErr("添加失败");
                            }
                        }
                    }
                );

                // 关闭窗口
                dialogItself.close();
            }
        },
            {
                label: '取消',
                action: function (dialogItself) {
                    dialogItself.close();
                }
            }]
    });
});

// 根据ajax返回值更新页面
function refresh_view(data) {
    if (data.success === true) {
        window.location.reload();
    } else {
        $.showErr("添加失败");
    }
}

// 在表格中增加用户
function add_row(user_id, user_account, manager_right, create_time) {
    var user_control = '是';
    if (!manager_right) {
        user_control = '否';
    }

    var table = $("#user_list_result");
    var tr = $('<tr>' +
        '<td style="text-align:center;"><input name="user_list[]" type="checkbox" value="' + user_id + '"></td>' +
        '<td style="text-align:center;">' + user_id + '</td>' +
        '<td style="text-align:center;">' + user_account + '</td>' +
        '<td style="text-align:center;">' + user_control + '</td>' +
        '<td style="text-align:center;">' + create_time + '</td>' +
        '<td style="text-align:center;"><button type="button" class="btn btn-primary user-edit-button">编辑</button></td>');
    table.append(tr);
}

// 点击查找用户按钮
$(document).on("click", "#search_user_name_btn", function(e){
    // 清空数据并设置查找账号
    reset_save_data();

    // 查询数据并更新页面
    query_and_update_view();
});

