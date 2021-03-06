/*
Copyright 2011-2014 Paul Ruane.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package vfs

import (
	"fmt"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
	"tmsu/common/log"
	"tmsu/entities"
	"tmsu/query"
	"tmsu/storage"
)

const tagsDir = "tags"
const queriesDir = "queries"
const queryHelpFilename = "README.md"
const queryDirHelp = `Query Directories
-----------------

Navigate to any directory that is a valid query to see a view of the files that
match the query:

    $ ls
    README.md
    $ ls "cheese and wine"
    pinot_cheddar.12  edam_blanc.14
    $ ls "cheese and (tomato or mushroom)
    margherita.7  funghi.11
    $ ls
    cheese and (tomato or mushroom)  cheese and wine 

Query directories are saved automatically and can be removed with ` + "`rmdir`."

type FuseVfs struct {
	store     *storage.Storage
	mountPath string
	server    *fuse.Server
}

func MountVfs(databasePath string, mountPath string, options []string) (*FuseVfs, error) {
	fuseVfs := FuseVfs{}
	pathFs := pathfs.NewPathNodeFs(&fuseVfs, nil)
	conn := nodefs.NewFileSystemConnector(pathFs.Root(), nil)
	mountOptions := &fuse.MountOptions{Options: options}

	server, err := fuse.NewServer(conn.RawFS(), mountPath, mountOptions)
	if err != nil {
		return nil, fmt.Errorf("could not mount virtual filesystem at '%v': %v", mountPath, err)
	}

	store, err := storage.OpenAt(databasePath)
	if err != nil {
		return nil, fmt.Errorf("could not open database at '%v': %v", databasePath, err)
	}

	fuseVfs.store = store
	fuseVfs.mountPath = mountPath
	fuseVfs.server = server

	return &fuseVfs, nil
}

func (vfs FuseVfs) Unmount() {
	vfs.server.Unmount()
}

func (vfs FuseVfs) Serve() {
	vfs.server.Serve()
}

func (vfs FuseVfs) SetDebug(debug bool) {
	vfs.SetDebug(debug)
}

func (vfs FuseVfs) Access(name string, mode uint32, context *fuse.Context) fuse.Status {
	log.Infof(2, "BEGIN Access(%v, %v)", name, mode)
	defer log.Infof(2, "END Access(%v, %v)", name, mode)

	return fuse.ENOSYS
}

func (vfs FuseVfs) Chmod(name string, mode uint32, context *fuse.Context) fuse.Status {
	log.Infof(2, "BEGIN Chmod(%v, %v)", name, mode)
	defer log.Infof(2, "BEGIN Chmod(%v, %v)", name, mode)

	return fuse.ENOSYS
}

func (vfs FuseVfs) Chown(name string, uid uint32, gid uint32, context *fuse.Context) fuse.Status {
	log.Infof(2, "BEGIN Chown(%v, %v, %v)", name, uid, gid)
	defer log.Infof(2, "BEGIN Chown(%v, %v)", name, uid, gid)

	return fuse.ENOSYS
}

func (vfs FuseVfs) Create(name string, flags uint32, mode uint32, context *fuse.Context) (nodefs.File, fuse.Status) {
	log.Infof(2, "BEGIN Create(%v, %v, %v)", name, flags, mode)
	defer log.Infof(2, "BEGIN Create(%v, %v)", name, flags, mode)

	return nil, fuse.ENOSYS
}

func (vfs FuseVfs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	log.Infof(2, "BEGIN GetAttr(%v)", name)
	defer log.Infof(2, "END GetAttr(%v)", name)

	switch name {
	case "":
		fallthrough
	case tagsDir:
		return vfs.getTagsAttr()
	case queriesDir:
		return vfs.getQueryAttr()
	}

	path := vfs.splitPath(name)

	switch path[0] {
	case tagsDir:
		return vfs.getTaggedEntryAttr(path[1:])
	case queriesDir:
		return vfs.getQueryEntryAttr(path[1:])
	}

	return nil, fuse.ENOENT
}

func (vfs FuseVfs) GetXAttr(name string, attr string, context *fuse.Context) ([]byte, fuse.Status) {
	log.Infof(2, "BEGIN GetXAttr(%v, %v)", name, attr)
	defer log.Infof(2, "END GetAttr(%v, %v)", name, attr)

	return nil, fuse.ENOSYS
}

func (vfs FuseVfs) Link(oldName string, newName string, context *fuse.Context) fuse.Status {
	log.Infof(2, "BEGIN Link(%v, %v)", oldName, newName)
	defer log.Infof(2, "END Link(%v, %v)", oldName, newName)

	return fuse.ENOSYS
}

func (vfs FuseVfs) ListXAttr(name string, context *fuse.Context) ([]string, fuse.Status) {
	log.Infof(2, "BEGIN ListXAttr(%v)", name)
	defer log.Infof(2, "END ListXAttr(%v)", name)

	return nil, fuse.ENOSYS
}

func (vfs FuseVfs) Mkdir(name string, mode uint32, context *fuse.Context) fuse.Status {
	log.Infof(2, "BEGIN Mkdir(%v)", name)
	defer log.Infof(2, "END Mkdir(%v)", name)

	path := vfs.splitPath(name)

	if len(path) != 2 {
		return fuse.EPERM
	}

	switch path[0] {
	case tagsDir:
		name := path[1]

		if _, err := vfs.store.AddTag(name); err != nil {
			log.Fatalf("could not create tag '%v': %v", name, err)
		}

		if err := vfs.store.Commit(); err != nil {
			log.Fatalf("could not commit transaction: %v", err)
		}

		return fuse.OK
	case queriesDir:
		return fuse.EINVAL
	}

	return fuse.ENOSYS
}

func (vfs FuseVfs) Mknod(name string, mode uint32, dev uint32, context *fuse.Context) fuse.Status {
	log.Infof(2, "BEGIN Mknod(%v)", name)
	defer log.Infof(2, "END Mknod(%v)", name)

	return fuse.ENOSYS
}

func (vfs FuseVfs) OnMount(nodeFs *pathfs.PathNodeFs) {
	log.Infof(2, "BEGIN OnMount()")
	defer log.Infof(2, "END OnMount()")
}

func (vfs FuseVfs) OnUnmount() {
	log.Infof(2, "BEGIN OnUnmount()")
	defer log.Infof(2, "END OnUnmount()")
}

func (vfs FuseVfs) Open(name string, flags uint32, context *fuse.Context) (nodefs.File, fuse.Status) {
	log.Infof(2, "BEGIN Open(%v)", name)
	defer log.Infof(2, "END Open(%v)", name)

	if name == filepath.Join(queriesDir, queryHelpFilename) {
		return nodefs.NewDataFile([]byte(queryDirHelp)), fuse.OK
	}

	return nil, fuse.ENOSYS
}

func (vfs FuseVfs) OpenDir(name string, context *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	log.Infof(2, "BEGIN OpenDir(%v)", name)
	defer log.Infof(2, "END OpenDir(%v)", name)

	switch name {
	case "":
		return vfs.topDirectories()
	case tagsDir:
		return vfs.tagDirectories()
	case queriesDir:
		return vfs.queriesDirectories()
	}

	path := vfs.splitPath(name)
	switch path[0] {
	case tagsDir:
		return vfs.openTaggedEntryDir(path[1:])
	case queriesDir:
		return vfs.openQueryEntryDir(path[1:])
	}

	return nil, fuse.ENOENT
}

func (vfs FuseVfs) Readlink(name string, context *fuse.Context) (string, fuse.Status) {
	log.Infof(2, "BEGIN Readlink(%v)", name)
	defer log.Infof(2, "END Readlink(%v)", name)

	path := vfs.splitPath(name)
	switch path[0] {
	case tagsDir, queriesDir:
		return vfs.readTaggedEntryLink(path[1:])
	}

	return "", fuse.ENOENT
}

func (vfs FuseVfs) RemoveXAttr(name string, attr string, context *fuse.Context) fuse.Status {
	log.Infof(2, "BEGIN RemoveXAttr(%v, %v)", name, attr)
	defer log.Infof(2, "END RemoveXAttr(%v, %v)", name, attr)

	return fuse.ENOSYS
}

func (vfs FuseVfs) Rename(oldName string, newName string, context *fuse.Context) fuse.Status {
	log.Infof(2, "BEGIN Rename(%v, %v)", oldName, newName)
	defer log.Infof(2, "END Rename(%v, %v)", oldName, newName)

	oldPath := vfs.splitPath(oldName)
	newPath := vfs.splitPath(newName)

	if len(oldPath) != 2 || len(newPath) != 2 {
		return fuse.EPERM
	}

	if oldPath[0] != tagsDir || newPath[0] != tagsDir {
		return fuse.EPERM
	}

	oldTagName := oldPath[1]
	newTagName := newPath[1]

	tag, err := vfs.store.TagByName(oldTagName)
	if err != nil {
		log.Fatalf("could not retrieve tag '%v': %v", oldTagName, err)
	}
	if tag == nil {
		return fuse.ENOENT
	}

	if _, err := vfs.store.RenameTag(tag.Id, newTagName); err != nil {
		log.Fatalf("could not rename tag '%v' to '%v': %v", oldTagName, newTagName, err)
	}

	if err := vfs.store.Commit(); err != nil {
		log.Fatalf("could not commit transaction: %v", err)
	}

	return fuse.OK
}

func (vfs FuseVfs) Rmdir(name string, context *fuse.Context) fuse.Status {
	log.Infof(2, "BEGIN Rmdir(%v)", name)
	defer log.Infof(2, "END Rmdir(%v)", name)

	path := vfs.splitPath(name)

	switch path[0] {
	case tagsDir:
		if len(path) != 2 {
			// can only remove top-level tag directories
			return fuse.EPERM
		}

		tagName := path[1]
		tag, err := vfs.store.TagByName(tagName)
		if err != nil {
			log.Fatalf("could not retrieve tag '%v': %v", tagName, err)
		}
		if tag == nil {
			return fuse.ENOENT
		}

		count, err := vfs.store.FileTagCountByTagId(tag.Id)
		if err != nil {
			log.Fatalf("could not retrieve file-tag count for tag '%v': %v", tagName, err)
		}
		if count > 0 {
			return fuse.Status(syscall.ENOTEMPTY)
		}

		if err := vfs.store.DeleteTag(tag.Id); err != nil {
			log.Fatalf("could not delete tag '%v': %v", tagName, err)
		}

		if err := vfs.store.Commit(); err != nil {
			log.Fatalf("could not commit transaction: %v", err)
		}

		return fuse.OK
	case queriesDir:
		if len(path) != 2 {
			// can only remove top-level queries directories
			return fuse.EPERM
		}

		text := path[1]

		if err := vfs.store.DeleteQuery(text); err != nil {
			log.Fatalf("could not remove tag '%v': %v", name, err)
		}

		if err := vfs.store.Commit(); err != nil {
			log.Fatalf("could not commit transaction: %v", err)
		}

		return fuse.OK
	}

	return fuse.ENOSYS
}

func (vfs FuseVfs) SetXAttr(name string, attr string, data []byte, flags int, context *fuse.Context) fuse.Status {
	log.Infof(2, "BEGIN SetXAttr(%v, %v)", name, attr)
	defer log.Infof(2, "END SetXAttr(%v, %v)", name, attr)

	return fuse.ENOSYS
}

func (vfs FuseVfs) StatFs(name string) *fuse.StatfsOut {
	log.Infof(2, "BEGIN StatFs(%v)", name)
	defer log.Infof(2, "END StatFs(%v)", name)

	return &fuse.StatfsOut{}
}

func (vfs FuseVfs) String() string {
	return "tmsu"
}

func (vfs FuseVfs) Symlink(value string, linkName string, context *fuse.Context) fuse.Status {
	log.Infof(2, "BEGIN Symlink(%v, %v)", value, linkName)
	defer log.Infof(2, "END Symlink(%v, %v)", value, linkName)

	return fuse.ENOSYS
}

func (vfs FuseVfs) Truncate(name string, offset uint64, context *fuse.Context) fuse.Status {
	log.Infof(2, "BEGIN Truncate(%v)", name)
	defer log.Infof(2, "END Truncate(%v)", name)

	return fuse.ENOSYS
}

func (vfs FuseVfs) Unlink(name string, context *fuse.Context) fuse.Status {
	log.Infof(2, "BEGIN Unlink(%v)", name)
	defer log.Infof(2, "END Unlink(%v)", name)

	fileId := vfs.parseFileId(name)
	if fileId == 0 {
		// can only unlink file symbolic links
		return fuse.EPERM
	}

	file, err := vfs.store.File(fileId)
	if err != nil {
		log.Fatal("could not retrieve file '%v': %v", fileId, err)
	}
	if file == nil {
		// reply ok if file doesn't exist otherwise recursive deletes fail
		return fuse.OK
	}
	path := vfs.splitPath(name)

	switch path[0] {
	case tagsDir:
		tagName := path[len(path)-2]
		//TODO value name

		tag, err := vfs.store.TagByName(tagName)
		if err != nil {
			log.Fatal(err)
		}
		if tag == nil {
			log.Fatalf("could not retrieve tag '%v'.", tagName)
		}

		if err = vfs.store.DeleteFileTag(fileId, tag.Id, 0); err != nil {
			log.Fatal(err)
		}

		if err := vfs.store.Commit(); err != nil {
			log.Fatalf("could not commit transaction: %v", err)
		}

		return fuse.OK
	case queriesDir:
		return fuse.EPERM
	}

	return fuse.ENOSYS
}

func (vfs FuseVfs) Utimens(name string, Atime *time.Time, Mtime *time.Time, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}

// non-exported

func (vfs FuseVfs) splitPath(path string) []string {
	return strings.Split(path, string(filepath.Separator))
}

func (vfs FuseVfs) parseFileId(name string) uint {
	parts := strings.Split(name, ".")
	count := len(parts)

	if count == 1 {
		return 0
	}

	id, err := atoui(parts[count-2])
	if err != nil {
		id, err = atoui(parts[count-1])
		if err != nil {
			return 0
		}
	}

	return id
}

func (vfs FuseVfs) topDirectories() ([]fuse.DirEntry, fuse.Status) {
	log.Infof(2, "BEGIN topDirectories")
	defer log.Infof(2, "END topDirectories")

	entries := []fuse.DirEntry{fuse.DirEntry{Name: tagsDir, Mode: fuse.S_IFDIR},
		fuse.DirEntry{Name: queriesDir, Mode: fuse.S_IFDIR}}
	return entries, fuse.OK
}

func (vfs FuseVfs) tagDirectories() ([]fuse.DirEntry, fuse.Status) {
	log.Infof(2, "BEGIN tagDirectories")
	defer log.Infof(2, "END tagDirectories")

	tags, err := vfs.store.Tags()
	if err != nil {
		log.Fatalf("Could not retrieve tags: %v", err)
	}

	entries := make([]fuse.DirEntry, len(tags))
	for index, tag := range tags {
		entries[index] = fuse.DirEntry{Name: tag.Name, Mode: fuse.S_IFDIR}
	}

	return entries, fuse.OK
}

func (vfs FuseVfs) queriesDirectories() ([]fuse.DirEntry, fuse.Status) {
	log.Infof(2, "BEGIN queriesDirectories")
	defer log.Infof(2, "END queriesDirectories")

	queries, err := vfs.store.Queries()
	if err != nil {
		log.Fatalf("could not retrieve queries: %v", err)
	}

	if len(queries) == 0 {
		return []fuse.DirEntry{fuse.DirEntry{Name: queryHelpFilename, Mode: fuse.S_IFREG}}, fuse.OK
	}

	entries := make([]fuse.DirEntry, len(queries))
	for index, query := range queries {
		entries[index] = fuse.DirEntry{Name: query.Text, Mode: fuse.S_IFDIR}
	}

	return entries, fuse.OK
}

func (vfs FuseVfs) getTagsAttr() (*fuse.Attr, fuse.Status) {
	log.Infof(2, "BEGIN getTagsAttr")
	defer log.Infof(2, "END getTagsAttr")

	tagCount, err := vfs.store.TagCount()
	if err != nil {
		log.Fatalf("could not get tag count: %v", err)
	}

	now := time.Now()
	return &fuse.Attr{Mode: fuse.S_IFDIR | 0755, Nlink: 2, Size: uint64(tagCount), Mtime: uint64(now.Unix()), Mtimensec: uint32(now.Nanosecond())}, fuse.OK
}

func (vfs FuseVfs) getQueryAttr() (*fuse.Attr, fuse.Status) {
	log.Infof(2, "BEGIN getQueryAttr")
	defer log.Infof(2, "END getQueryAttr")

	now := time.Now()
	return &fuse.Attr{Mode: fuse.S_IFDIR | 0755, Nlink: 2, Size: 0, Mtime: uint64(now.Unix()), Mtimensec: uint32(now.Nanosecond())}, fuse.OK
}

func (vfs FuseVfs) getTaggedEntryAttr(path []string) (*fuse.Attr, fuse.Status) {
	log.Infof(2, "BEGIN getTaggedEntryAttr(%v)", path)
	defer log.Infof(2, "END getTaggedEntryAttr(%v)", path)

	pathLength := len(path)
	name := path[pathLength-1]

	fileId := vfs.parseFileId(name)
	if fileId != 0 {
		return vfs.getFileEntryAttr(fileId)
	}

	// tag directory
	tagIds, err := vfs.tagNamesToIds(path)
	if err != nil {
		log.Fatalf("could not lookup tag IDs: %v.", err)
	}
	if tagIds == nil {
		return nil, fuse.ENOENT
	}

	now := time.Now()
	return &fuse.Attr{Mode: fuse.S_IFDIR | 0755, Nlink: 2, Size: uint64(0), Mtime: uint64(now.Unix()), Mtimensec: uint32(now.Nanosecond())}, fuse.OK
}

func (vfs FuseVfs) getQueryEntryAttr(path []string) (*fuse.Attr, fuse.Status) {
	log.Infof(2, "BEGIN getQueryEntryAttr(%v)", path)
	defer log.Infof(2, "END getQueryEntryAttr(%v)", path)

	pathLength := len(path)
	name := path[pathLength-1]

	if len(path) == 1 && path[0] == queryHelpFilename {
		now := time.Now()
		return &fuse.Attr{Mode: fuse.S_IFREG | 0444, Nlink: 1, Size: uint64(len(queryDirHelp)), Mtime: uint64(now.Unix()), Mtimensec: uint32(now.Nanosecond())}, fuse.OK
	}

	if len(path) > 1 {
		fileId := vfs.parseFileId(name)
		if fileId != 0 {
			return vfs.getFileEntryAttr(fileId)
		}

		return nil, fuse.ENOENT
	}

	queryText := path[0]

	if queryText[len(queryText)-1] == ' ' {
		// prevent multiple entries for same query when typing path in a GUI
		return nil, fuse.ENOENT
	}

	expression, err := query.Parse(queryText)
	if err != nil {
		return nil, fuse.ENOENT
	}

	tagNames := query.TagNames(expression)
	tags, err := vfs.store.TagsByNames(tagNames)
	for _, tagName := range tagNames {
		if !containsTag(tags, tagName) {
			return nil, fuse.ENOENT
		}
	}

	_, _ = vfs.store.AddQuery(queryText)

	if err := vfs.store.Commit(); err != nil {
		log.Fatalf("could not commit transaction: %v", err)
	}

	now := time.Now()
	return &fuse.Attr{Mode: fuse.S_IFDIR | 0755, Nlink: 2, Size: uint64(0), Mtime: uint64(now.Unix()), Mtimensec: uint32(now.Nanosecond())}, fuse.OK
}

func (vfs FuseVfs) getFileEntryAttr(fileId uint) (*fuse.Attr, fuse.Status) {
	file, err := vfs.store.File(fileId)
	if err != nil {
		log.Fatalf("could not retrieve file #%v: %v", fileId, err)
	}
	if file == nil {
		return &fuse.Attr{Mode: fuse.S_IFREG}, fuse.ENOENT
	}

	fileInfo, err := os.Stat(file.Path())
	var size int64
	var modTime time.Time
	if err == nil {
		size = fileInfo.Size()
		modTime = fileInfo.ModTime()
	} else {
		size = 0
		modTime = time.Time{}
	}

	return &fuse.Attr{Mode: fuse.S_IFLNK | 0755, Size: uint64(size), Mtime: uint64(modTime.Unix()), Mtimensec: uint32(modTime.Nanosecond())}, fuse.OK
}

func (vfs FuseVfs) openTaggedEntryDir(path []string) ([]fuse.DirEntry, fuse.Status) {
	log.Infof(2, "BEGIN openTaggedEntryDir(%v)", path)
	defer log.Infof(2, "END openTaggedEntryDir(%v)", path)

	expression := query.HasAll(path)
	files, err := vfs.store.QueryFiles(expression)
	if err != nil {
		log.Fatalf("could not query files: %v", err)
	}

	tagNames := make(map[string]interface{}, len(path))
	for _, tagName := range path {
		tagNames[tagName] = nil
	}

	furtherTagNames := make([]string, 0, 10)
	for _, file := range files {
		fileTags, err := vfs.store.FileTagsByFileId(file.Id)
		if err != nil {
			log.Fatalf("could not retrieve file-tags for file '%v': %v", file.Id, err)
		}

		tagIds := make([]uint, len(fileTags))
		for index, fileTag := range fileTags {
			tagIds[index] = fileTag.TagId
		}

		tags, err := vfs.store.TagsByIds(tagIds)
		if err != nil {
			log.Fatalf("could not retrieve tags: %v", err)
		}

		for _, tag := range tags {
			_, has := tagNames[tag.Name]
			if !has {
				if !containsName(furtherTagNames, tag.Name) {
					furtherTagNames = append(furtherTagNames, tag.Name)
				}
			}
		}
	}

	entries := make([]fuse.DirEntry, 0, len(files)+len(furtherTagNames))
	for _, tagName := range furtherTagNames {
		entries = append(entries, fuse.DirEntry{Name: tagName, Mode: fuse.S_IFDIR | 0755})
	}
	for _, file := range files {
		linkName := vfs.getLinkName(file)
		entries = append(entries, fuse.DirEntry{Name: linkName, Mode: fuse.S_IFLNK})
	}

	return entries, fuse.OK
}

func (vfs FuseVfs) openQueryEntryDir(path []string) ([]fuse.DirEntry, fuse.Status) {
	log.Infof(2, "BEGIN openQueryEntryDir(%v)", path)
	defer log.Infof(2, "END openQueryEntryDir(%v)", path)

	queryText := path[0]

	expression, err := query.Parse(queryText)
	if err != nil {
		return nil, fuse.ENOENT
	}

	tagNames := query.TagNames(expression)
	tags, err := vfs.store.TagsByNames(tagNames)
	for _, tagName := range tagNames {
		if !containsTag(tags, tagName) {
			return nil, fuse.ENOENT
		}
	}

	files, err := vfs.store.QueryFiles(expression)
	if err != nil {
		log.Fatalf("could not query files: %v", err)
	}

	entries := make([]fuse.DirEntry, 0, len(files))
	for _, file := range files {
		linkName := vfs.getLinkName(file)
		entries = append(entries, fuse.DirEntry{Name: linkName, Mode: fuse.S_IFLNK})
	}

	return entries, fuse.OK
}

func (vfs FuseVfs) readTaggedEntryLink(path []string) (string, fuse.Status) {
	log.Infof(2, "BEGIN readTaggedEntryLink(%v)", path)
	defer log.Infof(2, "END readTaggedEntryLink(%v)", path)

	name := path[len(path)-1]

	fileId := vfs.parseFileId(name)
	if fileId == 0 {
		return "", fuse.ENOENT
	}

	file, err := vfs.store.File(fileId)
	if err != nil {
		log.Fatalf("could not find file %v in database.", fileId)
	}

	return file.Path(), fuse.OK
}

func (vfs FuseVfs) getLinkName(file *entities.File) string {
	extension := filepath.Ext(file.Path())
	fileName := filepath.Base(file.Path())
	linkName := fileName[0 : len(fileName)-len(extension)]
	suffix := "." + uitoa(file.Id) + extension

	if len(linkName)+len(suffix) > 255 {
		linkName = linkName[0 : 255-len(suffix)]
	}

	return linkName + suffix
}

func (vfs FuseVfs) tagNamesToIds(tagNames []string) ([]uint, error) {
	tagIds := make([]uint, len(tagNames))

	for index, tagName := range tagNames {
		tag, err := vfs.store.TagByName(tagName)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve tag '%v': %v", tagName, err)
		}
		if tag == nil {
			return nil, nil
		}

		tagIds[index] = tag.Id
	}

	return tagIds, nil
}

func (vfs FuseVfs) saveQuery(query string) {
}

func uitoa(ui uint) string {
	return strconv.FormatUint(uint64(ui), 10)
}

func atoui(str string) (uint, error) {
	ui64, err := strconv.ParseUint(str, 10, 0)
	return uint(ui64), err
}

func containsTag(tags entities.Tags, tagName string) bool {
	for _, tag := range tags {
		if tag.Name == tagName {
			return true
		}
	}

	return false
}

func containsName(tagNames []string, tagName string) bool {
	for _, name := range tagNames {
		if name == tagName {
			return true
		}
	}

	return false
}
