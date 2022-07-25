window.addEventListener('DOMContentLoaded', (event) => {
    document.querySelector('#tbl-companies').querySelectorAll('td').forEach(item => {
        item.addEventListener('click', event => {
            let company = event.target.attributes['data-company'].value;
            showLeagues(company);
        })
    });
    document.getElementById('btn-back-to-companies').addEventListener('click', backToCompaniesClick);
    document.getElementById('btn-back-to-leagues').addEventListener('click', backToLeaguesClick);
    document.getElementById('btn-back-to-teams').addEventListener('click', backToTeamsClick);
});

var currentScheduleData = {};

var showLeagues = function (company) {
    showLoading()
    fetch('getleagues?company=' + company)
    .then(function(resp) {
        return resp.json();
    })
    .then(function(json) {
        var tblBody = document.querySelector('#tbl-leagues').querySelector('tbody');
        tblBody.innerHTML = '';

        for (var i = 0; i < json.length; i++) {
            var row = document.createElement('tr');
            var cell = document.createElement('td');
            cell.setAttribute('data-company', company);
            cell.setAttribute('data-leagueId', json[i].Id);
            cell.innerHTML = json[i].Name;
            row.appendChild(cell);
            tblBody.appendChild(row);
            cell.addEventListener('click', event => {
                let company = event.target.attributes['data-company'].value;
                let leagueId = event.target.attributes['data-leagueId'].value;

                showTeams(company, leagueId);
            })
        }

        hideLoading();
        document.getElementById('tbl-companies').style.display = 'none';
        document.getElementById('tbl-teams').style.display = 'none';
        document.getElementById('tbl-leagues').style.display = 'block';
    });
};

var showTeams = function (company, leagueId) {
    showLoading()
    fetch('getteams?company=' + company + '&leagueId=' + leagueId)
    .then(function(resp) {
        return resp.json();
    })
    .then(function(json) {
        var tblBody = document.querySelector('#tbl-teams').querySelector('tbody');
        tblBody.innerHTML = '';

        for (var i = 0; i < json.length; i++) {
            var row = document.createElement('tr');
            var cell = document.createElement('td');
            cell.setAttribute('data-company', company);
            cell.setAttribute('data-leagueId', leagueId);
            cell.setAttribute('data-teamId', json[i].Id);
            cell.setAttribute('data-teamName', json[i].Name);
            cell.innerHTML = json[i].Name;
            row.appendChild(cell);
            tblBody.appendChild(row);
            cell.addEventListener('click', event => {
                let company = event.target.attributes['data-company'].value;
                let leagueId = event.target.attributes['data-leagueId'].value;
                let teamId = event.target.attributes['data-teamId'].value;
                let teamName = event.target.attributes['data-teamName'].value;

                showSchedule(company, leagueId, teamId, teamName);
            })
        }

        hideLoading();
        document.getElementById('tbl-schedule').style.display = 'none';
        document.getElementById('tbl-teams').style.display = 'block';
        document.getElementById('tbl-leagues').style.display = 'none';
    });
};

var showSchedule = function (company, leagueId, teamId, teamName) {
    showLoading()
    fetch('getschedule?company=' + company + '&leagueId=' + leagueId + '&teamId=' + teamId)
    .then(function(resp) {
        return resp.json();
    })
    .then(function(json) {
        currentScheduleData = json;
        document.getElementById('schedule-header').innerHTML = "Schedule for " + teamName;

        var tblBody = document.querySelector('#tbl-schedule').querySelector('tbody');
        tblBody.innerHTML = '';        

        var btnExportCSV = document.getElementById('btn-export-csv');      
        btnExportCSV.setAttribute('data-company', company);
        btnExportCSV.setAttribute('data-leagueId', leagueId);
        btnExportCSV.setAttribute('data-teamId', teamId);
        btnExportCSV.setAttribute('data-teamName', teamName);

        for (var i = 0; i < json.length; i++) {
            var row = document.createElement('tr');
            var cell = document.createElement('td');
            cell.innerHTML = json[i].HomeTeam;
            row.appendChild(cell);
            cell = document.createElement('td');
            cell.innerHTML = json[i].VisitorTeam;
            row.appendChild(cell);
            cell = document.createElement('td');
            cell.innerHTML = json[i].Location;
            row.appendChild(cell);
            cell = document.createElement('td');
            cell.innerHTML = (new Date(json[i].EventTime.substr(0, json[i].EventTime.length - 1))).toLocaleString('en-US', {year: 'numeric', month: 'numeric', day: 'numeric', hour: '2-digit', minute: '2-digit'});
            row.appendChild(cell);
            tblBody.appendChild(row);
        }

        hideLoading();
        document.getElementById('tbl-teams').style.display = 'none';
        document.getElementById('tbl-schedule').style.display = 'block';

        btnExportCSV.removeEventListener('click', exportCSVOnClick);
        btnExportCSV.addEventListener('click', exportCSVOnClick);
    });
};

var exportCSVOnClick = function (event) {    
    let teamName = event.target.attributes['data-teamName'].value;

    // Building the CSV from the Data two-dimensional array
    // Each column is separated by ";" and new line "\n" for next row
    var csvContent = templateStart;
    var json = currentScheduleData;
    for (var i = 0; i < json.length; i++) {
        let dt = new Date(json[i].EventTime.substr(0, json[i].EventTime.length - 1));
        csvContent += "GAME,REGULAR,," + json[i].HomeTeam + "," + json[i].VisitorTeam + "," + dt.toLocaleDateString() + "," + dt.toLocaleString('en-US', {hour: '2-digit', minute: '2-digit'}) + ',1:30,' + json[i].Location + ',,\r\n';
    };
    csvContent += templateEnd;

    // The download function takes a CSV string, the filename and mimeType as parameters
    // Scroll/look down at the bottom of this snippet to see how download is called
    var download = function (content, fileName, mimeType) {
        var a = document.createElement('a');
        mimeType = mimeType || 'application/octet-stream';

        if (navigator.msSaveBlob) { // IE10
            navigator.msSaveBlob(new Blob([content], {
                type: mimeType
            }), fileName);
        } else if (URL && 'download' in a) { //html5 A[download]
            a.href = URL.createObjectURL(new Blob([content], {
                type: mimeType
            }));
            a.setAttribute('download', fileName);
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
        } else {
            location.href = 'data:application/octet-stream,' + encodeURIComponent(content); // only this mime type is supported
        }
    }

    download(csvContent, teamName + '-Schedule.csv', 'text/csv;encoding:utf-8');
};

var backToCompaniesClick = function() {
    document.getElementById('tbl-leagues').style.display = 'none';
    document.getElementById('tbl-companies').style.display = 'block';
};

var backToLeaguesClick = function() {
    document.getElementById('tbl-teams').style.display = 'none';
    document.getElementById('tbl-leagues').style.display = 'block';
};

var backToTeamsClick = function() {
    document.getElementById('tbl-teams').style.display = 'block';
    document.getElementById('tbl-schedule').style.display = 'none';
};

var showLoading = function() {
    document.getElementById('loading-indicator').style.display = 'block';
};

var hideLoading = function() {
    document.getElementById('loading-indicator').style.display = 'none';
};

const templateStart = "Schedule Template,,,,,,,,,,\r\n \
,,,,,,,,,,\r\n \
Type,Game Type,Title (Optional),Home,Away,Date,Time,Duration,Location (Optional),Address (Optional),Notes (Optional)\r\n";

const templateEnd = "\"Options: \
GAME\
SCRIMMAGE\
DROP-IN\
PRACTICE\
EVENT\",\"Options: \
PRE-SEASON\
REGULAR\
PLAYOFF\
TOURNAMENT\
(Only required for games)\",\"Example:\
Team BBQ\
(used for events only)\",\"Please Note:\
Name must be spelt exactly as in BenchApp for matching to work properly. If no match is found, a new team will be created and linked to this event.\
(Your team must ALWAYS be home or away)\",\"Please Note:\
Name must be spelt exactly as in BenchApp for matching to work properly. If no match is found, a new team will be created and linked to this event.\
(Your team must ALWAYS be home or away)\",\"Must be this format:\
DD/MM/YYYY\",\"Must be this format:\
 6:30 PM\",\"Duration of the event\
H:MM\",\"The name of the facility:\
ex. Rink name, etc\
(Optional)\",\"The Address of the facility:\
Street address with zip/postal code\
(Optional)\",\"Notes visible to your players\
(Optional)\",,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,";