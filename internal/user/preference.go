package user

import "time"

// Theme represents the type of theme preference
type Theme int

const (
	ThemeLight Theme = iota
	ThemeDark
	ThemeAuto
)

// String returns the string representation of the Theme enum
func (t Theme) String() string {
	return [...]string{"Light", "Dark", "Auto"}[t]
}

// Preference represents user preference settings
type Preference struct {
	ID        int64
	UserID    int64
	Theme     Theme
	Language  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewPreference creates a new user preference with default values
func NewPreference(userID int64) *Preference {
	now := time.Now()
	return &Preference{
		UserID:    userID,
		Theme:     ThemeLight,
		Language:  "en",
		CreatedAt: now,
		UpdatedAt: now,
	}
}
