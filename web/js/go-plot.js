////////////////////////////////////////////////////////////////////////////////
//  go-plot.js  -  Dec-11-2022 by aldebap
//
//  Go-Plot Web Application
////////////////////////////////////////////////////////////////////////////////

//  show the math function plot form
function showMathFunctionPlot() {
    $('#mathPlotParameters').show();
    $('#dataSetPlotParameters').hide();
}

//  add a function to "function list"
function addFuncion() {
    let title = $('#function-title').val();
    let math_function = $('#function').val();

    //  validate function
    let validationAlert = $('#mathPlotAlert');

    validationAlert.empty();

    if (math_function == '') {
        validationAlert.append('<div class="alert alert-warning d-flex align-items-center alert-dismissible" role="alert">'
            + '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-exclamation-triangle-fill" viewBox="0 0 16 16">'
            + '<path d="M8.982 1.566a1.13 1.13 0 0 0-1.96 0L.165 13.233c-.457.778.091 1.767.98 1.767h13.713c.889 0 1.438-.99.98-1.767L8.982 1.566zM8 5c.535 0 .954.462.9.995l-.35 3.507a.552.552 0 0 1-1.1 0L7.1 5.995A.905.905 0 0 1 8 5zm.002 6a1 1 0 1 1 0 2 1 1 0 0 1 0-2z"/>'
            + '</svg>'
            + '<div>&nbsp; Need to type a function before add it</div>'
            + '<button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>'
            + '</div>');
        return;
    }

    let functionList = $('#function-list');
    let newIndex = functionList.children().length + 1;

    if (title == '') {
        functionList.append('<li class="list-group-item">'
            + '<input class="form-check-input me-1" type="checkbox" id="check-function-' + newIndex + '" onclick="funcionCheckBoxClicked();" \>'
            + '<label class="form-check-label" id="function-' + newIndex + '">' + math_function + '</label>'
            + '</li>');
    } else {
        functionList.append('<li class="list-group-item">'
            + '<input class="form-check-input me-1" type="checkbox" id="check-function-' + newIndex + '" onclick="funcionCheckBoxClicked();" \>'
            + '<label class="form-check-label" id="function-' + newIndex + '">' + title + ' : ' + math_function + '</label>'
            + '</li>');
    }

    //  enable plot button
    let plotButton = $('#btn-math-function-plot');

    plotButton.prop("disabled", false);
}

//  event function called when a "function checkbox" is clicked
function funcionCheckBoxClicked() {
    let functionList = $('#function-list');
    let deleteButton = $('#btn-delete-function');
    let enableDeleteButton = false;

    for (let i = 1; i <= functionList.children().length; i++) {
        let functionCheckbox = $('#check-function-' + i);

        if (functionCheckbox.is(':checked')) {
            enableDeleteButton = true;
            break;
        }
    }

    //  if any function is checked, enable delete button
    if (enableDeleteButton) {
        deleteButton.prop("disabled", false);
    } else {
        deleteButton.prop("disabled", true);
    }
}

//  delete all selected functions from "function list"
function deleteFuncion() {
    let functionList = $('#function-list');
    let deleteButton = $('#btn-delete-function');
    let leftFunctions = [];

    for (let i = 1; i <= functionList.children().length; i++) {
        let functionCheckbox = $('#check-function-' + i);
        let math_function = $('#function-' + i).text();

        if (!functionCheckbox.is(':checked')) {
            leftFunctions.push(math_function);
        }
    }

    //  clear the function list and add only uncheckd items
    functionList.empty();

    for (let i = 1; i <= leftFunctions.length; i++) {
        functionList.append('<li class="list-group-item">'
            + '<input class="form-check-input me-1" type="checkbox" id="check-function-' + i + '" onclick="funcionCheckBoxClicked();" \>'
            + '<label class="form-check-label" for="function-' + i + '" id="function-' + i + '">' + leftFunctions[i - 1] + '</label>'
            + '</li>');
    }

    deleteButton.prop("disabled", true);

    //  enable plot button
    let plotButton = $('#btn-math-function-plot');

    if (leftFunctions.length > 0) {
        plotButton.prop("disabled", false);
    } else {
        plotButton.prop("disabled", true);
    }
}

//  invoke plot API with all functions and parameters
function doMathFuncionPlot() {
    let x_label = $('#x-label-func').val();
    let y_label = $('#y-label-func').val();
    let min_x = $('#min-x').val();
    let max_x = $('#max-x').val();
    let functionList = $('#function-list');
    let plotCanvas = $('#canvas_plot');
    let plots = [];

    //  validate min_x and max_x
    let validationAlert = $('#mathPlotAlert');
    let formValidated = true;

    validationAlert.empty();

    if (min_x != '' && isNaN(min_x)) {
        validationAlert.append('<div class="alert alert-warning d-flex align-items-center alert-dismissible" role="alert">'
            + '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-exclamation-triangle-fill" viewBox="0 0 16 16">'
            + '<path d="M8.982 1.566a1.13 1.13 0 0 0-1.96 0L.165 13.233c-.457.778.091 1.767.98 1.767h13.713c.889 0 1.438-.99.98-1.767L8.982 1.566zM8 5c.535 0 .954.462.9.995l-.35 3.507a.552.552 0 0 1-1.1 0L7.1 5.995A.905.905 0 0 1 8 5zm.002 6a1 1 0 1 1 0 2 1 1 0 0 1 0-2z"/>'
            + '</svg>'
            + '<div>&nbsp;Min_x needs to be a valid number</div>'
            + '<button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>'
            + '</div>');
        formValidated = false;
    }
    if (max_x != '' && isNaN(max_x)) {
        validationAlert.append('<div class="alert alert-warning d-flex align-items-center alert-dismissible" role="alert">'
            + '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-exclamation-triangle-fill" viewBox="0 0 16 16">'
            + '<path d="M8.982 1.566a1.13 1.13 0 0 0-1.96 0L.165 13.233c-.457.778.091 1.767.98 1.767h13.713c.889 0 1.438-.99.98-1.767L8.982 1.566zM8 5c.535 0 .954.462.9.995l-.35 3.507a.552.552 0 0 1-1.1 0L7.1 5.995A.905.905 0 0 1 8 5zm.002 6a1 1 0 1 1 0 2 1 1 0 0 1 0-2z"/>'
            + '</svg>'
            + '<div>&nbsp;Max_x needs to be a valid number</div>'
            + '<button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>'
            + '</div>');
        formValidated = false;
    }
    if (!formValidated) {
        return;
    }

    min_x = Number($('#min-x').val());
    max_x = Number($('#max-x').val());

    //  get every function from the list
    for (let i = 1; i <= functionList.children().length; i++) {
        let functionItem = $('#function-' + i).text();
        let math_function = functionItem;
        let title = '';
        let tokens = functionItem.match(/^(.+)\s:\s(.+)$/);

        if (tokens != null) {
            math_function = tokens[2];
            title = tokens[1];
        }

        plots.push({
            title: title,
            math_function: {
                function: math_function,
                min_x: min_x,
                max_x: max_x,
            }
        });
    }

    const GOPLOT_API_URL = '/plot/api/canvas';
    const REQUEST_HEADERS = {
        'Content-Type': 'application/json',
    };
    const plotRequest = {
        x_label: x_label,
        y_label: y_label,
        plot: plots,
        width: Math.floor(plotCanvas.width()),
        height: Math.floor(plotCanvas.height()),
    };
    console.debug('request payload: ' + JSON.stringify(plotRequest))

    //  send a request to the server
    axios.post(GOPLOT_API_URL, plotRequest, { headers: REQUEST_HEADERS })
        .then(response => {
            const canvasScript = response.data;

            eval(canvasScript);
            canvas_plot();
        })
        .catch(error => {
            validationAlert.append('<div class="alert alert-danger d-flex align-items-center" role="alert">'
                + '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-sign-stop-fill" viewBox="0 0 16 16">'
                + '<path d="M10.371 8.277v-.553c0-.827-.422-1.234-.987-1.234-.572 0-.99.407-.99 1.234v.553c0 .83.418 1.237.99 1.237.565 0 .987-.408.987-1.237Zm2.586-.24c.463 0 .735-.272.735-.744s-.272-.741-.735-.741h-.774v1.485h.774Z"/>'
                + '<path d="M4.893 0a.5.5 0 0 0-.353.146L.146 4.54A.5.5 0 0 0 0 4.893v6.214a.5.5 0 0 0 .146.353l4.394 4.394a.5.5 0 0 0 .353.146h6.214a.5.5 0 0 0 .353-.146l4.394-4.394a.5.5 0 0 0 .146-.353V4.893a.5.5 0 0 0-.146-.353L11.46.146A.5.5 0 0 0 11.107 0H4.893ZM3.16 10.08c-.931 0-1.447-.493-1.494-1.132h.653c.065.346.396.583.891.583.524 0 .83-.246.83-.62 0-.303-.203-.467-.637-.572l-.656-.164c-.61-.147-.978-.51-.978-1.078 0-.706.597-1.184 1.444-1.184.853 0 1.386.475 1.436 1.087h-.645c-.064-.32-.352-.542-.797-.542-.472 0-.77.246-.77.6 0 .261.196.437.553.522l.654.161c.673.164 1.06.487 1.06 1.11 0 .736-.574 1.228-1.544 1.228Zm3.427-3.51V10h-.665V6.57H4.753V6h3.006v.568H6.587Zm4.458 1.16v.544c0 1.131-.636 1.805-1.661 1.805-1.026 0-1.664-.674-1.664-1.805V7.73c0-1.136.638-1.807 1.664-1.807 1.025 0 1.66.674 1.66 1.807ZM11.52 6h1.535c.82 0 1.316.55 1.316 1.292 0 .747-.501 1.289-1.321 1.289h-.865V10h-.665V6.001Z"/>'
                + '</svg>'
                + '<div>&nbsp;Error generating Go-Plot graphic: ' + error + '</div>'
                + '<button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>'
                + '</div>');
            console.error('On Generating Go-Plot Graphic', error)
        });
}

//  show the data set plot form
function showDataSetPlot() {
    $('#mathPlotParameters').hide();
    $('#dataSetPlotParameters').show();
}

//  upload content from a file into dataSet
function doUploadData() {
    let filenameSelector = $('#dataset-filename')[0].files;
    let validationAlert = $('#dataSetPlotAlert');

    if (filenameSelector.length == 0) {
        validationAlert.append('<div class="alert alert-warning d-flex align-items-center alert-dismissible" role="alert">'
            + '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-exclamation-triangle-fill" viewBox="0 0 16 16">'
            + '<path d="M8.982 1.566a1.13 1.13 0 0 0-1.96 0L.165 13.233c-.457.778.091 1.767.98 1.767h13.713c.889 0 1.438-.99.98-1.767L8.982 1.566zM8 5c.535 0 .954.462.9.995l-.35 3.507a.552.552 0 0 1-1.1 0L7.1 5.995A.905.905 0 0 1 8 5zm.002 6a1 1 0 1 1 0 2 1 1 0 0 1 0-2z"/>'
            + '</svg>'
            + '<div>&nbsp;A file needs to be selected</div>'
            + '<button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>'
            + '</div>');
    }

    let file = filenameSelector[0];
    let datasetReader = new FileReader();

    datasetReader.onload = (event) => {
        const datasetContent = event.target.result;
        let dataSet = $('#dataSet');

        dataSet.val(datasetContent);
    };

    datasetReader.onerror = (event) => {
        validationAlert.append('<div class="alert alert-warning d-flex align-items-center alert-dismissible" role="alert">'
            + '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-exclamation-triangle-fill" viewBox="0 0 16 16">'
            + '<path d="M8.982 1.566a1.13 1.13 0 0 0-1.96 0L.165 13.233c-.457.778.091 1.767.98 1.767h13.713c.889 0 1.438-.99.98-1.767L8.982 1.566zM8 5c.535 0 .954.462.9.995l-.35 3.507a.552.552 0 0 1-1.1 0L7.1 5.995A.905.905 0 0 1 8 5zm.002 6a1 1 0 1 1 0 2 1 1 0 0 1 0-2z"/>'
            + '</svg>'
            + '<div>&nbsp;Error reading dataset file: ' + event.target.error.name + '</div>'
            + '<button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>'
            + '</div>');
    };

    datasetReader.readAsText(file);
}

//  add a pair of columns from dataSet to "dataSet list"
function addDataSet() {
    let title = $('#dataset-title').val();
    let column_x = $('#column-x').val();
    let column_y = $('#column-y').val();
    let style = $('#point-style').val();

    //  validate column_x and column_y
    let validationAlert = $('#dataSetPlotAlert');
    let formValidated = true;

    validationAlert.empty();

    if (isNaN(column_x) || Number(column_x) <= 0) {
        validationAlert.append('<div class="alert alert-warning d-flex align-items-center alert-dismissible" role="alert">'
            + '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-exclamation-triangle-fill" viewBox="0 0 16 16">'
            + '<path d="M8.982 1.566a1.13 1.13 0 0 0-1.96 0L.165 13.233c-.457.778.091 1.767.98 1.767h13.713c.889 0 1.438-.99.98-1.767L8.982 1.566zM8 5c.535 0 .954.462.9.995l-.35 3.507a.552.552 0 0 1-1.1 0L7.1 5.995A.905.905 0 0 1 8 5zm.002 6a1 1 0 1 1 0 2 1 1 0 0 1 0-2z"/>'
            + '</svg>'
            + '<div>&nbsp;Column X needs to be a valid number and greater than zero</div>'
            + '<button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>'
            + '</div>');
        formValidated = false;
    }
    if (isNaN(column_y) || Number(column_y) <= 0) {
        validationAlert.append('<div class="alert alert-warning d-flex align-items-center alert-dismissible" role="alert">'
            + '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-exclamation-triangle-fill" viewBox="0 0 16 16">'
            + '<path d="M8.982 1.566a1.13 1.13 0 0 0-1.96 0L.165 13.233c-.457.778.091 1.767.98 1.767h13.713c.889 0 1.438-.99.98-1.767L8.982 1.566zM8 5c.535 0 .954.462.9.995l-.35 3.507a.552.552 0 0 1-1.1 0L7.1 5.995A.905.905 0 0 1 8 5zm.002 6a1 1 0 1 1 0 2 1 1 0 0 1 0-2z"/>'
            + '</svg>'
            + '<div>&nbsp;Column Y needs to be a valid number and greater than zero</div>'
            + '<button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>'
            + '</div>');
        formValidated = false;
    }
    if (!formValidated) {
        return;
    }

    let dataSetList = $('#dataSet-list');
    let newIndex = dataSetList.children().length + 1;

    if (title == '') {
        dataSetList.append('<li class="list-group-item">'
            + '<input class="form-check-input me-1" type="checkbox" id="check-dataSet-' + newIndex + '" onclick="dataSetCheckBoxClicked();" \>'
            + '<label class="form-check-label" id="dataSet-' + newIndex + '">[' + column_x + ':' + column_y + '] ' + style + '</label>'
            + '</li>');
    } else {
        dataSetList.append('<li class="list-group-item">'
            + '<input class="form-check-input me-1" type="checkbox" id="check-dataSet-' + newIndex + '" onclick="dataSetCheckBoxClicked();" \>'
            + '<label class="form-check-label" id="dataSet-' + newIndex + '">' + title + ' : [' + column_x + ':' + column_y + '] ' + style + '</label>'
            + '</li>');
    }

    //  enable plot button
    let plotButton = $('#btn-dataSet-plot');

    plotButton.prop("disabled", false);
}

//  event function called when a "dataSet checkbox" is clicked
function dataSetCheckBoxClicked() {
    let dataSetList = $('#dataSet-list');
    let deleteButton = $('#btn-delete-dataSet');
    let enableDeleteButton = false;

    for (let i = 1; i <= dataSetList.children().length; i++) {
        let dataSetCheckbox = $('#check-dataSet-' + i);

        if (dataSetCheckbox.is(':checked')) {
            enableDeleteButton = true;
            break;
        }
    }

    //  if any dataSet is checked, enable delete button
    if (enableDeleteButton) {
        deleteButton.prop("disabled", false);
    } else {
        deleteButton.prop("disabled", true);
    }
}

//  delete all selected dataSets from "dataSet list"
function deleteDataSet() {
    let dataSetList = $('#dataSet-list');
    let deleteButton = $('#btn-delete-dataSet');
    let leftDataSets = [];

    for (let i = 1; i <= dataSetList.children().length; i++) {
        let dataSetCheckbox = $('#check-dataSet-' + i);
        let dataSet = $('#dataSet-' + i).text();

        if (!dataSetCheckbox.is(':checked')) {
            leftDataSets.push(dataSet);
        }
    }

    //  clear the dataSet list and add only uncheckd items
    dataSetList.empty();

    for (let i = 1; i <= leftDataSets.length; i++) {
        dataSetList.append('<li class="list-group-item">'
            + '<input class="form-check-input me-1" type="checkbox" id="check-dataSet-' + i + '" onclick="dataSetCheckBoxClicked();" \>'
            + '<label class="form-check-label" id="dataSet-' + i + '">' + leftDataSets[i - 1] + '</label>'
            + '</li>');
    }

    deleteButton.prop("disabled", true);

    //  enable plot button
    let plotButton = $('#btn-dataSet-plot');

    if (leftDataSets.length > 0) {
        plotButton.prop("disabled", false);
    } else {
        plotButton.prop("disabled", true);
    }
}

function doDataSetPlot() {

    let x_label = $('#x-label-data').val();
    let y_label = $('#y-label-data').val();
    let dataSetList = $('#dataSet-list');
    let dataSet = $('#dataSet').val();
    let lines = dataSet.split("\n");
    let plotCanvas = $('#canvas_plot');
    let plots = [];

    //  get every dataSet item from the list and parse column_x, column_y and style
    for (let i = 1; i <= dataSetList.children().length; i++) {
        let dataSetItem = $('#dataSet-' + i).text();
        let title = '';
        let column_x = '';
        let column_y = '';
        let style = '';
        let tokens = dataSetItem.match(/^(.+)\s:\s\[(\d+):(\d+)\]\s(.+)$/);

        if (tokens != null) {
            title = tokens[1];
            column_x = tokens[2];
            column_y = tokens[3];
            style = tokens[4];
        } else {
            tokens = dataSetItem.match(/\[(\d+):(\d+)\]\s(.+)$/);

            if (tokens != null) {
                column_x = tokens[1];
                column_y = tokens[2];
                style = tokens[3];
            }
        }

        //  parse dataset to get points data and add it to the list of plots
        let points = [];
        let lineNumber = 0;

        lines.forEach(line => {
            let linePoints = line.split(" ");

            //  ignore empty lines
            if (line.length == 0) {
                return;
            }

            //  ignore first line if it's a comment
            if (0 == lineNumber && line[0] == '#') {
                return;
            }

            points.push({
                x: Number(linePoints[column_x - 1]),
                y: Number(linePoints[column_y - 1]),
            });

            lineNumber++;
        });

        plots.push({
            title: title,
            data_set: {
                points: points,
                style: style,
            }
        });
    }

    const GOPLOT_API_URL = '/plot/api/canvas';
    const REQUEST_HEADERS = {
        'Content-Type': 'application/json',
    };
    const plotRequest = {
        x_label: x_label,
        y_label: y_label,
        plot: plots,
        width: Math.floor(plotCanvas.width()),
        height: Math.floor(plotCanvas.height()),
    };
    console.debug('request payload: ' + JSON.stringify(plotRequest))

    //  send a request to the server
    axios.post(GOPLOT_API_URL, plotRequest, { headers: REQUEST_HEADERS })
        .then(response => {
            const canvasScript = response.data;

            eval(canvasScript);
            canvas_plot();
        })
        .catch(error => {
            validationAlert.append('<div class="alert alert-danger d-flex align-items-center" role="alert">'
                + '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-sign-stop-fill" viewBox="0 0 16 16">'
                + '<path d="M10.371 8.277v-.553c0-.827-.422-1.234-.987-1.234-.572 0-.99.407-.99 1.234v.553c0 .83.418 1.237.99 1.237.565 0 .987-.408.987-1.237Zm2.586-.24c.463 0 .735-.272.735-.744s-.272-.741-.735-.741h-.774v1.485h.774Z"/>'
                + '<path d="M4.893 0a.5.5 0 0 0-.353.146L.146 4.54A.5.5 0 0 0 0 4.893v6.214a.5.5 0 0 0 .146.353l4.394 4.394a.5.5 0 0 0 .353.146h6.214a.5.5 0 0 0 .353-.146l4.394-4.394a.5.5 0 0 0 .146-.353V4.893a.5.5 0 0 0-.146-.353L11.46.146A.5.5 0 0 0 11.107 0H4.893ZM3.16 10.08c-.931 0-1.447-.493-1.494-1.132h.653c.065.346.396.583.891.583.524 0 .83-.246.83-.62 0-.303-.203-.467-.637-.572l-.656-.164c-.61-.147-.978-.51-.978-1.078 0-.706.597-1.184 1.444-1.184.853 0 1.386.475 1.436 1.087h-.645c-.064-.32-.352-.542-.797-.542-.472 0-.77.246-.77.6 0 .261.196.437.553.522l.654.161c.673.164 1.06.487 1.06 1.11 0 .736-.574 1.228-1.544 1.228Zm3.427-3.51V10h-.665V6.57H4.753V6h3.006v.568H6.587Zm4.458 1.16v.544c0 1.131-.636 1.805-1.661 1.805-1.026 0-1.664-.674-1.664-1.805V7.73c0-1.136.638-1.807 1.664-1.807 1.025 0 1.66.674 1.66 1.807ZM11.52 6h1.535c.82 0 1.316.55 1.316 1.292 0 .747-.501 1.289-1.321 1.289h-.865V10h-.665V6.001Z"/>'
                + '</svg>'
                + '<div>&nbsp;Error generating Go-Plot graphic: ' + error + '</div>'
                + '<button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>'
                + '</div>');
            console.error('On Generating Go-Plot Graphic', error)
        });
}
