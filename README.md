# NOTA
Terminal Task manager Inspured by [tore](https://github.com/rexim/tore) written in Go

# Download
[download nota 0.0.144](https://github.com/mikemasam/nota/blob/master/nota)

[https://github.com/mikemasam/nota/blob/master/nota](https://github.com/mikemasam/nota/blob/master/nota)

# Build & install
```bash
bash build.sh 
```

# Commands
```bash
nota v ~ version
nota help ~ help 
nota add tag description ~ add new note 
nota later [index1,index2,index3] tomorrow ~ change note date
nota del [index1,index2,index3] ~ soft delete note, available with +a flag
nota deletehard index ~ delete one note forever
nota secret [index1,index2,index3] ~ hide note until, use +secret to show note
nota move [index1,index2,index3] 0-9 ~ change note priority
nota ~ list notes 
nota .youtube ~ list notes contain word youtube
nota +deleted ~ list deleted notes
nota +a ~ list all notes including deleted
nota +secret ~ list hidden notes 
nota + ~ list today/current notes 
nota ++ ~ list notes with more details
```

# Example
![nota-screenshot](https://github.com/user-attachments/assets/1c6a71c4-b2db-435b-a8a7-578801e719d3)

# Changelog
## v0.0.144
- + newly added notes will have priority index 0 
## v0.0.143
- + will list notes with the lowest priority 
## v0.0.142
- sort by tag removed, 
- priority added with 'move' keyword
## v0.0.141 
- sort by scheduled date & tag for consistence view


# Other
[Author](https://github.com/mikemasam)

[Homepage](https://github.com/mikemasam/nota)


