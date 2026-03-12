package main

import (
	"encoding/hex"
	"os"
	"time"

	"github.com/tobischo/gokeepasslib/v3"
)

// KDBXEntry represents a password entry from KeePass database (internal structure)
type KDBXEntry struct {
	UUID         string
	Title        string
	UserName     string
	Password     string
	URL          string
	Notes        string
	Icon         string
	Tags         []string
	CreatedTime  int64
	ModifiedTime int64
	ExpiryTime   *int64
	CustomFields map[string]string
}

// EntryResponse represents a password entry for JSON API response
type EntryResponse struct {
	UUID         string            `json:"uuid"`
	Title        string            `json:"title"`
	UserName     string            `json:"userName"`
	Password     string            `json:"password"`
	URL          string            `json:"url"`
	Notes        string            `json:"notes"`
	Icon         string            `json:"icon"`
	Tags         []string          `json:"tags"`
	CreatedTime  string            `json:"createdTime"`
	ModifiedTime string            `json:"modifiedTime"`
	ExpiryTime   string            `json:"expiryTime,omitempty"`
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

// GroupResponse represents a group for JSON API response
type GroupResponse struct {
	UUID     string          `json:"uuid"`
	Name     string          `json:"name"`
	Icon     string          `json:"icon"`
	Entries  []EntryResponse `json:"entries,omitempty"`
	Groups   []GroupResponse `json:"groups,omitempty"`
}

// KDBXDatabase represents the parsed KeePass database (internal structure)
type KDBXDatabase struct {
	Name       string
	Path       string
	Groups     []KDBXGroup
	AllEntries []KDBXEntry
}

// DatabaseResponse represents the database for JSON API response
type DatabaseResponse struct {
	Name       string          `json:"name"`
	Path       string          `json:"path"`
	Groups     []GroupResponse `json:"groups"`
	AllEntries []EntryResponse `json:"allEntries"`
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

		// Get times as Unix timestamps
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

// entryToResponse converts KDBXEntry to EntryResponse with formatted dates
func entryToResponse(entry KDBXEntry) EntryResponse {
	resp := EntryResponse{
		UUID:         entry.UUID,
		Title:        entry.Title,
		UserName:     entry.UserName,
		Password:     entry.Password,
		URL:          entry.URL,
		Notes:        entry.Notes,
		Icon:         entry.Icon,
		Tags:         entry.Tags,
		CustomFields: entry.CustomFields,
	}

	// Format dates only if they are not zero
	if entry.CreatedTime != 0 {
		resp.CreatedTime = formatDateRU(time.Unix(entry.CreatedTime, 0))
	}
	if entry.ModifiedTime != 0 {
		resp.ModifiedTime = formatDateRU(time.Unix(entry.ModifiedTime, 0))
	}
	if entry.ExpiryTime != nil && *entry.ExpiryTime != 0 {
		resp.ExpiryTime = formatDateRU(time.Unix(*entry.ExpiryTime, 0))
	}

	return resp
}

// entriesToResponse converts a slice of KDBXEntry to EntryResponse
func entriesToResponse(entries []KDBXEntry) []EntryResponse {
	result := make([]EntryResponse, 0, len(entries))
	for _, entry := range entries {
		result = append(result, entryToResponse(entry))
	}
	return result
}

// groupToResponse converts KDBXGroup to GroupResponse with formatted dates
func groupToResponse(group KDBXGroup) GroupResponse {
	return GroupResponse{
		UUID:     group.UUID,
		Name:     group.Name,
		Icon:     group.Icon,
		Entries:  entriesToResponse(group.Entries),
		Groups:   groupsToResponse(group.Groups),
	}
}

// groupsToResponse converts a slice of KDBXGroup to GroupResponse
func groupsToResponse(groups []KDBXGroup) []GroupResponse {
	result := make([]GroupResponse, 0, len(groups))
	for _, group := range groups {
		result = append(result, groupToResponse(group))
	}
	return result
}

// databaseToResponse converts KDBXDatabase to DatabaseResponse
func databaseToResponse(db *KDBXDatabase) *DatabaseResponse {
	return &DatabaseResponse{
		Name:       db.Name,
		Path:       db.Path,
		Groups:     groupsToResponse(db.Groups),
		AllEntries: entriesToResponse(db.AllEntries),
	}
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

// formatDateRU форматирует timestamp в российский формат даты (ДД.ММ.ГГГГ ЧЧ:ММ)
func formatDateRU(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02.01.2006 15:04")
}
