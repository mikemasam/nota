package main

import (
	"database/sql"
	"fmt"
	"log"
	"nota/notalib"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func main() {
	db, err := sql.Open("sqlite3", notalib.ResolveHomeDir(".nota.db"))
	if err != nil {
		log.Fatal(err)
	}
	notalib.Migrate(db)

	action := ""
	// Print the arguments
	if len(os.Args) > 1 {
		action = os.Args[1]
	}
	switch action {
	case "add", "a":
		addRemind(db)
		printReminds(db)
	case "later", "l":
		laterRemind(db)
		printReminds(db)
	case "secret":
		secretRemind(db)
		printReminds(db)
	case "del", "d":
		deleteRemind(db)
		printReminds(db)
	case "deletehard":
		hardDeleteRemind(db)
		printReminds(db)
	case "help", "h", "version", "v", "--help", "-h", "--version", "-v":
		printHelp()
	default:
		printReminds(db)
	}
	defer db.Close()
}

func addRemind(db *sql.DB) {
	arr := os.Args[2:]
	if len(arr) < 2 {
		log.Fatal("[add] require tag and content")
	}
	remind := Remind{
		createdAt: time.Now(),
		tag:       arr[0],
	}
	last := arr[len(arr)-1]
	date, _ := notalib.ParseDateTime(last)
	start := 1
	end := len(arr)
	if date != nil {
		remind.scheduledAt = date
		end -= 1
	}
	remind.title = strings.Join(arr[start:end], " ")

	_, err := db.Exec("INSERT INTO reminds (tag, title, scheduled_at) VALUES(?,?,?);", remind.tag, remind.title, remind.scheduledAt)
	if err != nil {
		log.Fatal(err)
	}
}

func laterRemind(db *sql.DB) {
	arr := os.Args[2:]
	last := arr[len(arr)-1]
	date, _ := notalib.ParseDateTime(last)
	start := 0
	end := len(arr)
	if date != nil {
		end -= 1
	}

	reminds, err := loadReminds(db)
	if err != nil {
		log.Fatal(err)
	}
	posList := strings.Split(strings.Join(arr[start:end], ","), ",")
	for _, p := range posList {
		_, err := strconv.Atoi(p)
		if err != nil {
			log.Fatalf("Invalid index %s", p)
		}
	}
	for _, p := range posList {
		pos, _ := strconv.Atoi(p)
		if len(reminds) < pos {
			continue
		}
		remind := reminds[pos]
		_, err := db.Exec("update reminds set scheduled_at = ? where id = ?;", date, remind.id)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func deleteRemind(db *sql.DB) {
	arr := os.Args[2:]
	reminds, err := loadReminds(db)
	if err != nil {
		log.Fatal(err)
	}
	posList := strings.Split(arr[0], ",")
	for _, p := range posList {
		_, err := strconv.Atoi(p)
		if err != nil {
			log.Fatalf("Invalid index %s", p)
		}
	}
	for _, p := range posList {
		pos, _ := strconv.Atoi(p)
		if len(reminds) < pos {
			continue
		}
		remind := reminds[pos]
		_, err := db.Exec("update reminds set deleted_at = ? where id = ?;", time.Now(), remind.id)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func hardDeleteRemind(db *sql.DB) {
	arr := os.Args[2:]
	reminds, err := loadReminds(db)
	if err != nil {
		log.Fatal(err)
	}
	p := arr[0]
	pos, err := strconv.Atoi(p)
	if err != nil {
		log.Fatalf("Invalid index %s", p)
	}
	if len(reminds) < pos {
		return
	}
	remind := reminds[pos]
	_, err = db.Exec("delete from reminds where id = ?;", remind.id)
	if err != nil {
		log.Fatal(err)
	}
}

func secretRemind(db *sql.DB) {
	arr := os.Args[2:]
	reminds, err := loadReminds(db)
	if err != nil {
		log.Fatal(err)
	}
	posList := strings.Split(arr[0], ",")
	for _, p := range posList {
		_, err := strconv.Atoi(p)
		if err != nil {
			log.Fatalf("Invalid index %s", p)
		}
	}
	for _, p := range posList {
		pos, _ := strconv.Atoi(p)
		if len(reminds) < pos {
			continue
		}
		remind := reminds[pos]
		_, err := db.Exec("update reminds set priority = -1 where id = ?;", remind.id)
		if err != nil {
			log.Fatal(err)
		}
	}
}

type Remind struct {
	createdAt   time.Time
	scheduledAt *time.Time
	deletedAt   *time.Time
	is_old      *bool
	tag         string
	title       string
	id          int
}

func printReminds(db *sql.DB) {
	if notalib.FileExists(notalib.ResolveHomeDir(".silent")) {
		return
	}
	reminds, err := loadReminds(db)
	if err != nil {
		log.Fatal(err)
	}
	if len(reminds) == 0 {
		fmt.Printf("%s> %s\n", notalib.Color("0"), "Nothing to show")
	}

	showEverything := slices.Contains(os.Args, "+a")
	showCreatedAt := slices.Contains(os.Args, "+c") || showEverything
	showMore := slices.Contains(os.Args, "+") || showCreatedAt
	for i, p := range reminds {
		scheduledAt := "*"
		if p.scheduledAt != nil {
			if showMore {
				scheduledAt = p.scheduledAt.Format("2006-01-02 15:04")
			} else {
				scheduledAt = p.scheduledAt.Format("02/01 15:04")
			}
		}
		deletedAt := ""
		if p.deletedAt != nil {
			deletedAt = fmt.Sprintf("%sD[%s]", notalib.Color("218"), p.deletedAt.Format("2006-01-02 15:04"))
		}
		createdAt := ""
		if showCreatedAt {
			createdAt = fmt.Sprintf("[%s]", p.createdAt.Format("2006-01-02 15:04"))
		}

		fmt.Printf(
			"%s%d:%s[%s]%s %s: %s%s %s%s %s\n",
			notalib.Color("243"), i,
			notalib.Color("245"), scheduledAt,
			notalib.Color("248"), p.tag,
			notalib.Color("231"), p.title,
			notalib.Color("0"), createdAt,
			deletedAt,
		)
	}
}

func printHelp() {
	fmt.Printf(
		`
%sversion: v0.0.11
webpage: https://github.com/mikemasam/nota
? datetime formats: [2024-12-10+11:46/today/now/tomorrow+morning/1week/+2weeks]
$ nota add/a/r tag description datetime ~ add new note
$ nota later index       			 datetime ~ move note datetime
$ nota del index           							~ softdelete a note, use +a to all 
$ nota secret index           					~ hide a note, use +secret to show secrets
$ nota delelehard index           			~ delete a single note forever
$ nota .youtube 												~ list notes contain word youtube
$ nota +deleted 												~ list deleted notes
$ nota +a 															~ list all notes including deleted
$ nota +secret 													~ list secrets 
`,
		notalib.Color("248"),
	)
}

func loadReminds(db *sql.DB) ([]Remind, error) {
	var reminds []Remind
	builder := []string{}
	if slices.Contains(os.Args, "+secret") {
		builder = append(builder, `priority = -1`)
	} else {
		builder = append(builder, `priority != -1`)
	}
	if slices.Contains(os.Args, "+deleted") {
		builder = append(builder, `deleted_at is not null`)
	} else if slices.Contains(os.Args, "+a") {
	} else {
		builder = append(builder, `deleted_at is null`)
	}

	matchIdx := slices.IndexFunc(os.Args, func(s string) bool { return strings.HasPrefix(s, ".") })
	if matchIdx > -1 {
		builder = append(builder, fmt.Sprintf(`title like '%%%s%%'`, os.Args[matchIdx][1:]))
	}
	querySQL := fmt.Sprintf("select id, tag, title, scheduled_at, created_at, deleted_at, (scheduled_at <= date('now', 'localtime')) as is_old from reminds where %s order by scheduled_at asc", strings.Join(builder, " and "))
	rows, err := db.Query(querySQL)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var remind Remind
		err = rows.Scan(&remind.id, &remind.tag, &remind.title, &remind.scheduledAt, &remind.createdAt, &remind.deletedAt, &remind.is_old)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		reminds = append(reminds, remind)
	}
	// Check for errors during iteration
	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return reminds, nil
}
