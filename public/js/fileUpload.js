function GetFileSizeNameAndType() {
    let fi = document.getElementById('uploadFile');
    let sb = document.getElementById('submitButton');
    let fp = document.getElementById('fileProperties');

    let totalFileSize = 0;
    if (fi.files.length > 0) {
        sb.disabled = false;
        sb.style.backgroundColor = '#404044'
        for (var i = 0; i <= fi.files.length - 1; i++) {
            let fsize = fi.files.item(i).size;
            totalFileSize = totalFileSize + fsize;
            fp.style.paddingTop = "1.2vw";
            fp.innerHTML =
                'File Name: <b>' + fi.files.item(i).name + '</b>'
                +
                '<br />' + 'File Size: <b>' + Math.round((fsize / 1024)) + '</b> KB'
                +
                '<br />' + 'File Type: <b>' + fi.files.item(i).type + '</b>';

        }
    }
}