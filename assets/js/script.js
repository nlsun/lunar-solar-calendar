function main() {
  setDefaultYear()
  document.getElementById("bd-submit").onclick = function() {
    getLunarBirthdayForYear()
  }
  document.getElementById("greg-submit").onclick = function() {
    getLunarBirthdayFromGregorian()
  }
  document.getElementById("cal-download").onclick = function() {
    getLunarBirthdayCalendar()
  }
}

function setDefaultYear() {
  document.getElementById("year").value = (new Date()).getFullYear()
}

function getLunarBirthdayForYear() {
  const lunarBirthDate = new Date(document.getElementById("bd-lunar-birth-date").value)
  const isLeapMonth = document.getElementById("bd-is-leap-month").checked
  const year = parseInt(document.getElementById("year").value)

  document.getElementById("bd-birthday").value = ""

  req = new XMLHttpRequest();
  req.onreadystatechange = function() {
    if (req.readyState === XMLHttpRequest.DONE) {
      if (req.status === 200) {
        const resp = JSON.parse(req.responseText)
        const date = new Date(resp.year, resp.month-1, resp.day)
        document.getElementById("bd-birthday").value = date.toLocaleDateString()
      }
    }
  }
  reqBody = {
    lunar_birth_date: lunarBirthDate.toISOString(),
    is_leap_month: isLeapMonth,
    year: year,
  }
  req.open('POST', 'api/v1/lunar-birthday-for-year/')
  req.send(JSON.stringify(reqBody))
}

function getLunarBirthdayFromGregorian() {
  const gregBirthDate = new Date(document.getElementById("greg-birth-date").value)

  document.getElementById("greg-birthday").value = ""
  document.getElementById("greg-is-leap-month").checked = false

  req = new XMLHttpRequest();
  req.onreadystatechange = function() {
    if (req.readyState === XMLHttpRequest.DONE) {
      if (req.status === 200) {
        const resp = JSON.parse(req.responseText)
        const date = new Date(resp.year, resp.month-1, resp.day)
        document.getElementById("greg-birthday").value = date.toLocaleDateString()
        if (resp.is_leap) {
          document.getElementById("greg-is-leap-month").checked = true
        }
      }
    }
  }
  reqBody = {
    solar_birth_date: gregBirthDate.toISOString(),
  }
  req.open('POST', 'api/v1/solar-to-lunar-birthday/')
  req.send(JSON.stringify(reqBody))
}

function getLunarBirthdayCalendar() {
  const personName = document.getElementById("person-name").value
  const lunarBirthDate = new Date(document.getElementById("cal-lunar-birth-date").value)
  const isLeapMonth = document.getElementById("cal-is-leap-month").checked
  const numYears = parseInt(document.getElementById("num-years").value)
  const notifications = JSON.parse(document.getElementById("notifications").value.trim())

  req = new XMLHttpRequest();
  req.onreadystatechange = function() {
    if (req.readyState === XMLHttpRequest.DONE) {
      if (req.status === 200) {
        const resp = JSON.parse(req.responseText)
        offerDownload(`${personName}_lunar_birthday.ics`, resp.calendar)
      } else {
        alert(`Generating calendar failed: ${req.responseText}`)
      }
    }
  }
  reqBody = {
    lunar_birth_date: lunarBirthDate.toISOString(),
    is_leap_month: isLeapMonth,
    last_year: lunarBirthDate.getFullYear() + numYears,
    title: `Birthday: ${personName}`,
    description: `Birth Date: ${lunarBirthDate.toLocaleDateString()}`,
    notifications: notifications,
  }
  req.open('POST', 'api/v1/lunar-birthday-calendar/')
  req.send(JSON.stringify(reqBody))
}

function offerDownload(filename, text) {
  var element = document.createElement('a');
  element.setAttribute('href', 'data:application/octet-stream;charset=utf-8,' + encodeURIComponent(text));
  element.setAttribute('download', filename);

  element.style.display = 'none';
  document.body.appendChild(element);

  element.click();

  document.body.removeChild(element);
}

main()
