(function () {
    "use strict";

    var cellToHighlight;
    var messageBanner;

    // 每次加载新页面时都必须运行初始化函数。
    Office.initialize = function (reason) {
        $(document).ready(function () {
            // 初始化通知机制并隐藏它
            // var element = document.querySelector('.MessageBanner');
            // messageBanner = new components.MessageBanner(element);
            // messageBanner.hideBanner();
            
            // 如果未使用 Excel 2016，请使用回退逻辑。
            // if (!Office.context.requirements.isSetSupported('ExcelApi', '1.1')) {
            //     $("#template-description").text("此示例将显示电子表格中选定单元格的值。");
            //     $('#button-text').text("显示!");
            //     $('#button-desc').text("显示所选内容");

            //     $('#highlight-button').click(displaySelectedCells);
            //     return;
            // }

            // $("#template-description").text("此示例将突出显示电子表格中选定单元格的最高值。");
            // $('#button-text').text("突出显示!");
            // $('#button-desc').text("突出显示最大数字。");
                
            // loadSampleData();
            // getCompanyData()
            // 为突出显示按钮添加单击事件处理程序。
            // $('#highlight-button').click(hightlightHighestValue);
            $('#get-data').click(getCompanyData);
        });
    };

    function loadSampleData() {
        var values = [
            [Math.floor(Math.random() * 1000), Math.floor(Math.random() * 1000), Math.floor(Math.random() * 1000)],
            [Math.floor(Math.random() * 1000), Math.floor(Math.random() * 1000), Math.floor(Math.random() * 1000)],
            [Math.floor(Math.random() * 1000), Math.floor(Math.random() * 1000), Math.floor(Math.random() * 1000)]
        ];

        // 针对 Excel 对象模型运行批处理操作
        Excel.run(function (ctx) {
            // 为活动工作表创建代理对象
            var sheet = ctx.workbook.worksheets.getActiveWorksheet();
            // 将向电子表格写入示例数据的命令插入队列
            // var resp = webRequest()
            values[0][0] = 1
            sheet.getRange("B3:D5").values = values;

            // 运行排队的命令，并返回承诺表示任务完成
            return ctx.sync();
        })
        .catch(errorHandler);
    }

    function hightlightHighestValue() {
        // 针对 Excel 对象模型运行批处理操作
        Excel.run(function (ctx) {
            // 创建选定范围的代理对象，并加载其属性
            var sourceRange = ctx.workbook.getSelectedRange().load("values, rowCount, columnCount");

            // 运行排队的命令，并返回承诺表示任务完成
            return ctx.sync()
                .then(function () {
                    var highestRow = 0;
                    var highestCol = 0;
                    var highestValue = sourceRange.values[0][0];

                    // 找到要突出显示的单元格
                    for (var i = 0; i < sourceRange.rowCount; i++) {
                        for (var j = 0; j < sourceRange.columnCount; j++) {
                            if (!isNaN(sourceRange.values[i][j]) && sourceRange.values[i][j] > highestValue) {
                                highestRow = i;
                                highestCol = j;
                                highestValue = sourceRange.values[i][j];
                            }
                        }
                    }

                    cellToHighlight = sourceRange.getCell(highestRow, highestCol);
                    sourceRange.worksheet.getUsedRange().format.fill.clear();
                    sourceRange.worksheet.getUsedRange().format.font.bold = false;

                    // 突出显示该单元格
                    cellToHighlight.format.fill.color = "orange";
                    cellToHighlight.format.font.bold = true;
                })
                .then(ctx.sync);
        })
        .catch(errorHandler);
    }

    function getCompanyData() {
        Excel.run(function (ctx) {
            var sourceRange = ctx.workbook.getSelectedRange().load("values, rowCount, columnCount, columnIndex, rowIndex");
            // var sourceRange = ctx.workbook.getSelectedRange().load("values, rowCount, columnCount");
            console.log(sourceRange)
            return ctx.sync().then(function () {
                var selectValues = sourceRange.values;
                var columnIndex = sourceRange.columnIndex
                var rowIndex = sourceRange.rowIndex
                getData(selectValues).then(function(res) {
                    console.log("---------------")
                    console.log(res)
                    var values = []
                    for (var i = 0; i < res.length; i++) {
                        values.push([res[i].saic_legal_name, res[i].founded_on, res[i].registered_address])
                    }
                    console.log(values)
                    displayData(values,columnIndex,rowIndex)
                })
            }).then(ctx.sync);
        }).catch(errorHandler);
    }

    function getData(selectValues) {
        var ids = []
        for (var i = 0; i < selectValues.length; i++) {
            var id = selectValues[i][0]+""
            if (id.length == 0){
                continue
            }
            ids.push(id)
        }
        let url = "/match/" + JSON.stringify(ids); // This is a hypothetical URL.
        return new Promise(function (resolve, reject) {
            fetch(url)
            .then(function (response){
                resolve(response.json());
            })
        })
    }

    function displayData(values, columnIndex, rowIndex) {
        // 针对 Excel 对象模型运行批处理操作
        Excel.run(function (ctx) {
            // 为活动工作表创建代理对象
            var sheet = ctx.workbook.worksheets.getActiveWorksheet();
            // 将向电子表格写入示例数据的命令插入队列
            // 有几个指标
            var factorNum = 3
            
            var colChar = numCharMap[columnIndex+1]
            var colNum = rowIndex+1
            var rowChar = numCharMap[columnIndex+factorNum]
            var rowNum = rowIndex + values.length
            console.log("colChar ",colChar, " colNum ", colNum, " rowChar ",rowChar ," rowNum ",rowNum)
            // sheet.getRange("B1:D"+rowNum).values = values;
            sheet.getRange(colChar+colNum+":"+rowChar+rowNum).values = values;
            // 运行排队的命令，并返回承诺表示任务完成
            return ctx.sync();
        })
        .catch(errorHandler);
    }

    function displaySelectedCells() {
        Office.context.document.getSelectedDataAsync(Office.CoercionType.Text,
            function (result) {
                if (result.status === Office.AsyncResultStatus.Succeeded) {
                    showNotification('选定的文本为:', '"' + result.value + '"');
                } else {
                    showNotification('错误', result.error.message);
                }
            });
    }

    // 处理错误的帮助程序函数
    function errorHandler(error) {
        // 请务必捕获 Excel.run 执行过程中出现的所有累积错误

        console.log("Error: " + error);
        if (error instanceof OfficeExtension.Error) {
            console.log("Debug info: " + JSON.stringify(error.debugInfo));
        }
    }

    // 用于显示通知的帮助程序函数
    function showNotification(header, content) {
        $("#notification-header").text(header);
        $("#notification-body").text(content);
        messageBanner.showBanner();
        messageBanner.toggleExpansion();
    }

    function webRequest() {
        let url = "https://rimeindex.com/login"; // This is a hypothetical URL.
        return new Promise(function (resolve, reject) {
            fetch(url)
                .then(function (response) {
                    return resolve(response);
                }).then(function (response) {
                    return reject(response.json());
                })
        })
    }
    var numCharMap = {
        0: "A",
        1: "B",
        2: "C",
        3: "D",
        4: "E",
        5: "F",
        6: "G",
        7: "H",
        8: "I",
        9: "J",
        10: "K",
        11: "L",
        12: "N",
        13: "M",
        14: "O",
        15: "P",
        16: "Q",
        17: "R",
        18: "S",
        19: "T",
        20: "U",
        21: "V",
        22: "W",
        23: "S",
        24: "Y",
        25: "Z",
    }
})();
