function main() {
  document.getElementById("submit").onclick = function() {
    getLunarBirthdayForYear()
  }
}

function getLunarBirthdayForYear() {
  const lunarBirthDate = document.getElementById("lunar-birth-date").value
  const isLeapMonth = document.getElementById("is-leap-month").checked
  const year = document.getElementById("year").value

  req = new XMLHttpRequest();
  req.onreadystatechange = function() {
    if (req.readyState === XMLHttpRequest.DONE) {
      if (req.status === 200) {
        const resp = JSON.parse(req.responseText)
        const date = new Date(resp.year, resp.month-1, resp.day)
        setBirthdayField(date.toLocaleDateString())
      }
    }
  }
  reqBody = {
    lunar_birth_date: (new Date(lunarBirthDate)).toISOString(),
    is_leap_month: isLeapMonth,
    year: parseInt(year),
  }
  req.open('POST', 'api/v1/lunar-birthday-for-year/')
  req.send(JSON.stringify(reqBody))
}

function setBirthdayField(birthday) {
  document.getElementById("birthday").value = birthday
}

main()
