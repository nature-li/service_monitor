// 保存数据并更新页面
function save_data_and_update_page_view(data) {
    if (data.success == "true") {
        // db 最大idx
        var db_page_count = Math.ceil(data.item_count / window.save_data.view_item_count_per_page);
        window.save_data.db_total_item_count = data.item_count;
        window.save_data.db_max_page_idx = db_page_count - 1;

        // 清空部分数据
        window.save_data.item_list = [];

        // 保存所有数据
        if (data.content != null) {
            var length = data.content.length;
            window.save_data.db_return_item_count = length;
            for (var i = 0; i < length; i++) {
                var item = data.content[i];
                window.save_data.item_list.push(item);
            }
        }

        // 更新view
        update_page_view(window.save_data.view_current_page_idx);
    } else {
        $.showErr("发生错误");
    }
}

// 更新分页标签
function update_page_partition(view_current_page_idx) {
    var db_max_page_idx = window.save_data.db_max_page_idx;
    var db_total_item_count = window.save_data.db_total_item_count;
    var db_return_item_count = window.save_data.db_return_item_count;
    var view_max_page_count = window.save_data.view_max_page_count;
    var view_start_page_idx = window.save_data.view_start_page_idx;
    var view_item_count_per_page = window.save_data.view_item_count_per_page;

    var view_current_page_count = view_max_page_count;
    if (db_max_page_idx - view_start_page_idx + 1 < view_max_page_count) {
        view_current_page_count = db_max_page_idx - view_start_page_idx + 1;
    }
    window.save_data.view_current_page_count = view_current_page_count;
    window.save_data.view_current_page_idx = view_current_page_idx;

    // 清空元素
    var $partition = $("#page_partition_id");
    $partition.find("li").remove();

    // 首页
    if (view_start_page_idx > 0) {
        $partition.append('<li class="hive-page-first"><a>首页</a></li>');
    }

    // 上一批
    if (view_start_page_idx > 0) {
        $partition.append('<li class="hive-page-pre-batch"><a>&laquo;</a></li>');
    }

    // 上一页
    if (view_current_page_idx > 0)
    {
        $partition.append('<li class="hive-page-pre"><a>上一页</a></li>');
    }

    // 页码
    for (var i = 0; i < view_current_page_count; i++) {
        var page_number = view_start_page_idx + i + 1;

        var page_class = "hive-page-partition";
        if (view_start_page_idx + i == view_current_page_idx) {
            page_class = "disabled hive-page-partition";
        }
        $partition.append('<li class="' + page_class + '"><a>' + page_number + '</a></li>');
    }

    // 下一页
    if (view_current_page_idx < db_max_page_idx) {
        $partition.append('<li class="hive-page-next"><a>下一页</a></li>');
    }

    // 下一批
    if (view_start_page_idx + view_current_page_count < db_max_page_idx) {
        $partition.append('<li class="hive-page-next-batch"><a>&raquo;</a></li>');
    }

    // 尾页
    if (view_start_page_idx + view_current_page_count <= db_max_page_idx) {
        $partition.append('<li class="hive-page-tail"><a>尾页</a></li>');
    }

    // 更新tips
    var data_count_tip = "";
    if (db_return_item_count > 0) {
        var start_item_number = view_current_page_idx * view_item_count_per_page + 1;
        var end_item_number = start_item_number + db_return_item_count - 1;
        data_count_tip = "第" + start_item_number + "-" + end_item_number + "条, 共" + db_total_item_count + "条数据";
    } else {
        data_count_tip = "共0条数据";
    }
    $("#data_count_tip").html(data_count_tip);
}

// 点击首页
$(document).on("click", ".hive-page-first", function () {
    window.save_data.view_current_page_idx = 0;
    window.save_data.view_start_page_idx = 0;
    query_and_update_view();
});

// 点击上一页
$(document).on("click", ".hive-page-pre", function () {
    window.save_data.view_current_page_idx -= 1;
    if (window.save_data.view_current_page_idx < window.save_data.view_start_page_idx) {
        window.save_data.view_start_page_idx -= window.save_data.view_max_page_count;
    }
    query_and_update_view();
});

// 点击上一批
$(document).on("click", ".hive-page-pre-batch", function () {
    window.save_data.view_start_page_idx -= window.save_data.view_max_page_count;
    window.save_data.view_current_page_idx = window.save_data.view_start_page_idx;
    query_and_update_view();
});

// 点击页码
$(document).on("click", ".hive-page-partition", function () {
    window.save_data.view_current_page_idx = $(this).find('a').text() - 1;
    query_and_update_view();
});

// 点击下一批
$(document).on("click", ".hive-page-next-batch", function () {
    window.save_data.view_start_page_idx += window.save_data.view_max_page_count;
    window.save_data.view_current_page_idx = window.save_data.view_start_page_idx;
    query_and_update_view();
});

// 点击下一页
$(document).on("click", ".hive-page-next", function () {
    window.save_data.view_current_page_idx += 1;
    if (window.save_data.view_current_page_idx >= window.save_data.view_start_page_idx + window.save_data.view_current_page_count) {
        window.save_data.view_start_page_idx += window.save_data.view_current_page_count;
    }
    query_and_update_view();
});


// 点击尾页
$(document).on("click", ".hive-page-tail", function () {
    window.save_data.view_current_page_idx = window.save_data.db_max_page_idx;

    var db_page_count = window.save_data.db_max_page_idx + 1;
    var batch = Math.ceil(db_page_count / window.save_data.view_max_page_count);
    var pre_batch = batch - 1;
    window.save_data.view_start_page_idx = pre_batch * window.save_data.view_max_page_count;
    query_and_update_view();
});