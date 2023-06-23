export function isHomeBrokerClosed() {
  const currentDate = new Date()
  const closedDate = new Date(
    currentDate.getFullYear(),
    currentDate.getMonth(),
    currentDate.getDate(),
    18,
    0,
    0
  )

  return currentDate > closedDate
}
