var submitting = 0;

// 页面加载时
$(document).ready(function () {
    $("#upload_file_btn").click(function (e) {
        submitting += 1;
        if (submitting > 1) {
            e.preventDefault();
            submitting -= 1;
            return;
        }

        var file_info = $("#choose_file_for_upload").val();
        if (file_info == "") {
            $("#upload_file_error_label").removeClass("hidden-self");
            $("#upload_file_error_label").find("span").html("请选择文件");
            e.preventDefault();
            submitting -= 1;
            return;
        }

        var file_size = $('#choose_file_for_upload')[0].files[0].size;
        var file_limit = $("#max_file_limit").val();
        if (file_size > file_limit) {
            e.preventDefault();
            submitting -= 1;
            return;
        }

        var bar = $('#upload_file_progress_bar');
        var percent = $('#upload_file_process_percent');

        $('form').ajaxForm({
            dataType: 'json',
            beforeSend: function() {
                $("#upload_file_progress").removeClass("no-display");
                var percentVal = '0%';
                bar.width(percentVal);
                percent.html(percentVal);
            },
            uploadProgress: function(event, position, total, percentComplete) {
                var percentVal = percentComplete + '%';
                bar.width(percentVal);
                percent.html(percentVal);
            },
            success: function(data) {
                if (data.code == "200") {
                    submitting = 0;
                    $("#upload_file_form").addClass("no-display");
                    $("#upload_file_form_again").removeClass("no-display");
                } else {
                    submitting = 0;
                    $("#upload_file_progress").addClass("no-display");
                    $("#upload_file_error_label").removeClass("hidden-self");
                    $("#upload_file_error_label").find("span").html("上传失败");
                }
            },
            error: function (jqXHR, textStatus, errorThrown) {
                submitting = 0;
                $("#upload_file_progress").addClass("no-display");
                $("#upload_file_error_label").removeClass("hidden-self");
                $("#upload_file_error_label").find("span").html("上传失败");
            }
        });
    });

    // 再传一个
    $("#upload_again_btn").click(function (e) {
       window.location.reload();
    });

    // 检测文件大小
    $("#choose_file_for_upload").change(function () {
        if (this.files == null) {
            return;
        }

        if (this.files.length < 1) {
            return;
        }
        var file = this.files[0];

        // 清楚错误提示信息
        $("#upload_file_error_label").addClass("hidden-self");
        $("#upload_file_error_label").find("span").html("");

        // 显示文件名
        $('#upload-file-info').html(file.name);

        // 计算并显示长度
        var text = "该文件大小为: ";
        if (file.size > 1024 * 1024 * 1024) {
            text += (file.size / 1024.0 / 1024.0 / 1024.0).toFixed(2) + "G";
        }
        else if (file.size > 1024 * 1024) {
            text += (file.size / 1024.0 / 1024.0).toFixed(2) + "M";
        }
        else if (file.size > 1024) {
            text += (file.size / 1024.0).toFixed(2) + "K";
        } else {
            text += length + "B";
        }
        $("#check_file_size_label").html(text);

        var file_limit = $("#max_file_limit").val();
        if (file.size > file_limit) {
            $("#upload_file_error_label").removeClass("hidden-self");
            $("#upload_file_error_label").find("span").html("文件过大");
        } else {
            $("#upload_file_error_label").addClass("hidden-self");
            $("#upload_file_error_label").find("span").html("");
        }
    });
});

// 返回首页
$(document).on("click", "#back_to_home", function (e) {
    window.location.replace("/");
});