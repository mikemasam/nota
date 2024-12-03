import { Database } from "bun:sqlite";
import { $ } from "bun";
import { parseArgs } from "util";
import os from "os";
import { parseDateTime } from "./date-parser";
const db: Database = new Database(`${os.homedir()}/.nota.db`, {
  create: true,
});
const args = Bun.argv.slice(2);
await migrate(db);
let action = args.shift();
if (action == "add" || action == "a" || action == "r") {
  const tag = args.shift();
  const last: string | undefined = args.pop();
  const [date, other] = parseDateTime(last);
  const str = [...args, other].join(" ");
  const query = db.query(
    `insert into reminds (tag, title, scheduled_at) values(?,?,${date == "now" ? "CURRENT_TIMESTAMP" : "?"})`,
  );
  if (date == "now") {
    await query.run(tag, str as any);
  } else {
    await query.run(tag, date ? (str as any) : `${str} ${last}`, date);
  }
  await printState();
} else if (action == "priority" || action == "prio" || action == "p") {
  const priority = args.shift();
  const csv_pos = args.join(",");
  if (!csv_pos) throw "invalid selection";
  let posList = csv_pos.split(",").filter((i) => !isNaN(parseInt(i)));
  const list = await load_reminders();
  for (let pos of posList) {
    if (!list[pos]) throw `${pos} not found`;
  }
  for (let pos of posList) {
    if (!list[pos]) continue;
    await db
      .query("update reminds set priority = ? where id = ?")
      .run(priority, list[pos].id);
  }
  await printState();
} else if (action == "later") {
  const last: string | undefined = args.pop();
  const [date, other] = parseDateTime(last);
  const csv_pos = [...args, other].join(",");
  if (!csv_pos) throw "invalid selection";
  let posList = csv_pos.split(",").filter((i) => !isNaN(parseInt(i)));
  const list = await load_reminders();
  for (let pos of posList) {
    if (!list[pos]) throw `${pos} not found`;
  }
  for (let pos of posList) {
    if (!list[pos]) continue;
    const item = list[pos];
    await db
      .query("update reminds set scheduled_at = ? where id = ?")
      .run(date, item.id);
    console.log(`${Bun.color("yellow", "ansi")} ${item.scheduled_at ?? '*'} -> ${date ?? '*'} ${Bun.color("gray", "ansi")}~ ${item.title}`);
  }
  await printState();
} else if (action == "pop" || action == "del" || action == "d") {
  const csv_pos = args.join(",");
  if (!csv_pos) throw "invalid selection";
  let posList = csv_pos.split(",").filter((i) => !isNaN(parseInt(i)));
  const list = await load_reminders();
  for (let pos of posList) {
    if (!list[pos]) throw `${pos} not found`;
  }
  for (let pos of posList) {
    if (!list[pos]) continue;
    await db
      .query("update reminds set deleted_at = CURRENT_TIMESTAMP where id = ?")
      .run(list[pos].id);
  }
  await printState();
} else if (action == "print" || action == "show") {
  const cond = args.shift();
  await printState(cond);
} else if (action == "version" || action == "v") {
  console.log(`${Bun.color("grey", "ansi")}version: v0.0.1`);
  console.log(`${Bun.color("grey", "ansi")}webpage: https://github.com/mikemasam/nota`);
} else if (action == "help") {
  console.log(`${Bun.color("grey", "ansi")}version: v0.0.1`);
  console.log(`${Bun.color("grey", "ansi")}webpage: https://github.com/mikemasam/nota`);
  console.log(`${Bun.color("grey", "ansi")}\t$ datetime formats: [2024-12-10+11:46/today/tomorrow+morning/1week/+2weeks]`);
  console.log(`${Bun.color("grey", "ansi")}\t$ nota add/a/r tag title datetime ~ add new note`);
  console.log(`${Bun.color("grey", "ansi")}\t$ nota later index       datetime ~ move note forward`);
  console.log(`${Bun.color("grey", "ansi")}\t$ nota del/pop index              ~ remove note`);
} else {
  await printState();
}
async function load_reminders(cond?: string) {
  const items = await db
    .query(
      `select *, (scheduled_at <= date('now', 'localtime')) as is_old from reminds where ${cond ? `${cond}` : "deleted_at is null "} order by priority desc, scheduled_at asc`,
    )
    .all();
  return items;
}

async function printState(cond?: string) {
  const reminds = await load_reminders(cond);
  if (!reminds?.length) console.log(`${Bun.color("grey", "ansi")}> Nothing to show`);
  for (let i = 0; i < reminds.length; i++) {
    const item: any = reminds[i];
    console.log(
      `${Bun.color("grey", "ansi")}${i}:${Bun.color(item.is_old ? "grey" : "pink", "ansi")}[${item.scheduled_at ?? "*"}]> ${Bun.color("grey", "ansi")}${item.tag}: ${Bun.color("white", "ansi")}${item.title} ${item.priority > 0 ? "*" : ""} ${Bun.color("#494949", "ansi")}[${item.created_at}]`,
    );
  }
}
async function migrate(db: Database) {
  const versions = [
    `CREATE TABLE IF NOT EXISTS reminds (
        id INTEGER PRIMARY KEY ASC,
        tag TEXT NOT NULL,
        title TEXT NOT NULL,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        deleted_at DATETIME NULL DEFAULT NULL,
        priority INTEGER NOT NULL DEFAULT 0,
        scheduled_at DATETIME NULL,
        period TEXT DEFAULT NULL,
        finished_at DATETIME DEFAULT NULL
    );`,
  ];
  //await db.query(`PRAGMA user_version = 0`).get();
  const out = await db.query<any, any>("PRAGMA user_version;").get();
  for (let i = Number(out.user_version); i < versions.length; i++) {
    const sql = versions[i];
    await db.query(sql).run();
    await db.query(`PRAGMA user_version = ${i + 1}`).get();
  }
  //console.log(out);
}
