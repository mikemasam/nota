
const timers = ["morning", "afternoon", "evening"];
const dayers = ["tomorrow", "today"];
export function parseDateTime(input: string | undefined): [string | null, string] {
  if (!input) return [null, ""];

  const DEFAULT_TIMES = {
    morning: "08:00:00",
    afternoon: "12:00:00",
    evening: "18:00:00",
    default: "06:00:00",
  };

  const today = new Date();
  let date = new Date(today);

  let [dayPart, timePart] = input.replaceAll("+", " ").split(" ");

  // Handle relative date adjustments (e.g., +2days, +3weeks)
  if (/^\+?\d*(day|days|week|weeks)$/.test(input)) {
    const relativeDate = parseRelativeDate(input, date);
    if (relativeDate) date = relativeDate;
    else return [null, input];
  } else if (dayers.includes(dayPart)) {
    // Handle "today" or "tomorrow"
    if (dayPart === "tomorrow") date.setDate(today.getDate() + 1);
  } else if (timers.includes(dayPart)) {
    // If only a timer (e.g., "morning") is provided, assume "today"
    timePart = dayPart;
  } else if (isISODate(dayPart)) {
    return [constructISODate(dayPart, timePart), ""];
  } else {
    return [null, input]; // Invalid input
  }

  const dateStr = date.toISOString().split("T")[0];

  // Determine the time
  const time = timePart && timers.includes(timePart) ? DEFAULT_TIMES[timePart] : DEFAULT_TIMES.default;

  return [`${dateStr} ${time}`, ""];
}

function parseRelativeDate(input: string, baseDate: Date): Date | null {
  const relativeMatch = /^\+?(\d+)?(day|days|week|weeks)$/.exec(input);
  if (!relativeMatch) return null;

  const [, amount, unit] = relativeMatch;
  const value = parseInt(amount || "1", 10); // Default to 1 if no number provided

  const resultDate = new Date(baseDate);

  switch (unit) {
    case "day":
    case "days":
      resultDate.setDate(baseDate.getDate() + value);
      break;
    case "week":
    case "weeks":
      resultDate.setDate(baseDate.getDate() + value * 7);
      break;
    default:
      return null;
  }

  return resultDate;
}

// Helper to check if a string is a valid ISO date (YYYY-MM-DD)
function isISODate(str: string): boolean {
  return /^\d{4}-\d{2}-\d{2}$/.test(str);
}

// Helper to construct ISO date string
function constructISODate(date: string, time: string | undefined): string | null {
  if (!time) return `${date} 00:00:00`;
  if (/^\d{2}:\d{2}$/.test(time)) return `${date} ${time}:00`;
  if (/^\d{2}:\d{2}:\d{2}$/.test(time)) return `${date} ${time}`;
  return null; // Invalid time format
}
