let log = console.log;

$("#search_file").on('submit', function (e) {
    e.preventDefault();
    $.ajax({
        type: "POST",
        url: "/search_file",
        data: JSON.stringify({
            query: $(this).find('input[name="query"]').val()
        }),
        success: function (fileRecords) { // res: Array({filename: string, owner_ip: string})
            let foundFilesDiv = $("#found_files");
            foundFilesDiv.html("");
            foundFilesDiv.append("<h3>Search results: </h3>");
            if (fileRecords.length === 0) {
                foundFilesDiv.html("No files found");
            } else {
                fileRecords.map(fileRecord => {
                    let fileRecPar = $('<a>', {
                        href: "/download_file?filename=" + fileRecord.filename + "&owner_ip=" + fileRecord.owner_ip,
                        text: fileRecord.filename + " @ " + fileRecord.owner_ip,
                        style: "font-size: 18px"
                    });
                    foundFilesDiv.append(fileRecPar);
                    let info = $('p', {
                        text: "Available at " + fileRecord.owner_ip
                    });
                    foundFilesDiv.append(info);
                    foundFilesDiv.append("<br>");
                    foundFilesDiv.append("<br>");
                })
            }
        },
        error: function () {
            log("Sending message failed. Please, try again.")
        }
    });
});

$("#share_file_btn").click(function (e) {
    e.preventDefault();
    let fd = new FormData($("#share_file_form")[0]);
    $.ajax({
        url: "/upload_file",
        type: "POST",
        data: fd,
        success: () => {
            let uploadResultDiv = $("#upload_result");
            uploadResultDiv.html("Upload successful! 😃");
            setTimeout(() => uploadResultDiv.html(""), 5000)
        },
        error: function (err) {
            let uploadResultDiv = $("#upload_result");
            uploadResultDiv.html("Upload unsuccessful! 💔");
            setTimeout(() => uploadResultDiv.html(""), 5000)
        },
        cache: false,
        contentType: false,
        processData: false,
    });
});