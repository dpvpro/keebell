package main

import (
	"encoding/hex"
	"os"

	"github.com/tobischo/gokeepasslib/v3"
)

// KDBXEntry represents a password entry from KeePass database
type KDBXEntry struct {
	UUID         string            `json:"uuid"`
	Title        string            `json:"title"`
	UserName     string            `json:"userName"`
	Password     string            `json:"password"`
	URL          string            `json:"url"`
	Notes        string            `json:"notes"`
	Icon         string            `json:"icon"`
	Tags         []string          `json:"tags"`
	CreatedTime  int64             `json:"createdTime"`
	ModifiedTime int64             `json:"modifiedTime"`
	ExpiryTime   *int64            `json:"expiryTime,omitempty"`
	CustomFields map[string]string `json:"customFields,omitempty"`
}

// KDBXGroup represents a group from KeePass database
type KDBXGroup struct {
	UUID     string        `json:"uuid"`
	Name     string        `json:"name"`
	Icon     string        `json:"icon"`
	Entries  []KDBXEntry   `json:"entries,omitempty"`
	Groups   []KDBXGroup   `json:"groups,omitempty"`
}

// KDBXDatabase represents the parsed KeePass database
type KDBXDatabase struct {
	Name       string      `json:"name"`
	Path       string      `json:"path"`
	Groups     []KDBXGroup `json:"groups"`
	AllEntries []KDBXEntry `json:"allEntries"`
}

// OpenKDBX opens and decrypts a KDBX database file
func OpenKDBX(path, password string) (*KDBXDatabase, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create new database with credentials
	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(password)

	// Decode the database
	if err := gokeepasslib.NewDecoder(file).Decode(db); err != nil {
		return nil, err
	}

	// Convert to our format
	kdbx := &KDBXDatabase{
		Name:       path,
		Path:       path,
		Groups:     convertGroups(db.Content.Root.Groups),
		AllEntries: collectAllEntries(db.Content.Root.Groups),
	}

	return kdbx, nil
}

// convertGroups converts gokeepasslib groups to our format
func convertGroups(groups []gokeepasslib.Group) []KDBXGroup {
	result := make([]KDBXGroup, 0, len(groups))
	for _, g := range groups {
		kdbxGroup := KDBXGroup{
			UUID:    uuidToString(g.UUID),
			Name:    g.Name,
			Icon:    uuidToString(g.CustomIconUUID),
			Entries: convertEntries(g.Entries),
			Groups:  convertGroups(g.Groups),
		}
		result = append(result, kdbxGroup)
	}
	return result
}

// convertEntries converts gokeepasslib entries to our format
func convertEntries(entries []gokeepasslib.Entry) []KDBXEntry {
	result := make([]KDBXEntry, 0, len(entries))
	for _, e := range entries {
		kdbxEntry := KDBXEntry{
			UUID:     uuidToString(e.UUID),
			Title:    getEntryValue(e, "Title"),
			UserName: getEntryValue(e, "UserName"),
			Password: getEntryValue(e, "Password"),
			URL:      getEntryValue(e, "URL"),
			Notes:    getEntryValue(e, "Notes"),
			Icon:     uuidToString(e.CustomIconUUID),
			Tags:     parseTags(e.Tags),
			CustomFields: make(map[string]string),
		}

		// Get times
		if e.Times.CreationTime != nil {
			kdbxEntry.CreatedTime = e.Times.CreationTime.Time.Unix()
		}
		if e.Times.LastModificationTime != nil {
			kdbxEntry.ModifiedTime = e.Times.LastModificationTime.Time.Unix()
		}
		if e.Times.ExpiryTime != nil {
			expiry := e.Times.ExpiryTime.Time.Unix()
			kdbxEntry.ExpiryTime = &expiry
		}

		// Get custom fields
		for _, fv := range e.Values {
			key := fv.Key
			if key != "Title" && key != "UserName" && key != "Password" && key != "URL" && key != "Notes" {
				kdbxEntry.CustomFields[key] = fv.Value.Content
			}
		}

		result = append(result, kdbxEntry)
	}
	return result
}

// uuidToString converts UUID to string
func uuidToString(uuid gokeepasslib.UUID) string {
	return hex.EncodeToString(uuid[:])
}

// parseTags parses tags string into slice
func parseTags(tags string) []string {
	if tags == "" {
		return []string{}
	}
	result := make([]string, 0)
	// Tags are semicolon-separated in KeePass
	for _, tag := range splitTags(tags) {
		if tag != "" {
			result = append(result, tag)
		}
	}
	return result
}

// splitTags splits tags by semicolon or comma
func splitTags(tags string) []string {
	result := make([]string, 0)
	current := ""
	for _, r := range tags {
		if r == ';' || r == ',' {
			result = append(result, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// getEntryValue gets a value from entry by key
func getEntryValue(entry gokeepasslib.Entry, key string) string {
	for _, v := range entry.Values {
		if v.Key == key {
			return v.Value.Content
		}
	}
	return ""
}

// collectAllEntries collects all entries from all groups recursively
func collectAllEntries(groups []gokeepasslib.Group) []KDBXEntry {
	var entries []KDBXEntry
	for _, g := range groups {
		entries = append(entries, convertEntries(g.Entries)...)
		entries = append(entries, collectAllEntries(g.Groups)...)
	}
	return entries
}
