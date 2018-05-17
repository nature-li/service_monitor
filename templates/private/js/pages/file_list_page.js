// 初始化全局变量
function reset_save_data() {
    window.save_data = {
        'item_list': [],
        'db_total_item_count': 0,
        'db_return_item_count': 0,
        'db_max_page_idx': 0,
        'view_max_page_count': 5,
        'view_item_count_per_page': 15,
        'view_start_page_idx': 0,
        'view_current_page_idx': 0,
        'view_current_page_count': 0
    };
}

// 查询并更新页面
function query_and_update_view() {
    var off_set = window.save_data.view_current_page_idx * window.save_data.view_item_count_per_page;
    var limit = window.save_data.view_item_count_per_page;

    // 加载数据
    $.ajax({
            url: '/file_list_api',
            type: "post",
            data: {
                'file_name': $("#search_file_name").val(),
                'off_set': off_set,
                'limit': limit
            },
            dataType: 'json',
            success: function (data) {
                save_data_and_update_page_view(data);
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

// 更新表格和分页
function update_page_view(page_idx) {
    // 更新表格
    var html = "";
    for (var i = 0; i < window.save_data.item_list.length; i++) {
        var item = window.save_data.item_list[i];
        var file_size;
        if (item.file_size > 1073741824) {
            length = item.file_size / 1073741824.0;
            file_size = length.toFixed(2) + "G"
        } else if (item.file_size > 1048576) {
            length = item.file_size / 1048576.0;
            file_size = length.toFixed(2) + "M"
        } else if (item.file_size > 1024) {
            length = item.file_size / 1024.0;
            file_size = length.toFixed(2) + "K"
        } else {
            file_size = item.file_size;
        }

        var length = item.file_size / 1024.0 / 1024.0;
        length = length.toFixed(2);
        html += "<tr>" +
            "<td style='display:none'>" + item.id + "</td>" +
            "<td><a href='" + item.file_url + "'download='" + item.file_name + "'>" + item.file_name + " <span class='glyphicon glyphicon-download-alt' aria-hidden='true'></span></a></td>" +
            "<td>" + file_size + "</td>" +
            "<td>" + item.version + "</td>" +
            "<td>" + item.create_time + "</td>" +
            "<td>" + item.md5_value + "</td>" +
            "<td><a href='#' class='a_file_detail'>详细信息</a></td>" +
            "</tr>";
    }
    $("#file_list_result").find("tr:gt(0)").remove();
    $("#file_list_result").append(html);

    // 更新分页标签
    update_page_partition(page_idx);
}

// 条件搜索
function init_page() {
    // 定义全局变量
    if (!window.save_data) {
        reset_save_data();
    }

    // 查询数据并更新页面
    query_and_update_view();
}

// 页面加载时
$(document).ready(function () {
    init_page();

    $("#search_file_name_btn").click(function () {
        init_page();
    });

    $("#search_file_name").keydown(function (event) {
        if (event.keyCode == "13") {
            event.preventDefault();
            init_page();
        }
    });
});

// 详细信息绑定函数
$(document).on("click", ".a_file_detail", function () {
    var file_id = $(this).parent().parent().find("td:eq(0)").text();
    var edit_url = "/edit_file?id=" + file_id;
    window.location.replace(edit_url);
});