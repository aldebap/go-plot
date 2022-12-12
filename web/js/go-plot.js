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

    functionList.append('<li class="list-group-item">'
        + '<input class="form-check-input me-1" type="checkbox" id="check-function-' + newIndex + '" onclick="funcionCheckBoxClicked();" \>'
        + '<label class="form-check-label" for="function-' + newIndex + '" id="function-' + newIndex + '">' + math_function + '</label>'
        + '</li>');

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

    let title = $('#math-plot-title').val();
    let min_x = $('#min-x').val();
    let max_x = $('#max-x').val();
    let functionList = $('#function-list');
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
        let math_function = $('#function-' + i).text();

        plots.push({
            title: title,
            function: math_function,
            min_x: min_x,
            max_x: max_x,
        });
    }

    const GOPLOT_API_URL = '/plot/api/canvas';
    const REQUEST_HEADERS = {
        'Content-Type': 'application/json',
    };
    const plotRequest = {
        plot: plots,
    };
    console.debug('request payload: ' + JSON.stringify(plotRequest))

    //  send a request to the server
    axios.post(GOPLOT_API_URL, plotRequest, { headers: REQUEST_HEADERS })
        .then(response => {
            const canvasScript = response.data;

            eval(canvasScript);
            canvas_plot();
        })
        .catch(error => console.error('On Generating Go-Plot Graphic', error));
}

//  show the data set plot form
function showDataSetPlot() {
    $('#mathPlotParameters').hide();
    $('#dataSetPlotParameters').show();
}

function doDataSetPlot() {

    let title = $('#dataset-title').val();
    let style = $('#point-style').val();
    let column_x = Number($('#column-x').val());
    let column_y = Number($('#column-y').val());
    let dataSet = $('#dataSet').val();
    let lines = dataSet.split("\n");
    let points = [];

    //  parse dataset to get points data
    lines.forEach(line => {
        let linePoints = line.split(" ");

        points.push({
            x: Number(linePoints[column_x - 1]),
            y: Number(linePoints[column_y - 1]),
        });
    });

    const GOPLOT_API_URL = '/plot/api/canvas';
    const REQUEST_HEADERS = {
        'Content-Type': 'application/json',
    };
    const plotRequest = {
        plot: [
            {
                title: title,
                style: style,
                points: points,
            },
        ],
    };
    console.debug('request payload: ' + JSON.stringify(plotRequest))

    //  send a request to the server
    axios.post(GOPLOT_API_URL, plotRequest, { headers: REQUEST_HEADERS })
        .then(response => {
            const canvasScript = response.data;

            eval(canvasScript);
            canvas_plot();
        })
        .catch(error => console.error('On Generating Go-Plot Graphic', error));
}
