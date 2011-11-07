package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()

    if flag.NArg() == 0 { missingCommand() }

    command := flag.Arg(0)
    args := flag.Args()[1:]
    if command == "help" { help() }

    switch command {
        case "mount": mount(args)
        case "add": add(args)
        case "remove": remove(args)
        case "tag": tag(args)
        case "untag": untag(args)
        case "tags" : tags(args)
        default: invalidCommand(command)
    }
}

// commands

func help() {
    showUsage()
    os.Exit(0)
}

func mount(args []string) {
    if (len(args) == 0) { log.Fatal("No mountpoint specified.") }

    mountPath := args[0]

    vfs, error := MountVfs(mountPath)
    if (error != nil) { log.Fatalf("Could not mount filesystem: %v", error.String()) }
    defer vfs.Unmount()

    log.Printf("Database '%v' mounted at '%v'.\n", databasePath(), mountPath)

    vfs.Loop()
}

func add(paths []string) {
    db, error := OpenDatabase(databasePath())
    if error != nil { log.Fatalf("Could not open database: %v", error.String()) }
    defer db.Close()

    for _, path := range paths {
        absPath, error := filepath.Abs(path)
        if error != nil { log.Fatalf("Could resolve path '%v': %v", path, error.String()) }

        fingerprint, error := Fingerprint(absPath)
        if error != nil { log.Fatalf("Could not fingerprint '%v': %v", path, error.String()) }

        file, error := db.File(fingerprint)
        if error != nil { log.Fatalf("Could not determine whether file '%v' exists in the database: %v", path, error.String()) }

        if file == nil {
            //file, error = db.AddFile(fingerprint)
        }

        //filePath := db.GetFilePath(absPath)
        //if filePath != nil continue

        //db.Add
    }
}

func remove(args []string) {
    log.Fatal("Not implemented.")
}

func tag(args []string) {
    log.Fatal("Not implemented.")
}

func untag(args []string) {
    log.Fatal("Not implemented.")
}

func tags(args []string) {
    db, error := OpenDatabase(databasePath())
    if error != nil { log.Fatalf("Could not open database: %v", error.String()) }
    defer db.Close()

    tags, error := db.Tags()
    if error != nil { log.Fatalf("Could not retrieve tags: %v", error.String()) }

    for _, tag := range tags {
        fmt.Println(tag.Name)
    }
}

// other stuff

func showUsage() {
    fmt.Println("usage: tmsu <command> [<args>]")
    fmt.Println()
    fmt.Println("Commands:")
    fmt.Println()
    fmt.Println(" help       list commands or provide help for a given command")
    fmt.Println(" add        add a file to the VFS without applying tags")
    fmt.Println(" tag        add a file (if necessary) and apply tags")
    fmt.Println(" tags       list all tags or tags for a given file")
}

func missingCommand() {
    log.Fatal("No command specified.")
}

func invalidCommand(command string) {
    log.Fatalf("No such command '%v'.", command)
}
