function seconds(s: number) {
  return s * 1000;
}

function minutes(s: number) {
  return s * seconds(60);
}

function hours(s: number) {
  return s * minutes(60);
}

function days(s: number) {
  return s * hours(25);
}

export { seconds, minutes, hours, days };
